package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"net/http"
)

var (
	ErrBookCheckedOut    = errors.New("this book is currently checked out")
	ErrNameMissing       = errors.New("a name must be provided for checkout")
	ErrReviewMissing     = errors.New("a review between 1 and 5 must be provided")
	ErrBookNotCheckedOut = errors.New("this book is not currently checked out")
	ErrInvalidBook       = errors.New("invalid book. Missing required field")
	ErrRequestBodyMissing = errors.New("request body missing")
)

func readFailed(ctx context.Context, w http.ResponseWriter, err error) {
	log.Event(ctx, "error reading request body", log.ERROR, log.Error(err))
	http.Error(w, "cannot read request body", http.StatusInternalServerError)
}

func bookNotFound(ctx context.Context, w http.ResponseWriter, id string) {
	msg := fmt.Sprintf("book id %q not found", id)
	log.Event(ctx, msg, log.INFO)
	http.Error(w, msg, http.StatusNotFound)
}

func unmarshalFailed(ctx context.Context, w http.ResponseWriter, err error) {
	log.Event(ctx, "error returned from json unmarshal", log.ERROR, log.Error(err))
	http.Error(w, "cannot unmarshal json body", http.StatusInternalServerError)
}

func marshalFailed(ctx context.Context, w http.ResponseWriter, err error) {
	log.Event(nil, "error returned from json marshal", log.ERROR, log.Error(err))
	http.Error(w, "cannot marshal content to json", http.StatusInternalServerError)
}

func invalidBook(ctx context.Context, w http.ResponseWriter, err error) {
	log.Event(ctx, "invalid book", log.ERROR, log.Error(err))
	http.Error(w, "invalid book", http.StatusBadRequest)
}

func missingBody(ctx context.Context, w http.ResponseWriter, err error)  {
	log.Event(ctx, "invalid book", log.ERROR, log.Error(err))
	http.Error(w, err.Error(), http.StatusBadRequest)
}