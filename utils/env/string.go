package env

import (
	"reflect"

	"github.com/spf13/viper"
)

func GetString(key string, defaultValues ...string) (resp string) {
	defaultValue := ""

	val := viper.GetString(key)
	if reflect.ValueOf(val).IsZero() {
		if len(defaultValues) > 0 {
			defaultValue = defaultValues[0]
		}

		return defaultValue
	}

	return val
}
