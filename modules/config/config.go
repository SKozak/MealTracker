package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration.
// Struct tags define environment variable mapping.
// Compare to Spring's @Value or @ConfigurationProperties.
type Config struct {
	Server ServerConfig
	//Database DatabaseConfig
	App AppConfig
}

type AppConfig struct {
	Env string `envconfig:"APP_ENV" default:"development"`
}

func (a AppConfig) IsDevelopment() bool {
	return a.Env == "development"
}

type ServerConfig struct {
	Host         string        `envconfig:"SERVER_HOST" default:"localhost"`
	Port         int           `envconfig:"SERVER_PORT" default:"8080"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"5s"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
	IdleTimeout  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"120s"`
}

// Load reads configuration from environment variables.
// In Go, we return (value, error) rather than throwing exceptions.
func Load() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	return &cfg, nil
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
