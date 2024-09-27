package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vizucode/gokit/factory"
	"github.com/vizucode/gokit/logger"
	"github.com/vizucode/gokit/tracer"
	"github.com/vizucode/gokit/types"
	"github.com/vizucode/gokit/utils/convert"
	"github.com/vizucode/gokit/utils/timezone"

	"github.com/streadway/amqp"
)

type rabbitMqWorker struct {
	ctx        context.Context
	cancelFunc func()
	opt        option
	tz         *time.Location
	ch         *amqp.Channel
	shutdown   chan struct{}
	isShutdown bool
	semaphore  []chan struct{}
	wg         sync.WaitGroup
	channels   []reflect.SelectCase
	handlers   map[string]types.BrokerHandler
}

// New create new rabbitmq consumer
func New(service factory.ServiceFactory, opts ...OptionFunc) factory.ApplicationFactory {
	if service.GetBroker(types.RabbitMQ) == nil {
		log.Fatalf("missing dependencies rabbitmq")
	}

	worker := &rabbitMqWorker{
		opt: getDefaultOption(),
		tz:  timezone.JakartaTz(),
	}
	for _, opt := range opts {
		opt(&worker.opt)
	}

	if reflect.ValueOf(worker.opt.serviceName).IsZero() {
		worker.opt.serviceName = service.Name()
	}

	worker.ctx, worker.cancelFunc = context.WithCancel(context.Background())
	worker.ch = service.GetBroker(types.RabbitMQ).GetConfiguration().(*amqp.Channel)
	worker.shutdown = make(chan struct{}, 1)
	worker.handlers = make(map[string]types.BrokerHandler)

	if h := service.BrokerHandler(types.RabbitMQ); h != nil {
		var hg types.BrokerHandlerGroup
		h.Register(&hg)

		for _, handler := range hg.Handlers {
			worker.opt.exchangeName, worker.opt.queue, worker.opt.isAutoAck = handler.Exchange, handler.Queue, handler.IsAutoAck

			queueChan, err := setupQueueConfig(worker.ch, worker.opt.exchangeName, worker.opt.queue)
			if err != nil {
				panic(err)
			}
			logger.Purple(fmt.Sprintf(`[RABBITMQ-CONSUMER] (exchange): %-15s (queue): %-15s`, `"`+worker.opt.exchangeName+`"`, `"`+worker.opt.queue+`"`))

			worker.channels = append(
				worker.channels, reflect.SelectCase{
					Dir: reflect.SelectRecv, Chan: reflect.ValueOf(queueChan),
				},
			)
			worker.handlers[worker.opt.queue] = handler
			worker.semaphore = append(worker.semaphore, make(chan struct{}, 1))
		}
	}
	logger.PurpleBold(fmt.Sprintf("⇨ RabbitMQ consumer running with %d queue", len(worker.channels)))
	return worker
}

func (r *rabbitMqWorker) Name() string {
	return types.RabbitMQ.String()
}

func (r *rabbitMqWorker) Shutdown(_ context.Context) {
	r.shutdown <- struct{}{}
	r.isShutdown = true
	var runningJob int
	for _, semp := range r.semaphore {
		runningJob += len(semp)
	}

	if runningJob != 0 {
		fmt.Printf("\x1b[34;1mRabbitMQ Broker:\x1b[0m waiting %d job until done...\x1b[0m\n", runningJob)
	}

	r.wg.Wait()
	defer logger.RedBold("Stopping RabbitMQ Broker")
	_ = r.ch.Close()
	r.cancelFunc()
}

func (r *rabbitMqWorker) Serve() {
	for {
		select {
		case <-r.shutdown:
			return
		default:
		}

		chosen, value, ok := reflect.Select(r.channels)
		if !ok {
			continue
		}

		// execute handler
		if msg, ok := value.Interface().(amqp.Delivery); ok {
			r.semaphore[chosen] <- struct{}{}
			if r.isShutdown {
				return
			}

			r.wg.Add(1)
			go func(message amqp.Delivery, index int) {
				r.processMessage(message)
				r.wg.Done()
				<-r.semaphore[index]
			}(msg, chosen)
		}
	}
}

func (r *rabbitMqWorker) processMessage(message amqp.Delivery) {
	start := time.Now().In(r.tz)

	if r.ctx.Err() != nil {
		log.Printf("rabbitmq_consumer > ctx root err: %s", r.ctx.Err())
		return
	}

	ctx := r.ctx
	selectedHandler := r.handlers[message.RoutingKey]

	header := map[string]string{}
	for key, val := range message.Headers {
		header[key] = convert.ToString(val)
	}

	var err error
	trace, ctx := tracer.StartTraceWithContext(ctx, "RabbitMqConsumer")

	// implement logging
	// init logger data
	ol := &logger.DataLogger{
		TimeStart:     start,
		RequestId:     uuid.NewString(),
		Type:          logger.ServiceType(types.RabbitMQ.String()),
		Service:       r.opt.serviceName,
		Endpoint:      fmt.Sprintf("queue: %s", r.opt.queue),
		RequestBody:   string(message.Body),
		RequestMethod: "CONSUME",
		RequestHeader: fmt.Sprintf("Exchange: %s | Routing Key: %s | Header: %v", message.Exchange, message.RoutingKey, header),
	}

	defer func() {
		if re := recover(); re != nil {
			err = fmt.Errorf("%s", re)
		}

		sc := http.StatusOK

		ack := false
		if r.opt.isAutoAck {
			ack = true
		}

		if err != nil {
			trace.SetError(err)

			sc = http.StatusInternalServerError
			ol.ErrorMessage = fmt.Sprintf("%s", err)
		} else {

			ol.Response = "success"
		}

		if ack {
			_ = message.Ack(true)
		} else {
			_ = message.Reject(true)
			_ = message.Nack(true, true)
		}

		trace.SetTag("trace_id", tracer.GetTraceID(ctx))
		ol.StatusCode = sc
		ol.ExecTime = time.Since(start).Seconds()
		logger.Response(ctx, sc, ol.Response, err)
		// finish trace and logging
		trace.Finish()
		ol.Finalize(ctx)
	}()

	// implement locking logging stdout
	var lock = new(logger.Locker)
	// set to context with logger.LogKey as a context key
	ctx = context.WithValue(ctx, logger.LogKey, lock)

	trace.SetTag("exchange", message.Exchange)
	trace.SetTag("routing_key", message.RoutingKey)
	trace.SetTag("body", message.Body)
	trace.SetTag("header", header)

	log.Printf("\x1b[35;3mRabbitMQ Consumer: message consumed, topic = %s\x1b[0m", message.RoutingKey)

	var ec = types.EventContext{}
	ec.SetContext(ctx)
	ec.SetWorkerType(types.RabbitMQ.String())
	ec.SetHandlerRoute(message.RoutingKey)
	ec.SetKey(message.Exchange)
	ec.SetHeader(header)
	_, _ = ec.Write(message.Body)

	handlerFunc := selectedHandler
	if err = handlerFunc.HandlerFunc(&ec); err != nil {
		ec.SetError(err)
	}
}
