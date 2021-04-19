package api

import (
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (api *API) addBookHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	if request.ContentLength == 0 {
		handleError(ctx, writer, apierrors.ErrEmptyRequestBody, nil)
		return
	}

	book := &models.Book{Links: &models.Link{}}
	if err := ReadJSONBody(ctx, request.Body, book); err != nil {
		handleError(ctx, writer, err, nil)
		return
	}

	logData := log.Data{"book": book}

	err := book.Validate()
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	book.ID = uuid.NewV4().String()
	book.Links.Self = fmt.Sprintf("/books/%s", book.ID)
	book.Links.Reviews = fmt.Sprintf("/books/%s/reviews", book.ID)

	api.dataStore.AddBook(ctx, book)

	if err := WriteJSONBody(book, writer, http.StatusCreated); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}
}

func (api *API) getBooksHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	books, err := api.dataStore.GetBooks(ctx)
	if err != nil {
		handleError(ctx, writer, err, nil)
		return
	}

	books.Count = len(books.Items)

	if err := WriteJSONBody(books, writer, http.StatusOK); err != nil {
		handleError(ctx, writer, err, nil)
		return
	}
	log.Event(ctx, "successfully retrieved list of books", log.INFO)
}

func (api *API) getBookHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	id := mux.Vars(request)["id"]
	logData := log.Data{"book_id": id}

	if id == "" {
		handleError(ctx, writer, apierrors.ErrEmptyBookID, logData)
		return
	}

	book, err := api.dataStore.GetBook(ctx, id)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	if err := WriteJSONBody(book, writer, http.StatusOK); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}
	log.Event(ctx, "successfully retrieved book", log.INFO, logData)
}
