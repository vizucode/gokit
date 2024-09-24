package env

import (
	"reflect"

	"github.com/spf13/viper"
)

func GetFloat(key string, defaultValues ...float64) float64 {
	var defaultValue float64 = -1

	val := viper.GetFloat64(key)
	if reflect.ValueOf(val).IsZero() {
		if len(defaultValues) > 0 {
			defaultValue = defaultValues[0]
		}

		return defaultValue
	}

	return val
}
