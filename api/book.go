package api

import (
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

var books models.Books

func init() {
	b := models.Book{
		Title:    "default book",
		Author:   "default author",
		Synopsis: "",
	}

	add(b)

}

func get(id string) (book *models.Book) {
	for i, l := range books.Items {
		if l.Id == id {
			book = &books.Items[i]
			break
		}
	}
	return
}

func getAll() models.Books {
	return books
}

func add(b models.Book) {
	count := len(books.Items)
	books.Count = count + 1

	b.Id = fmt.Sprint(books.Count)
	books.Items = append(books.Items, b)
}

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

func (api *BooksAPI) createBook(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		readFailed(w, err)
		return
	}

	var book models.Book
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

	add(book)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(getAll())
	if err != nil {
		marshalFailed(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	book := get(id)
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

func checkoutBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	book := get(id)
	if book == nil {
		bookNotFound(w, id)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		readFailed(w, err)
		return
	}

	var co models.Checkout
	err = json.Unmarshal(b, &co)
	if err != nil {
		unmarshalFailed(w, err)
		return
	}

	if err := checkout(book, co.Who); err != nil {
		log.Event(ctx, "could not check out book", log.ERROR, log.Error(err), log.Data{"book": book.History})
		http.Error(w, "invalid checkout details provided", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return
}

func checkinBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		readFailed(w, err)
		return
	}

	var co models.Checkout
	err = json.Unmarshal(b, &co)
	if err != nil {
		unmarshalFailed(w, err)
		return
	}

	id := mux.Vars(r)["id"]
	book := get(id)
	if book == nil {
		bookNotFound(w, id)
		return
	}

	if err := checkin(book, co.Review); err != nil {
		log.Event(ctx, "could not check in book", log.ERROR, log.Error(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return
}
