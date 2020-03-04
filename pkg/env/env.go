package env

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          string `default:"8080"`
	OauthClientID string `required:"true" split_words:"true"`
	LoginDomain   string `required:"true" split_words:"true"`
	ProjectID     string `required:"true" split_words:"true"`
	AdminEmails   string `required:"false" split_words:"true"`
}

func LoadConfig() *Config {
	config := &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		panic("Unable to load config: " + err.Error())
	}
	return config
}
