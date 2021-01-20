package api

import (
	"github.com/cadmiumcat/books-api/interfaces"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	Router    *mux.Router
	dataStore interfaces.DataStore
}

func Setup(host string, router *mux.Router, dataStore interfaces.DataStore) *API {
	api := &API{
		Router:    router,
		dataStore: dataStore,
	}

	api.Router.HandleFunc("/books", api.createBook).Methods("POST")
	api.Router.HandleFunc("/books", listBooks).Methods("GET")
	api.Router.HandleFunc("/books/{id}", getBook).Methods("GET")

	api.Router.HandleFunc("/library/{id}/checkout", checkoutBook).Methods("PUT")
	api.Router.HandleFunc("/library/{id}/checkin", checkinBook).Methods("PUT")

	http.ListenAndServe(host, api.Router)

	return api

}

func setupRoutes(api *API) {

	return
}
