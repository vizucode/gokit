package env

import (
	"log"

	"github.com/spf13/viper"
)

func Init() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetConfigName("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("failed to load env")
	}
}
