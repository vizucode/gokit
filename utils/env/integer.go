package env

import (
	"reflect"

	"github.com/spf13/viper"
)

func GetInteger(key string, defaultValues ...int) (resp int) {
	defaultValue := 0

	val := viper.GetInt(key)
	if reflect.ValueOf(val).IsZero() {
		if len(defaultValues) > 0 {
			defaultValue = defaultValues[0]
		}

		return defaultValue
	}

	return val
}
