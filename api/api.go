package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/interfaces"
	"github.com/cadmiumcat/books-api/mongo"
	"github.com/cadmiumcat/books-api/pagination"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

type API struct {
	host      string
	router    *mux.Router
	paginator interfaces.Paginator
	dataStore interfaces.DataStore
	hc        interfaces.HealthChecker
}

// Setup sets up the endpoints.
func Setup(ctx context.Context, host string, router *mux.Router, paginator interfaces.Paginator, dataStore interfaces.DataStore, hc interfaces.HealthChecker) *API {
	api := &API{
		host:      host,
		router:    router,
		paginator: paginator,
		dataStore: dataStore,
		hc:        hc,
	}

	// Endpoints
	api.router.HandleFunc("/books", api.addBookHandler).Methods("POST")
	api.router.HandleFunc("/books", api.getBooksHandler).Methods("GET")
	api.router.HandleFunc("/books/{id}", api.getBookHandler).Methods("GET")

	api.router.HandleFunc("/books/{id}/reviews", api.getReviewsHandler).Methods("GET")
	api.router.HandleFunc("/books/{id}/reviews", api.addReviewHandler).Methods("POST")
	api.router.HandleFunc("/books/{id}/reviews/{reviewID}", api.getReviewHandler).Methods("GET")
	api.router.HandleFunc("/books/{id}/reviews/{reviewID}", api.updateReviewHandler).Methods("PUT")

	api.router.HandleFunc("/health", api.hc.Handler).Methods("GET")

	log.Event(ctx, "enabling endpoints", log.INFO, log.Data{"bind_addr": api.host})

	return api

}

// WriteJSONBody marshals the provided interface into json, and writes it to the response body.
func WriteJSONBody(v interface{}, w http.ResponseWriter, httpStatus int) error {

	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)

	// Marshal provided model
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// Write payload to body
	if _, err := w.Write(payload); err != nil {
		return err
	}
	return nil
}

// ReadJSONBody reads the bytes from the provided body, and marshals it to the provided model interface.
func ReadJSONBody(ctx context.Context, body io.ReadCloser, v interface{}) error {
	defer body.Close()

	// Get Body bytes
	payload, err := ioutil.ReadAll(body)
	if err != nil {
		return apierrors.ErrUnableToReadMessage
	}

	// Unmarshal body bytes to model
	if err := json.Unmarshal(payload, v); err != nil {
		return apierrors.ErrUnableToParseJSON
	}

	return nil
}

func handleError(ctx context.Context, w http.ResponseWriter, err error, data log.Data) {
	var status int
	var apiError error

	if err != nil {
		apiError = err
		switch err {
		case mongo.ErrBookNotFound,
			mongo.ErrReviewNotFound:
			status = http.StatusNotFound
		case apierrors.ErrRequiredFieldMissing,
			apierrors.ErrEmptyRequestBody,
			apierrors.ErrEmptyBookID,
			apierrors.ErrEmptyReviewID,
			apierrors.ErrInvalidReview,
			apierrors.ErrEmptyReviewMessage,
			apierrors.ErrEmptyReviewUser,
			apierrors.ErrLongReviewMessage,
			pagination.ErrInvalidLimitParameter,
			pagination.ErrInvalidOffsetParameter,
			pagination.ErrLimitOverMax:
			status = http.StatusBadRequest
		default:
			apiError = apierrors.ErrInternalServer
			status = http.StatusInternalServerError
		}
	}

	if data == nil {
		data = log.Data{}
	}

	data["response_status"] = status
	log.Event(ctx, "request unsuccessful", log.ERROR, log.Error(err), data)

	http.Error(w, apiError.Error(), status)
}
