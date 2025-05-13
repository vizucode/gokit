package config

import (
	"log"

	"github.com/spf13/viper"
)

// Load any configuration like open connection database, open connection redis, monitoring, e.t.c
func Load(serviceName string, configPath string) {

	// serviceName = strings.ToLower(serviceName)
	// serviceName = strings.ReplaceAll(serviceName, "-", "_")
	// serviceName = strings.ReplaceAll(serviceName, " ", "_")

	// load all configuration needed
	// init viper first time
	Config(configPath)
}

func Config(configPath string) {
	viper.AutomaticEnv()
	viper.AddConfigPath(configPath)
	viper.SetConfigFile(".env")
	viper.SetConfigName("env")

	err := viper.ReadInConfig()
	if err == nil {
		log.Print("Config file loaded")
	}
}
