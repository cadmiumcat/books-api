package api

import (
	"bytes"
	"encoding/json"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"time"
)

const emptyJson = "{}"

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

	if request.Body == http.NoBody {
		missingBody(ctx, writer, ErrRequestBodyMissing)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		readFailed(ctx, writer, err)
		return
	}
	buffer := new(bytes.Buffer)
	json.Compact(buffer, body)

	if buffer.String() == emptyJson {
		emptyRequest(ctx, writer, ErrEmptyRequest)
		return
	}

	book := &models.Book{}

	err = json.Unmarshal(body, &book)
	if err != nil {
		unmarshalFailed(ctx, writer, err)
		return
	}

	err = book.Validate()
	if err != nil {
		invalidBook(ctx, writer, err)
		return
	}

	book.ID = uuid.NewV4().String()
	api.dataStore.AddBook(book)

	body, err = json.Marshal(book)
	if err != nil {
		marshalFailed(ctx, writer, err)
		return
	}

	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write(body)
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
	_, _ = writer.Write(bytes)
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
	_, _ = writer.Write(bytes)
}
