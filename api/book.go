package api

import (
	"encoding/json"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

func checkout(b *models.Book, name string) error {
	h := len(b.History)
	if h != 0 {
		lastCheckout := b.History[h-1]
		if lastCheckout.In.IsZero() {
			return ErrBookCheckedOut
		}
	}

	if len(name) == 0 {
		return ErrNameMissing
	}

	b.History = append(b.History, models.Checkout{
		Who: name,
		Out: time.Now(),
	})

	return nil
}

func checkin(b *models.Book, review int) error {
	h := len(b.History)
	if h == 0 {
		return ErrBookNotCheckedOut
	}

	if review < 1 || review > 5 {
		return ErrReviewMissing
	}

	lastCheckout := b.History[h-1]
	if !lastCheckout.In.IsZero() {
		return ErrBookNotCheckedOut
	}

	b.History[h-1] = models.Checkout{
		Who:    lastCheckout.Who,
		Out:    lastCheckout.Out,
		In:     time.Now(),
		Review: review,
	}

	return nil
}

func (api *API) createBook(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	bytes, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		readFailed(ctx, writer, err)
		return
	}

	book := &models.Book{}

	err = json.Unmarshal(bytes, &book)
	if err != nil {
		unmarshalFailed(ctx, writer, err)
		return
	}

	err = book.Validate()
	if err != nil {
		invalidBook(ctx, writer, err)
		return
	}

	api.dataStore.AddBook(book)

	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(bytes)
}

func (api *API) listBooks(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	books, err := api.dataStore.GetBooks()

	books.Count = len(books.Items)

	bytes, err := json.Marshal(books)
	if err != nil {
		marshalFailed(ctx, writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(bytes)
}

func (api *API) getBook(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	id := mux.Vars(request)["id"]

	book, err := api.dataStore.GetBook(id)
	if book == nil {
		bookNotFound(ctx, writer, id)
		return
	}

	bytes, err := json.Marshal(book)
	if err != nil {
		marshalFailed(ctx, writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(bytes)
}
