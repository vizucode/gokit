package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/vizucode/gokit/utils/timezone"
)

// InitializeCron for init first context from scheduler
func InitializeCron(endpoint string) (context.Context, DataLogger) {
	var (
		timezone = timezone.JakartaTz()
		start    = time.Now().In(timezone)
		lock     = new(Locker)
		dl       DataLogger
	)

	function, _, _, _ := runtime.Caller(1)
	functionName := runtime.FuncForPC(function).Name()

	dl.RequestId = uuid.New().String()
	dl.Type = cron
	dl.Service = getServiceName()
	dl.Host = functionName
	dl.Endpoint = endpoint
	dl.RequestMethod = http.MethodGet
	dl.TimeStart = start

	ctx := context.WithValue(context.Background(), LogKey, lock)

	lock.Set(RequestId, dl.RequestId)

	return ctx, dl
}

func getServiceName() string {
	return filepath.Base(os.Args[0])
}

// Store for storing data context to third parties
func (th ThirdParty) Store(ctx context.Context) {
	var (
		data []ThirdParty
		val  Values
	)

	val, ok := extract(ctx)
	if !ok {
		return
	}

	tmp, ok := val.LoadAndDelete(_ThirdParties)
	if ok {
		data = tmp.([]ThirdParty)
	}

	// check response size is more than 1k character
	var resp = th.Response
	b, _ := json.Marshal(resp)
	if len(b) > 1000 {
		// if character more than 1k, make it simple response
		resp = "success request"
	}
	th.Response = resp // replace with the new one, if logic is valid

	data = append(data, th)

	val.Set(_ThirdParties, data)
}

// DeviceUser is record data used device by user session
func DeviceUser(ctx context.Context, userCode, deviceType, deviceOs, brand, model string) {
	value, ok := extract(ctx)
	if !ok {
		return
	}

	device := fmt.Sprintf("%s %s - %s %s", deviceType, deviceOs, brand, model)

	value.Set(_UserCode, userCode)
	value.Set(_Device, device)
}

// Response is record data response to context
func Response(ctx context.Context, status int, res interface{}, err error) {
	value, ok := extract(ctx)
	if !ok {
		return
	}
	// check total string for response
	if res != nil {
		buf, err := json.Marshal(res)
		if err != nil {
			return
		}

		if len(buf) > 1000 {
			value.Set(_Response, success)
		} else {
			value.Set(_Response, string(buf))
		}
	}
	value.Set(_StatusCode, status)

	if err != nil {
		value.Set(_ErrorMessage, err.Error())
	}
}
