package main

import (
	"context"
	"errors"
	dpHealthCheck "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoDB "github.com/ONSdigital/dp-mongodb/health"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/api"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/initialiser"
	"github.com/cadmiumcat/books-api/mongo"
	"github.com/cadmiumcat/books-api/pagination"
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

	// ErrRegisterHealthCheck represents an error when registering a health checker to the healthcheck
	ErrRegisterHealthCheck = errors.New("error registering checkers for healthcheck")
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

	versionInfo, err := dpHealthCheck.NewVersionInfo(BuildTime, GitCommit, Version)
	if err != nil {
		log.Event(ctx, "could not instantiate health check", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	hc := dpHealthCheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)

	// Initialise database
	mongodb := &mongo.Mongo{}
	err = mongodb.Init(cfg.MongoConfig)
	if err != nil {
		log.Event(ctx, "failed to initialise mongo", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	databaseCollectionBuilder := make(map[dpMongoDB.Database][]dpMongoDB.Collection)
	databaseCollectionBuilder[(dpMongoDB.Database)(mongodb.Database)] =
		[]dpMongoDB.Collection{
			(dpMongoDB.Collection)(mongodb.BooksCollection),
			(dpMongoDB.Collection)(mongodb.ReviewsCollection),
		}

	mongoClient := dpMongoDB.NewClientWithCollections(mongodb.Session.Copy(), databaseCollectionBuilder)

	// Add API checks
	if err := registerCheckers(ctx, &hc, mongoClient); err != nil {
		log.Event(ctx, err.Error(), log.FATAL, log.Error(err))
		os.Exit(1)
	}
	hc.Start(ctx)

	// Initialise server
	svc := initialiser.Service{}
	router := mux.NewRouter()
	svc.Server = initialiser.GetHTTPServer(cfg.BindAddr, router)

	paginator := pagination.NewPaginator(cfg.DefaultLimit, cfg.DefaultOffset, cfg.DefaultMaximumLimit)

	svc.API = api.Setup(ctx, cfg.BindAddr, router, paginator, mongodb, &hc)

	svc.Server.ListenAndServe()

	hc.Stop()
}

// registerCheckers adds the checkers for the provided clients to the health check object
func registerCheckers(ctx context.Context, hc *dpHealthCheck.HealthCheck, mongoClient *dpMongoDB.Client) error {
	var hasErrors bool
	mongoHealth := dpMongoDB.CheckMongoClient{
		Client:      *mongoClient,
		Healthcheck: mongoClient.Healthcheck,
	}
	if err := hc.AddCheck("mongoDB", mongoHealth.Checker); err != nil {
		hasErrors = true
		log.Event(ctx, "error adding mongoDB checker", log.FATAL, log.Error(err))
	}

	if hasErrors {
		log.Event(ctx, ErrRegisterHealthCheck.Error(), log.ERROR)
		return ErrRegisterHealthCheck
	}

	return nil
}
