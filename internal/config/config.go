package config

import (
	"fmt"

	"github.com/vrischmann/envconfig"

	"webinar/graphql/server/internal/system/database/redis"
)

type Config struct {
	DSN              string
	Port             string
	HashSalt         string
	PlaygroundEnable bool `envconfig:"default=true"`

	Redis *redis.Config
}

func InitConfig(prefix string) (*Config, error) {
	conf := &Config{}
	if err := envconfig.InitWithPrefix(conf, prefix); err != nil {
		return nil, fmt.Errorf("init config error: %w", err)
	}

	return conf, nil
}
