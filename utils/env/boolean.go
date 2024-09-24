package env

import (
	"strconv"

	"github.com/spf13/viper"
)

func GetBool(key string, defaultValues ...bool) bool {
	defaultValue := false
	if len(defaultValues) > 0 {
		defaultValue = defaultValues[0]
	}

	val, err := strconv.ParseBool(viper.GetString(key))
	if err != nil {
		return defaultValue
	}

	return val
}
