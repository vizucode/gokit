package env

import (
	"reflect"
	"time"

	"github.com/spf13/viper"
)

func GetDuration(key string, defaultValues ...time.Duration) time.Duration {
	defaultValue := time.Duration(1) * time.Second
	if len(defaultValues) > 0 {
		defaultValue = defaultValues[0]
	}

	val, err := time.ParseDuration(viper.GetString(key))
	if err != nil {
		return defaultValue
	}

	if reflect.ValueOf(val).IsZero() {
		return defaultValue
	}

	return val
}
