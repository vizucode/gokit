package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/vizucode/gokit/utils/timezone"
)

const layoutDateTime = "2006-01-02T15:04:05-07:00"

type LogFormatted struct {
	logrus.Formatter
}

func Logrus() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&LogFormatted{
		&logrus.JSONFormatter{
			TimestampFormat: layoutDateTime,
		},
	})

	return log
}

func (l *LogFormatted) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.In(timezone.JakartaTz())
	return l.Formatter.Format(e)
}
