package api

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/mongo"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type BooksAPI struct {
	Router *mux.Router
}

func Setup(cfg *config.Configuration) {
	api := &BooksAPI{Router: mux.NewRouter()}

	setupRoutes(api)

	mongodb := &mongo.Mongo{}
	err := mongodb.Init(cfg.MongoConfig)
	if err != nil {
		log.Event(nil, "failed to initialise mongo", log.ERROR, log.Error(err))
		os.Exit(1)
	}

	http.ListenAndServe(cfg.BindAddr, api.Router)

}
func setupRoutes(api *BooksAPI)  {

	api.Router.HandleFunc("/books", createBook).Methods("POST")
	api.Router.HandleFunc("/books", listBooks).Methods("GET")
	api.Router.HandleFunc("/books/{id}", getBook).Methods("GET")

	api.Router.HandleFunc("/library/{id}/checkout", checkoutBook).Methods("PUT")
	api.Router.HandleFunc("/library/{id}/checkin", checkinBook).Methods("PUT")
	return
}
