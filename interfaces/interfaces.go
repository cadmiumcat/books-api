package interfaces

import (
	"context"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/models"
	"net/http"
)

//go:generate moq -out mock/paginator.go -pkg mock . Paginator
//go:generate moq -out mock/datastore.go -pkg mock . DataStore
//go:generate moq -out mock/healthcheck.go -pkg mock . HealthChecker
//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/initaliser.go -pkg mock . Initialiser

// Paginator defines the required methods from the paginator package
type Paginator interface {
	GetPaginationValues(r *http.Request) (offset int, limit int, err error)
}

// DataStore implements the methods required to interact with the database
type DataStore interface {
	Init(config.MongoConfig) (err error)
	Close(ctx context.Context) (err error)
	AddBook(ctx context.Context, book *models.Book) (err error)
	GetBook(ctx context.Context, id string) (*models.Book, error)
	GetBooks(ctx context.Context, offset, limit int) ([]models.Book, int, error)
	GetReview(ctx context.Context, reviewID string) (*models.Review, error)
	GetReviews(ctx context.Context, bookID string, offset, limit int) ([]models.Review, int, error)
	AddReview(ctx context.Context, review *models.Review) (err error)
	UpdateReview(ctx context.Context, reviewID string, review *models.Review) (err error)
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
