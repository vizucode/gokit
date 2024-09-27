package factory

import "context"

// ApplicationFactory factory for server and/or worker abstraction
type ApplicationFactory interface {
	// Name server application
	Name() string
	// Serve for running server or worker
	Serve()
	// Shutdown stop the server or worker
	Shutdown(ctx context.Context)
}
