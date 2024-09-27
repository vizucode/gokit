package types

// Server is the type returned by a classifier server (REST, gRPC)
type Server string

const (
	// REST server
	REST Server = "rest"
	// GRPC server
	GRPC Server = "grpc"
)

func (s Server) String() string {
	return string(s)
}
