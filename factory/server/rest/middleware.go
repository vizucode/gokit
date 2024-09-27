package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vizucode/gokit/logger"
	"github.com/vizucode/gokit/tracer"
	"github.com/vizucode/gokit/utils/timezone"
)

func (r *rest) restTraceLogger(c *fiber.Ctx) error {
	ctx := c.UserContext()
	start := time.Now().In(timezone.JakartaTz())

	var err error
	var sc = http.StatusOK
	var resp string

	requestId := c.Get("x-request-id")
	if reflect.ValueOf(requestId).IsZero() {
		requestId = uuid.NewString()
	}

	// dump and/or parse url, header, and body
	parseUrl := parseUrl(c)
	dumpHeader := string(dumpHeaderFromRequest(c))
	dumpBody := dumpBodyFromRequest(c)
	// init logger
	dl := logger.DataLogger{
		RequestId:     requestId,
		Type:          logger.ServiceType("rest_api"),
		TimeStart:     start,
		Service:       r.service.Name(),
		Host:          c.BaseURL(),
		RequestMethod: c.Method(),
		RequestHeader: dumpHeader,
		RequestBody:   dumpBody,
		Endpoint:      parseUrl,
	}

	// start open tracing with jaeger
	operationName := fmt.Sprintf("%s %s", c.Method(), parseUrl)
	trace, ctx := tracer.StartTraceWithContext(ctx, operationName)
	defer func() {
		if re := recover(); re != nil {
			err = fmt.Errorf("%s", re)
		}

		if err != nil {
			trace.SetError(err)
		}

		if sc < http.StatusOK {
			sc = http.StatusInternalServerError
		}

		// set response
		logger.Response(ctx, sc, resp, err)
		// get all data logging from context with mutext
		dl.Finalize(ctx)
		// finish the tracing
		trace.Finish()
	}()

	// set logger into context with key LogKey
	lock := new(logger.Locker)
	ctx = context.WithValue(ctx, logger.LogKey, lock)
	lock.Set(logger.RequestId, dl.RequestId)

	// set current context into fiber-context
	c.SetUserContext(ctx)

	// set log request
	trace.SetTag("tracer_id", tracer.GetTraceID(ctx))
	trace.SetTag("request_id", dl.RequestId)
	trace.SetTag("app_version", c.Get("x-app-version"))
	trace.SetTag("http.method", c.Method())
	trace.SetTag("http.url", parseUrl)
	trace.SetTag("http.original_url", c.OriginalURL())
	trace.SetTag("http.request", dumpHeader)
	trace.SetTag("http.request_body", dl.RequestBody)

	// next handler
	err = c.Next()
	trace.SetTag("user_code", dl.UserCode)
	trace.SetTag("device", dl.Device)

	// set open tracing response
	trace.SetTag("http.status_code", c.Response().StatusCode())
	sc = c.Response().StatusCode()
	var respBody = c.Response().Body()
	if len(respBody) > 1000 {
		trace.SetTag("response.body", "success request")
		trace.SetTag("response.body.size", len(respBody))
		resp = "success request"
	} else {
		trace.SetTag("response.body", respBody)
		resp = string(respBody)
	}

	return err
}

func dumpHeaderFromRequest(c *fiber.Ctx) []byte {
	var uri string
	var header []byte

	uri = c.OriginalURL()
	header, _ = json.Marshal(c.GetReqHeaders())

	s := fmt.Sprintf("%s %s %s", c.Method(), uri, string(header))
	return []byte(s)
}

func parseUrl(c *fiber.Ctx) string {
	var url = string(c.Request().URI().Path())

	for key, val := range c.AllParams() {
		// when url value contains with value from urlParams
		if strings.Contains(url, val) {
			// get index characters
			index := strings.Index(url, val)

			// replace value with key of params
			url = fmt.Sprintf("%s:%s%s", url[:index], key, url[len(val)+index:])
		}
	}

	return url
}

// dumpBodyFromRequest for getting all request from payload body http_rest_api
func dumpBodyFromRequest(c *fiber.Ctx) string {
	var reqBody string

	// NOTES:
	// - before version v1.1.0, only support formValue with key 'content' (application/x-www-formurlencoded)
	// - after version v1.1.0, support both. form-value with key 'content' or raw json
	reqBody = c.FormValue("content")
	// when formValue is not exists, check on raw json
	if reflect.ValueOf(reqBody).IsZero() || reqBody == "" {
		reqBody = string(c.Request().Body())
	}

	return reqBody
}
