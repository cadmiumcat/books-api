package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Configuration struct {
	BindAddr                   string `envconfig:"BIND_ADDR"`
	HealthCheckCriticalTimeout time.Duration
	HealthCheckInterval        time.Duration
	MongoConfig                MongoConfig
}

type MongoConfig struct {
	BindAddr   string `envconfig:"MONGODB_BIND_ADDR"   json:"-"`
	Collection string `envconfig:"MONGODB_COLLECTION"`
	Database   string `envconfig:"MONGODB_DATABASE"`
}

var cfg *Configuration

// Get configures the application and returns the configuration
func Get() (*Configuration, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Configuration{
		BindAddr:                   ":8080",
		HealthCheckCriticalTimeout: 90 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		MongoConfig: MongoConfig{
			BindAddr:   "localhost:27017",
			Collection: "books",
			Database:   "bookStore",
		},
	}

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
