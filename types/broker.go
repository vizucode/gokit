package types

// Broker is the type returned by a classifier broker
type Broker string

const (
	// RabbitMQ Broker
	RabbitMQ Broker = "rabbit-mq"
	// Solace Broker
	Solace Broker = "solace"
	// NSQ Broker
	NSQ Broker = "nsq"
	// Kafka Broker
	Kafka Broker = "kafka"
)

func (b Broker) String() string {
	return string(b)
}

// BrokerHandlerFunc type abstract for each broker implementation
type BrokerHandlerFunc func(ec *EventContext) error

type BrokerHandlerOption func(*BrokerHandler)

// BrokerHandler instance
type BrokerHandler struct {
	Topic            string // topic broker
	Exchange         string // exchange of broker
	Queue            string // queue message
	IsQueueDurable   bool   // durable of queue
	IsQueueExclusive bool   // queue exclusive
	Channel          string // channel app name
	IsAutoAck        bool   // auto acknowledgement
	HandlerFunc      BrokerHandlerFunc
}

// BrokerHandlerGroup group of broker handlers by topic, exchange, or queue with channels
type BrokerHandlerGroup struct {
	Handlers []BrokerHandler
}

// AddBrokerHandler method from BrokerHandlerGroup
func (bhg *BrokerHandlerGroup) AddBrokerHandler(handlerFunc BrokerHandlerFunc, opts ...BrokerHandlerOption) {
	bh := BrokerHandler{HandlerFunc: handlerFunc, IsQueueDurable: true, IsQueueExclusive: false}

	for _, opt := range opts {
		opt(&bh)
	}
	bhg.Handlers = append(bhg.Handlers, bh)
}

// SetBrokerTopic set topic into broker
func SetBrokerTopic(topic string) BrokerHandlerOption {
	return func(bh *BrokerHandler) {
		bh.Topic = topic
	}
}

// SetBrokerExchange set exchange of broker
func SetBrokerExchange(exchange string) BrokerHandlerOption {
	return func(bh *BrokerHandler) {
		bh.Exchange = exchange
	}
}

// SetBrokerQueue set queue
func SetBrokerQueue(queue string) BrokerHandlerOption {
	return func(bh *BrokerHandler) {
		bh.Queue = queue
	}
}

// SetBrokerChannel set channel to broker
func SetBrokerChannel(channel string) BrokerHandlerOption {
	return func(bh *BrokerHandler) {
		bh.Channel = channel
	}
}

// SetBrokerDurable set channel to broker
func SetBrokerDurable(durable bool) BrokerHandlerOption {
	return func(bh *BrokerHandler) {
		bh.IsQueueDurable = durable
	}
}

// SetBrokerExclusive set channel to broker
func SetBrokerExclusive(exclusive bool) BrokerHandlerOption {
	return func(bh *BrokerHandler) {
		bh.IsQueueExclusive = exclusive
	}
}

// SetBrokerAutoAck set channel to broker
func SetBrokerAutoAck(autoAck bool) BrokerHandlerOption {
	return func(bh *BrokerHandler) {
		bh.IsAutoAck = autoAck
	}
}
