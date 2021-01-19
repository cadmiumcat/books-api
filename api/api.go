package api

import (
	"github.com/cadmiumcat/books-api/config"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	Router    *mux.Router
	dataStore DataStore
}

func Setup(cfg *config.Configuration, dataStore DataStore) {
	api := &API{
		Router:    mux.NewRouter(),
		dataStore: dataStore,
	}

	setupRoutes(api)

	http.ListenAndServe(cfg.BindAddr, api.Router)

}
func setupRoutes(api *API)  {

	api.Router.HandleFunc("/books", api.createBook).Methods("POST")
	api.Router.HandleFunc("/books", listBooks).Methods("GET")
	api.Router.HandleFunc("/books/{id}", getBook).Methods("GET")

	api.Router.HandleFunc("/library/{id}/checkout", checkoutBook).Methods("PUT")
	api.Router.HandleFunc("/library/{id}/checkin", checkinBook).Methods("PUT")
	return
}
