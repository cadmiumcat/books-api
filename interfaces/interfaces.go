package interfaces

import (
	"context"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/models"
)

//go:generate moq -out datastoretest/datastore.go -pkg datastoretest . DataStore

type DataStore interface {
	Init(config.MongoConfig) (err error)
	Close(ctx context.Context) (err error)
	AddBook(book *models.Book)
	GetBook(id string) (*models.Book, error)
	GetBooks() (models.Books, error)
}
