package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoURI string `envconfig:"mongo_uri" json:"-"`
	Env      string
}

func LoadConfig(c *Config) error {
	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}
