package api

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

const emptyJson = "{}"

func checkout(b *models.Book, name string) error {
	h := len(b.History)
	if h != 0 {
		lastCheckout := b.History[h-1]
		if lastCheckout.In.IsZero() {
			return apierrors.ErrBookCheckedOut
		}
	}

	if len(name) == 0 {
		return apierrors.ErrNameMissing
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
		return apierrors.ErrBookNotCheckedOut
	}

	if review < 1 || review > 5 {
		return apierrors.ErrReviewMissing
	}

	lastCheckout := b.History[h-1]
	if !lastCheckout.In.IsZero() {
		return apierrors.ErrBookNotCheckedOut
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
		handleError(ctx, writer, apierrors.ErrEmptyRequest, nil)
		return
	}

	book := &models.Book{}
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
	api.dataStore.AddBook(book)

	if err := WriteJSONBody(book, writer, http.StatusCreated); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}
}

func (api *API) listBooks(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	books, err := api.dataStore.GetBooks(ctx)
	if err != nil {
		handleError(ctx, writer, err, nil)
	}

	books.Count = len(books.Items)

	if err := WriteJSONBody(books, writer, http.StatusOK); err != nil {
		handleError(ctx, writer, err, nil)
		return
	}
	log.Event(ctx, "successfully retrieved list of books", log.INFO)
}

func (api *API) getBook(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	id := mux.Vars(request)["id"]

	logData := log.Data{"book_id": id}

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
