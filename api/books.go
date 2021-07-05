package api

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/models"
	"github.com/cadmiumcat/books-api/pagination"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) addBookHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	if request.ContentLength == 0 {
		handleError(ctx, writer, apierrors.ErrEmptyRequestBody, nil)
		return
	}

	book := models.NewBook()
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

	api.dataStore.AddBook(ctx, book)

	if err := WriteJSONBody(book, writer, http.StatusCreated); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}
}

func (api *API) getBooksHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logData := log.Data{}

	offset, limit, err := api.paginator.GetPaginationValues(request)
	logData["offset"] = offset
	logData["limit"] = limit
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	books, totalCount, err := api.dataStore.GetBooks(ctx, offset, limit)
	if err != nil {
		handleError(ctx, writer, err, nil)
		return
	}

	response := models.BooksResponse{
		Items: books,
		Page: pagination.Page{
			Count:      len(books),
			Offset:     offset,
			Limit:      limit,
			TotalCount: totalCount,
		},
	}

	if err := WriteJSONBody(response, writer, http.StatusOK); err != nil {
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
