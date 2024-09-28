package grpc

import (
	"fmt"

	"github.com/vizucode/gokit/utils/env"
)

// OptionFunc setter to set grpc option
type OptionFunc func(*option)

// option grpc
type option struct {
	tcpPort string
	tcpHost string
}

func defaultOption() option {
	return option{
		tcpPort: fmt.Sprintf(":%d", env.GetInteger("GRPC_PORT", 6060)),
	}
}

// SetTCPPort set tcp port
func SetTCPPort(port int) OptionFunc {
	return func(o *option) {
		o.tcpPort = fmt.Sprintf("%d", port)
	}
}

func SetTCPHost(host string) OptionFunc {
	return func(o *option) {
		o.tcpHost = host
	}
}
