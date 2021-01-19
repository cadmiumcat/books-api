package api

import (
	"context"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/models"
)

type DataStore interface {
	Init(config.MongoConfig) (err error)
	Close(ctx context.Context) (err error)
	AddBook(book *models.Book)
}
