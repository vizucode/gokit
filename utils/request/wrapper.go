package request

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/vizucode/gokit/logger"
	"github.com/vizucode/gokit/tracer"
	"github.com/vizucode/gokit/utils/monitoring"
	"github.com/vizucode/gokit/utils/timezone"
)

func (r *request) wrapper(ctx context.Context, payload []byte, method string) ([]byte, int, error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, fmt.Sprintf("RequestClient:%s", strings.ToUpper(r.serviceTarget)))
	defer trace.Finish()

	var (
		timezone = timezone.JakartaTz()
		start    = time.Now().In(timezone)
		tp       logger.ThirdParty
	)
	_, shortUrl := filterUrl(r.url)

	tp.Method = method
	tp.URL = r.url
	tp.RequestHeader = parseHeader(r.header)
	tp.ServiceTarget = r.serviceTarget

	trace.SetTag("request_method", tp.Method)
	trace.SetTag("request_url", tp.URL)
	trace.SetTag("request_header", tp.RequestHeader)

	if payload != nil {
		tp.RequestBody = parseBodyPayload(payload)
		trace.SetTag("request_body", tp.RequestBody)
	}

	res, status, err := r.do(payload, method)

	trace.SetTag("response_status_code", status)

	tp.StatusCode = status
	if err != nil {
		tp.Response = err.Error()
		trace.SetError(err)
	}

	if res != nil {
		trace.SetTag("response_body", res)
		if len(res) > 1000 {
			tp.Response = "success request"
		} else {
			tp.Response = string(res)
		}
	}

	since := time.Since(start)
	tp.ExecTime = since.Seconds()
	// storing data third party request and response to context and prometheus
	tp.Store(ctx)
	monitoring.PrometheusRecord(tp.StatusCode, tp.Method, shortUrl, tp.ServiceTarget, since)

	return res, status, err
}

func filterUrl(url string) (string, string) {
	hideDynamicPath := `((628|08)(31|32|33|38|591|598)\d{6,10}|\d{5,13})`
	regex := regexp.MustCompile(hideDynamicPath)

	u := regex.ReplaceAllString(url, "_censored_")

	return u, parseURL(u)
}

func parseURL(urls string) string {
	u, err := url.Parse(urls)
	if err != nil {
		return urls
	}
	return u.Path
}

func parseBodyPayload(payload []byte) string {
	return string(payload)
}

func parseHeader(header http.Header) string {
	var h string
	for key, val := range header {
		h += fmt.Sprintf("%s: %s  ", key, strings.Join(val, ","))
	}

	return strings.TrimSpace(h)
}

// convertInterfaceToString converts any given interface to a JSON string representation.
// // It returns an empty string in case of an error during marshaling.
// func convertInterfaceToString(i interface{}) string {
// 	// Marshal the interface into JSON bytes
// 	buf, err := json.Marshal(i)
// 	if err != nil {
// 		// Log the error for debugging purposes
// 		logger.Log.Errorf(context.TODO(), "Error marshaling interface to JSON: %v", err)
// 		return "" // Return empty string on error
// 	}

// 	// Return the JSON string representation
// 	return string(buf)
// }
