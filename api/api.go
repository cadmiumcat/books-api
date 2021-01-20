package api

import (
	"github.com/cadmiumcat/books-api/interfaces"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	host string
	router    *mux.Router
	dataStore interfaces.DataStore
}

func Setup(host string, router *mux.Router, dataStore interfaces.DataStore) *API {
	api := &API{
		host: host,
		router:    router,
		dataStore: dataStore,
	}

	api.router.HandleFunc("/books", api.createBook).Methods("POST")
	api.router.HandleFunc("/books", listBooks).Methods("GET")
	api.router.HandleFunc("/books/{id}", getBook).Methods("GET")

	api.router.HandleFunc("/library/{id}/checkout", checkoutBook).Methods("PUT")
	api.router.HandleFunc("/library/{id}/checkin", checkinBook).Methods("PUT")

	http.ListenAndServe(api.host, api.router)

	return api

}
