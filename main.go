package main

import (
	"context"
	hc "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/api"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/initialiser"
	"github.com/cadmiumcat/books-api/interfaces"
	"github.com/cadmiumcat/books-api/mongo"
	"github.com/gorilla/mux"
	"os"
)

const serviceName = "books-api"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	ctx := context.Background()

	log.Namespace = serviceName
	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Event(ctx, "error retrieving the configuration", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	log.Event(ctx, "loaded configuration", log.INFO, log.Data{"config": cfg})

	// Initialise Health Check?
	versionInfo, err := hc.NewVersionInfo(BuildTime, GitCommit, Version)
	if err != nil {
		log.Event(ctx, "could not instantiate health check", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	hc := hc.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)

	// Initialise database
	var dataStore interfaces.DataStore
	dataStore = &mongo.Mongo{}
	err = dataStore.Init(cfg.MongoConfig)
	if err != nil {
		log.Event(ctx, "failed to initialise mongo", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	if err = hc.AddCheck("mongoDB", dataStore.Checker); err != nil {
		log.Event(ctx, "failed to add healthcheck", log.FATAL, log.Error(err))
	}

	// Initialise server
	svc := initialiser.Service{}
	router := mux.NewRouter()
	svc.Server = initialiser.GetHTTPServer(cfg.BindAddr, router)

	svc.API = api.Setup(ctx, cfg.BindAddr, router, dataStore, &hc)

	hc.Start(ctx)

	svc.Server.ListenAndServe()

	hc.Stop()
}
