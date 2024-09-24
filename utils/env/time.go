package env

import (
	"time"

	"github.com/spf13/viper"
)

type OptionTime func(t *times)

type times struct {
	defaultTime     time.Time
	format          string
	parseInLocation bool
	timezone        *time.Location
}

func defaultTimes() times {
	return times{
		format:          "2006-01-02 15:04:05",
		parseInLocation: false,
		timezone:        time.UTC,
	}
}

func SetFormatTime(format string) OptionTime {
	return func(t *times) {
		t.format = format
	}
}

func SetDefaultValue(dv time.Time) OptionTime {
	return func(t *times) {
		t.defaultTime = dv
	}
}

func SetParseInLocation(b bool) OptionTime {
	return func(t *times) {
		t.parseInLocation = b
	}
}

func SetTimezone(tz *time.Location) OptionTime {
	return func(t *times) {
		t.timezone = tz
	}
}

func GetTime(key string, options ...OptionTime) time.Time {
	t := defaultTimes()
	for _, option := range options {
		option(&t)
	}
	val, err := time.Parse(t.format, viper.GetString(key))
	if err != nil {
		return t.defaultTime
	}

	if t.parseInLocation {
		return val.In(t.timezone)
	}

	return val
}
