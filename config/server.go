package config

import (
	"bytes"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

//// Global koanf instance. Use . as the key path delimiter. This can be / or anything.
//var k = koanf.New(".")
//
//func ReadServer() ServerConfig {
//	// Load JSON config.
//	if err := k.Load(file.Provider("./serverConf.yml"), yaml.Parser()); err != nil {
//		log.Fatalf("error loading config: %v", err)
//	}
//
//	var out ServerConfig
//
//	// Quick unmarshal.
//	k.Unmarshal("parent1.child1", &out)
//	fmt.Println(out)
//
//	// Unmarshal with advanced config.
//	//out = childStruct{}
//	//k.UnmarshalWithConf("parent1.child1", &out, koanf.UnmarshalConf{Tag: "koanf"})
//	//fmt.Println(out)
//	return out
//}

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
