package api

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/mongo"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func Setup(cfg *config.Configuration) {
	router := setupRoutes()

	mongodb := &mongo.Mongo{}
	err := mongodb.Init(cfg.MongoConfig)
	if err != nil {
		log.Event(nil, "failed to initialise mongo", log.ERROR, log.Error(err))
		os.Exit(1)
	}

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
