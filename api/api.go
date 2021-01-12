package api

import (
	"github.com/cadmiumcat/books-api/config"
	"github.com/gorilla/mux"
	"net/http"
)

func Setup(cfg *config.Configuration) {
	router := setupRoutes()

	http.ListenAndServe(cfg.BindAddr, router)

}
func setupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books", listBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")

	router.HandleFunc("/library/{id}/checkout", checkoutBook).Methods("PUT")
	router.HandleFunc("/library/{id}/checkin", checkinBook).Methods("PUT")
	return router
}
