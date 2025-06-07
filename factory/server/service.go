package server

import (
	"github.com/vizucode/gokit/abstract"
	"github.com/vizucode/gokit/factory"
	"github.com/vizucode/gokit/factory/server/grpc"
	"github.com/vizucode/gokit/factory/server/rabbitmq"
	"github.com/vizucode/gokit/factory/server/rest"
	"github.com/vizucode/gokit/types"
)

// ServiceFunc setter to set service instance
type ServiceFunc func(*service)

// service instance
type service struct {
	name                 string
	brokerHandler        map[types.Broker]abstract.BrokerHandler
	brokerHandlerOptions map[types.Broker]interface{}
	brokers              map[types.Broker]abstract.Broker
	rest                 abstract.RestHandler
	restOptions          []rest.OptionFunc
	grpc                 abstract.GRPCHandler
	grpcOptions          []grpc.OptionFunc
	applications         map[string]factory.ApplicationFactory
}

// SetServiceName setter
func SetServiceName(svcName string) ServiceFunc {
	return func(s *service) {
		s.name = svcName
	}
}

// SetBrokerHandler setter brokerHandler
func SetBrokerHandler(broker types.Broker, handler abstract.BrokerHandler) ServiceFunc {
	return func(s *service) {
		if len(s.brokerHandler) < 1 || s.brokerHandler == nil {
			s.brokerHandler = make(map[types.Broker]abstract.BrokerHandler)
		}

		s.brokerHandler[broker] = handler
	}
}

// SetBrokerHandlerOptions setter broker handler options
func SetBrokerHandlerOptions(broker types.Broker, opts ...interface{}) ServiceFunc {
	return func(s *service) {
		if len(s.brokerHandlerOptions) < 1 || s.brokerHandlerOptions == nil {
			s.brokerHandlerOptions = make(map[types.Broker]interface{})
		}

		s.brokerHandlerOptions[broker] = opts
	}
}

// SetBroker setter broker
func SetBroker(brokerName types.Broker, broker abstract.Broker) ServiceFunc {
	return func(s *service) {
		if len(s.brokers) < 1 || s.brokers == nil {
			s.brokers = make(map[types.Broker]abstract.Broker)
		}

		s.brokers[brokerName] = broker
	}
}

// SetRestHandler setter
func SetRestHandler(restHandler abstract.RestHandler) ServiceFunc {
	return func(s *service) {
		s.rest = restHandler
	}
}

// SetRestHandlerOptions setter options for rest handler
func SetRestHandlerOptions(opts ...rest.OptionFunc) ServiceFunc {
	return func(s *service) {
		s.restOptions = opts
	}
}

// SetGrpcHandler setter
func SetGrpcHandler(grpcHandler abstract.GRPCHandler) ServiceFunc {
	return func(s *service) {
		s.grpc = grpcHandler
	}
}

// SetGrpcHandlerOptions setter options for grpc
func SetGrpcHandlerOptions(opts ...grpc.OptionFunc) ServiceFunc {
	return func(s *service) {
		s.grpcOptions = opts
	}
}

// NewService initiate service
func NewService(serviceFuncs ...ServiceFunc) factory.ServiceFactory {
	svc := &service{}
	for _, service := range serviceFuncs {
		service(svc)
	}

	return svc
}

func (s *service) Name() string {
	return s.name
}

func (s *service) GetApplications() map[string]factory.ApplicationFactory {
	// initiate when map not yet declare
	// handling error nil pointer reference
	if len(s.applications) < 1 || s.applications == nil {
		s.applications = make(map[string]factory.ApplicationFactory)
	}

	// set default rest handler when rest handler is nil
	if s.rest == nil {
		s.rest = defaultRestHandler()
	}

	// set rest handler into applications factory
	if _, ok := s.applications[types.REST.String()]; !ok {
		s.applications[types.REST.String()] = rest.New(s, s.restOptions...)
	}

	// set grpc handler into application factory
	if s.grpc != nil {
		if _, ok := s.applications[types.GRPC.String()]; !ok {
			s.applications[types.GRPC.String()] = grpc.New(s, s.grpcOptions...)
		}
	}

	// set rabbit-mq handler into applications factory
	if s.brokerHandler[types.RabbitMQ] != nil {
		if _, ok := s.applications[types.RabbitMQ.String()]; !ok {
			var rmqOpts = make([]rabbitmq.OptionFunc, 0)
			if in, ok := s.brokerHandlerOptions[types.RabbitMQ]; ok {
				if val, ok := in.([]rabbitmq.OptionFunc); ok {
					rmqOpts = val
				}
			}

			// initiate rabbit-mq server here
			s.applications[types.RabbitMQ.String()] = rabbitmq.New(s, rmqOpts...)
		}
	}

	// return all applications factory
	return s.applications
}

func (s *service) RESTHandler() abstract.RestHandler {
	return s.rest
}

func (s *service) GRPCHandler() abstract.GRPCHandler {
	return s.grpc
}

func (s *service) BrokerHandler(broker types.Broker) abstract.BrokerHandler {
	return s.brokerHandler[broker]
}

func (s *service) GetBroker(broker types.Broker) abstract.Broker {
	return s.brokers[broker]
}
