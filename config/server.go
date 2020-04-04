package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Address string `konaf:"address"`
}

func Read() Config {
	// Global koanf instance. Use . as the key path delimiter. This can be / or anything.
	var k = koanf.New(".")
	// Load default configuration from struct
	if err := k.Load(structs.Provider(Default(), "konaf"), nil); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	// Load configuration from file
	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		log.Println("No config file provided")
	}

	// Prefix indicates environments variables prefix
	const Prefix = "applifier_"

	// Load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "_", ".", -1)
	}), nil); err != nil {
		log.Println("No env variable provided")
	}

	var out Config

	if err := k.Unmarshal("", &out); err != nil {
		logrus.Fatalf("error unmarshalling config: %s", err)
	}

	return out
}
