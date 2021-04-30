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
	DefaultMaximumLimit        int `envconfig:"DEFAULT_MAXIMUM_LIMIT"`
	DefaultLimit               int `envconfig:"DEFAULT_LIMIT"`
	DefaultOffset              int `envconfig:"DEFAULT_OFFSET"`
}

type MongoConfig struct {
	BindAddr          string `envconfig:"MONGODB_BIND_ADDR"   json:"-"`
	Database          string `envconfig:"MONGODB_DATABASE"`
	BooksCollection   string `envconfig:"MONGODB_BOOKS_COLLECTION"`
	ReviewsCollection string `envconfig:"MONGODB_REVIEWS_COLLECTION"`
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
			BindAddr:          "localhost:27017",
			Database:          "bookStore",
			BooksCollection:   "books",
			ReviewsCollection: "reviews",
		},
		DefaultMaximumLimit: 1000,
		DefaultLimit:        20,
		DefaultOffset:       0,
	}

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
