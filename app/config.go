package app

import (
	"github.com/JoseBeteta/surfe/app/infrastructure/common/http"
	"time"
)

const (
	serviceName               = "surfe"
	defaultConcurrentRoutines = 20
)

// Config overall configuration struct
type Config struct {
	ServiceName string
	Server      http.Config
}

// ConfigLoader interface for the config loader
type ConfigLoader interface {
	Load(cfg any) error
}

// LoadConfig loads the config from the file into the struct
func LoadConfig(loader ConfigLoader) (*Config, error) {
	cfg := &Config{}
	err := loader.Load(cfg)
	cfg.ServiceName = serviceName

	cfg.Server.ReadTimeout = 5 * time.Second
	cfg.Server.HandlerTimeout = 5 * time.Second
	cfg.Server.WriteTimeout = cfg.Server.ReadTimeout + cfg.Server.HandlerTimeout
	cfg.Server.IdleTimeout = 15 * time.Second

	return cfg, err
}
