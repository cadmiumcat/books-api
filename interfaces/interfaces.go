package interfaces

import (
	"context"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/models"
	"net/http"
)

//go:generate moq -out datastoretest/datastore.go -pkg datastoretest . DataStore
//go:generate moq -out mock/healthcheck.go -pkg mock . HealthChecker

// DataStore implements the methods required to interact with the database
type DataStore interface {
	Init(config.MongoConfig) (err error)
	Close(ctx context.Context) (err error)
	AddBook(book *models.Book)
	GetBook(id string) (*models.Book, error)
	GetBooks() (models.Books, error)
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	Handler(w http.ResponseWriter, req *http.Request)
	Start(ctx context.Context)
	Stop()
	AddCheck(name string, checker healthcheck.Checker) (err error)
}