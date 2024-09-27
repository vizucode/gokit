package abstract

import "context"

// Closer abstraction to close the connection
type Closer interface {
	Disconnect(ctx context.Context) error
}
