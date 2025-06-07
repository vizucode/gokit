package factory

import (
	"github.com/vizucode/gokit/abstract"
	"github.com/vizucode/gokit/types"
)

// ServiceFactory abstraction of service
type ServiceFactory interface {
	// Name service name
	Name() string

	// GetApplications return all applications (servers, workers and/or brokers)
	GetApplications() map[string]ApplicationFactory

	// RESTHandler return abstraction of rest-api handler
	RESTHandler() abstract.RestHandler

	// GRPCHandler return abstraction of grpc handler
	GRPCHandler() abstract.GRPCHandler

	// BrokerHandler return abstraction of broker handler by types.Broker
	BrokerHandler(broker types.Broker) abstract.BrokerHandler

	// GetBroker return abstraction of broker configuration by types.Broker
	GetBroker(broker types.Broker) abstract.Broker
}
