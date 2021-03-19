package api

import (
	"context"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/interfaces"
	"github.com/gorilla/mux"
)

type API struct {
	host      string
	router    *mux.Router
	dataStore interfaces.DataStore
	hc        interfaces.HealthChecker
}

// Setup sets up the endpoints.
func Setup(ctx context.Context, host string, router *mux.Router, dataStore interfaces.DataStore, hc interfaces.HealthChecker) *API {
	api := &API{
		host:      host,
		router:    router,
		dataStore: dataStore,
		hc:        hc,
	}

	// Endpoints
	api.router.HandleFunc("/books", api.createBook).Methods("POST")
	api.router.HandleFunc("/books", api.listBooks).Methods("GET")
	api.router.HandleFunc("/books/{id}", api.getBook).Methods("GET")

	api.router.HandleFunc("/books/{id}/reviews/{reviewID}", api.getReview).Methods("GET")

	api.router.HandleFunc("/health", api.hc.Handler).Methods("GET")

	log.Event(ctx, "enabling endpoints", log.INFO, log.Data{"bind_addr": api.host})

	return api

}
