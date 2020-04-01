package config

import (
	"bytes"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type ClientConfig struct {
	IP   string `mapstructure:"ip"`
	Port string `mapstructure:"port"`
}

func ReadClient() ClientConfig {
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadConfig(bytes.NewBufferString(ClientDefault)); err != nil {
		log.Fatalf("err: %s", err)
	}

	viper.SetConfigName("config")

	if err := viper.MergeInConfig(); err != nil {
		log.Print("No config file found")
	}

	viper.SetEnvPrefix("applifier")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	var cfg ClientConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("err: %s", err)
	}

	return cfg
}
