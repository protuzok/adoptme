package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App  App
		HTTP HTTP
		Log  Log
		PG   PG
		JWT  JWT
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	JWT struct {
		Secret      string        `env:"JWT_SECRET,required"`
		TokenExpiry time.Duration `env:"JWT_TOKEN_EXPIRY" envDefault:"24h"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
