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

func (api *API) createBook(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		readFailed(w, err)
		return
	}

	book := &models.Book{}

	err = json.Unmarshal(b, &book)
	if err != nil {
		unmarshalFailed(w, err)
		return
	}

	err = book.Validate()
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.dataStore.AddBook(book)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func (api *API) listBooks(w http.ResponseWriter, r *http.Request) {
	books, err := api.dataStore.GetBooks()

	books.Count = len(books.Items)

	b, err := json.Marshal(books)
	if err != nil {
		marshalFailed(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (api *API) getBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	book, err := api.dataStore.GetBook(id)
	if book == nil {
		bookNotFound(w, id)
		return
	}

	b, err := json.Marshal(book)
	if err != nil {
		marshalFailed(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

