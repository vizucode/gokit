package logger

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vizucode/gokit/utils/env"
	"github.com/vizucode/gokit/utils/monitoring"
)

// Finalize load from context and delete data context
func (d *DataLogger) Finalize(ctx context.Context) {
	value, ok := extract(ctx)
	if !ok {
		return
	}

	if i, ok := value.LoadAndDelete(_StatusCode); ok && i != nil {
		d.StatusCode = i.(int)
	}

	if i, ok := value.LoadAndDelete(_Response); ok && i != nil {
		d.Response = i
	}

	if i, ok := value.LoadAndDelete(_ThirdParties); ok && i != nil {
		d.ThirdParties = i.([]ThirdParty)
	}

	if i, ok := value.LoadAndDelete(_LogMessages); ok && i != nil {
		d.LogMessages = i.([]LogMessage)
	}

	if i, ok := value.LoadAndDelete(_ErrorMessage); ok && i != nil {
		d.ErrorMessage = i.(string)
	}

	if i, ok := value.LoadAndDelete(_UserCode); ok && i != nil {
		d.UserCode = i.(string)
	}

	if i, ok := value.LoadAndDelete(_Device); ok && i != nil {
		d.Device = i.(string)
	}

	d.ExecTime = time.Since(d.TimeStart).Seconds()

	appEnv := strings.ToUpper(env.GetString("APP_ENV"))
	if len(d.LogMessages) > 5 && !reflect.ValueOf(appEnv).IsZero() && appEnv == "PRODUCTION" {
		d.LogMessages = d.LogMessages[len(d.LogMessages)-5:]
	}

	// delete context GetRequestId and GetSaltKey
	value.Delete(_SaltKey)
	value.Delete(RequestId)

	monitoring.PrometheusRecord(d.StatusCode, d.RequestMethod, d.Endpoint, d.Service, time.Since(d.TimeStart))
	d.write()
}

func (d *DataLogger) write() {
	var (
		level logrus.Level
		// elasticStatus = env.GetBool("ELASTICSEARCH_ENABLED", false)
		// errChan       = make(chan error, 1)
	)

	// if elasticStatus {
	// 	if err := adapter.PublishLogging(context.Background(), d); err != nil {
	// 		errChan <- err
	// 	}

	// 	select {
	// 	case err := <-errChan:
	// 		logrus.Errorf("error send data to elastic %v", err)
	// 		break
	// 	default:
	// 		close(errChan)
	// 		break
	// 	}
	// }

	if d.StatusCode >= 200 && d.StatusCode < 400 {
		level = logrus.InfoLevel
	} else if d.StatusCode >= 400 && d.StatusCode < 500 {
		level = logrus.WarnLevel
	} else {
		level = logrus.ErrorLevel
	}

	Logrus().WithField("data", d).Log(level, d.Type.String())
}
