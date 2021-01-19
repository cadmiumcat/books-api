package main

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/api"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/mongo"
	"os"
)

const serviceName = "books-api"

func main() {
	var dataStore api.DataStore

	log.Namespace = serviceName
	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Event(nil, "error retrieving the configuration", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	// TODO: figure out why this does not log the binding address for the database
	log.Event(nil, "loaded configuration", log.INFO, log.Data{"config": cfg})

	dataStore = &mongo.Mongo{}
	err = dataStore.Init(cfg.MongoConfig)
	if err != nil {
		log.Event(nil, "failed to initialise mongo", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	api.Setup(cfg, dataStore)

}
