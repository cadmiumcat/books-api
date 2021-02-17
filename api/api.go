package api

import (
	"context"
	dpHttp "github.com/ONSdigital/dp-net/http"
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

// Setup sets up the endpoints and starts the http  server.
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

	api.router.HandleFunc("/health", api.hc.Handler).Methods("GET")

	log.Event(ctx, "starting http server", log.INFO, log.Data{"bind_addr": api.host})

	httpServer := dpHttp.NewServer(api.host, api.router)
	httpServer.ListenAndServe()

	return api

}
