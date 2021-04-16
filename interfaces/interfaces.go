package interfaces

import (
	"context"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/models"
	"net/http"
)

//go:generate moq -out mock/datastore.go -pkg mock . DataStore
//go:generate moq -out mock/healthcheck.go -pkg mock . HealthChecker
//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/initaliser.go -pkg mock . Initialiser

// DataStore implements the methods required to interact with the database
type DataStore interface {
	Init(config.MongoConfig) (err error)
	Close(ctx context.Context) (err error)
	AddBook(book *models.Book) (err error)
	GetBook(ctx context.Context, id string) (*models.Book, error)
	GetBooks(ctx context.Context) (models.Books, error)
	GetReview(ctx context.Context, reviewID string) (*models.Review, error)
	GetReviews(ctx context.Context, bookID string) (models.Reviews, error)
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	Handler(w http.ResponseWriter, req *http.Request)
	Start(ctx context.Context)
	Stop()
	AddCheck(name string, checker healthcheck.Checker) (err error)
}

type HTTPServer interface {
	ListenAndServe() error
}

type Initialiser interface {
	GetHTTPServer(BindAddr string, router http.Handler) HTTPServer
}
