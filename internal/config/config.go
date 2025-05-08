package config

import (
	"github.com/kelseyhightower/envconfig"
)

const prefix = "HEALTHMONITOR"

type Config struct {
	ListenHTTP   string `envconfig:"LISTEN_HTTP" required:"true"`
	DatabaseHost string `envconfig:"DATABASE_HOST" required:"true"`
	DatabaseName string `envconfig:"DATABASE_NAME" required:"true"`
	DatabaseUser string `envconfig:"DATABASE_USER" required:"true"`
	DatabasePass string `envconfig:"DATABASE_PASS" required:"true"`
}

func New() (*Config, error) {
	cfg := Config{}

	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
