package api

import (
	"github.com/cadmiumcat/books-api/interfaces"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	host      string
	router    *mux.Router
	dataStore interfaces.DataStore
}

func Setup(host string, router *mux.Router, dataStore interfaces.DataStore) *API {
	api := &API{
		host:      host,
		router:    router,
		dataStore: dataStore,
	}

	api.router.HandleFunc("/books", api.createBook).Methods("POST")
	api.router.HandleFunc("/books", api.listBooks).Methods("GET")
	api.router.HandleFunc("/books/{id}", api.getBook).Methods("GET")

	http.ListenAndServe(api.host, api.router)

	return api

}
