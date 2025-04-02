package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	Pg       string        `yaml:"storage_path" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC     GRPCConfig    `yaml:"grpc"`
	Log      string        `yaml:"log_level"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
