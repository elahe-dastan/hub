package config

import (
	"bytes"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

func ReadServer() ServerConfig {
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadConfig(bytes.NewBufferString(ServerDefault)); err != nil {
		log.Fatalf("err: %s", err)
	}

	viper.SetConfigName("config")

	if err := viper.MergeInConfig(); err != nil {
		log.Print("No config file found")
	}

	viper.SetEnvPrefix("applifier")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	var cfg ServerConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("err: %s", err)
	}

	return cfg
}
