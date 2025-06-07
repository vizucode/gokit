package logger

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vizucode/gokit/tracer"
	"github.com/vizucode/gokit/utils/monitoring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// interceptor instance of intercept gRPC request
type interceptor struct {
	host string
}

// InterceptorInterface defines all method to be implemented by controllers/handlers
// for intercept gRPC request
type InterceptorInterface interface {
	// ChainUnaryClient start intercept grpc client request
	ChainUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor

	// UnaryClientTracerInterceptor trace the outcoming request (from client) to grpc server
	UnaryClientTracerInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error)
}

// NewInterceptor constructor interceptor grpc request
func NewInterceptor(host string) InterceptorInterface {
	return &interceptor{
		host: host,
	}
}

func (i *interceptor) ChainUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	n := len(interceptors)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		chainer := func(currentIntercept grpc.UnaryClientInterceptor, currentInvoker grpc.UnaryInvoker) grpc.UnaryInvoker {
			return func(currentCtx context.Context, currentMethod string, currentReq, currentResp interface{}, currentClientConn *grpc.ClientConn, currentOpts ...grpc.CallOption) error {
				return currentIntercept(currentCtx, currentMethod, currentReq, currentResp, currentClientConn, currentInvoker, currentOpts...)
			}
		}

		chainedInvoker := invoker
		for i := n - 1; i >= 0; i-- {
			chainedInvoker = chainer(interceptors[i], chainedInvoker)
		}

		return chainedInvoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *interceptor) UnaryClientTracerInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	start := time.Now()

	trace, ctx := tracer.StartTraceWithContext(ctx, fmt.Sprintf("GRPCRepository:%s", method))
	tp := ThirdParty{
		ServiceTarget: method,
		URL:           cc.Target(),
		RequestBody:   convertInterfaceToString(req),
		Method:        http.MethodPost,
	}

	var sc = http.StatusOK

	defer func() {
		if r := recover(); r != nil {
			err = status.Errorf(codes.Aborted, "%s", r)
		}
		if err != nil {
			sc = http.StatusBadRequest
			trace.SetError(err)
		}

		end := time.Since(start)
		tp.Response = convertInterfaceToString(reply)
		tp.StatusCode = sc
		tp.ExecTime = end.Seconds()
		tp.Store(ctx)

		trace.SetTag("response_body", reply)
		trace.SetTag("status_code", sc)
		trace.Finish()

		monitoring.PrometheusRecord(sc, tp.Method, tp.URL, getServiceName(), end)
	}()

	trace.SetTag("trace_id", tracer.GetTraceID(ctx))
	trace.SetTag("url", tp.URL)
	trace.SetTag("target", tp.ServiceTarget)
	trace.SetTag("request_body", req)

	err = invoker(ctx, method, req, reply, cc, opts...)
	return
}
