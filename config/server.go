package config

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/structs"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Address string `konaf:"address"`
}

func ReadServer() ServerConfig {
	// Global koanf instance. Use . as the key path delimiter. This can be / or anything.
	var k = koanf.New(".")
	// Load default configuration from file
	if err := k.Load(structs.Provider(ServerDefault(), "konaf"), nil); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	var out ServerConfig

	if err := k.Unmarshal("", &out); err != nil {
		logrus.Fatalf("error unmarshalling config: %s", err)
	}
	return out
}
