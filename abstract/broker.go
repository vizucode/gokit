package abstract

import "github.com/vizucode/gokit/types"

// Broker message broker abstraction
type Broker interface {
	GetPublisher() Publisher
	GetName() types.Broker
	GetConfiguration() interface{}

	Closer
}
