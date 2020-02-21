package env

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port string `default:"8080"`
}

func LoadConfig() *Config {
	config := &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		panic("Unable to load config: " + err.Error())
	}
	return config
}
