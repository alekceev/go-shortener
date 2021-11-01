package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port              string `envconfig:"PORT" default:"8080" required:"true"`
	ReadTimeout       int    `envconfig:"READ_TIMEOUT" default:"30" required:"true"`
	WriteTimeout      int    `envconfig:"WRITE_TIMEOUT" default:"30" required:"true"`
	ReadHeaderTimeout int    `envconfig:"READ_HEADER_TIMEOUT" default:"30" required:"true"`
	DSN               string `envconfig:"DSN" default:"" required:"true"`
}

// init config
func Get() (Config, error) {
	config := Config{}

	err := envconfig.Process("", &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// Get Server host:port
func (c *Config) Addr() string {
	return ":" + c.Port
}
