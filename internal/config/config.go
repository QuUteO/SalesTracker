package config

import (
	"fmt"

	"github.com/wb-go/wbf/config/cleanenv-port"
)

type Config struct {
	Env  string `yaml:"env"`
	DSN  string `yaml:"dsn"`
	Addr string `yaml:"addr"`
}

func New(path string) (*Config, error) {
	var cfg Config
	if err := cleanenvport.LoadPath(path, &cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &cfg, nil
}
