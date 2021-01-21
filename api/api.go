package api

import (
	"context"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/interfaces"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	host      string
	router    *mux.Router
	dataStore interfaces.DataStore
}

func Setup(ctx context.Context, host string, router *mux.Router, dataStore interfaces.DataStore) *API {
	api := &API{
		host:      host,
		router:    router,
		dataStore: dataStore,
	}

	api.router.HandleFunc("/books", api.createBook).Methods("POST")
	api.router.HandleFunc("/books", api.listBooks).Methods("GET")
	api.router.HandleFunc("/books/{id}", api.getBook).Methods("GET")

	log.Event(ctx, "starting http server", log.INFO, log.Data{"bind_addr": api.host})
	http.ListenAndServe(api.host, api.router)

	return api

}
