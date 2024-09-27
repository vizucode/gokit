package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vizucode/gokit/logger"
	"github.com/vizucode/gokit/tracer"
	"github.com/vizucode/gokit/types"
	"github.com/vizucode/gokit/utils/convert"
	errs "github.com/vizucode/gokit/utils/errors"
	"github.com/vizucode/gokit/utils/monitoring"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// interceptor instance of grpc interceptor
type interceptor struct {
	serviceName string
	host        string
	opt         *option
}

// newInterceptor init an instance interceptor
func newInterceptor(host, sn string) *interceptor {
	return &interceptor{
		host:        host,
		serviceName: sn,
	}
}

func (i *interceptor) chainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		chainer := func(currentInterceptor grpc.UnaryServerInterceptor, currentHandler grpc.UnaryHandler) grpc.UnaryHandler {
			return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return currentInterceptor(currentCtx, currentReq, info, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, req)
	}
}

func (i *interceptor) unaryServerTracerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()

	dl := logger.DataLogger{
		RequestId:     logger.GetRequestId(ctx),
		Type:          logger.ServiceType(string(types.GRPC)),
		Service:       i.serviceName,
		Host:          i.host,
		Endpoint:      info.FullMethod,
		RequestMethod: http.MethodPost,
		TimeStart:     start,
		RequestBody:   convert.ToString(req),
	}

	trace, ctx := tracer.StartTraceWithContext(ctx, fmt.Sprintf("GRPC: %s", info.FullMethod))
	defer func() {
		if r := recover(); r != nil {
			err = status.Errorf(codes.Aborted, "%s", r)
		}
		var sc = http.StatusOK
		if err != nil {
			switch er := err.(type) {
			case *errs.ErrorResponse:
				sc = er.StatusCode()
			default:
				c := status.Code(err)

				switch c {
				case codes.FailedPrecondition, codes.InvalidArgument, codes.Unimplemented:
					sc = http.StatusBadRequest
				case codes.Unauthenticated:
					sc = http.StatusUnauthorized
				case codes.PermissionDenied:
					sc = http.StatusForbidden
				case codes.Unknown, codes.NotFound:
					sc = http.StatusNotFound
				case codes.AlreadyExists:
					sc = http.StatusConflict
				case codes.Aborted, codes.Canceled, codes.DeadlineExceeded, codes.Internal, codes.DataLoss:
					sc = http.StatusInternalServerError
				case codes.OutOfRange:
					sc = http.StatusBadGateway
				case codes.Unavailable:
					sc = http.StatusServiceUnavailable
				case codes.ResourceExhausted:
					sc = http.StatusGatewayTimeout
				default:
					sc = http.StatusOK
				}
			}
		}

		if sc < 1 {
			sc = http.StatusInternalServerError
		}
		logger.Response(ctx, sc, resp, err)

		trace.SetError(err)
		// set error logging
		respBody, _ := json.Marshal(resp)
		if len(respBody) > 1000 {
			trace.Log("response.body.size", len(respBody))
		} else {
			trace.Log("response.body", respBody)
		}
		trace.SetTag("request_id", dl.RequestId)
		trace.SetTag("trace_id", tracer.GetTraceID(ctx))
		trace.Finish()
		dl.Finalize(ctx)
		monitoring.PrometheusRecord(dl.StatusCode, dl.RequestMethod, dl.Endpoint, dl.Service, time.Since(dl.TimeStart))
	}()

	lock := new(logger.Locker)
	ctx = context.WithValue(ctx, logger.LogKey, lock)
	lock.Set(logger.RequestId, dl.RequestId)

	reqBody, _ := json.Marshal(req)
	if len(reqBody) > 1000 {
		trace.Log("request.body.size", len(reqBody))
	} else {
		trace.Log("request.body", len(reqBody))
	}

	resp, err = handler(ctx, req)
	return
}
