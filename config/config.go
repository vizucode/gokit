package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Load any configuration like open connection database, open connection redis, monitoring, e.t.c
func Load(serviceName string) {

	serviceName = strings.ToLower(serviceName)
	serviceName = strings.ReplaceAll(serviceName, "-", "_")
	serviceName = strings.ReplaceAll(serviceName, " ", "_")

	// load all configuration needed
	// init viper first time
	Config()
	fmt.Println(serviceName)
}

func Config() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetConfigName("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("failed to load env")
	}
}
