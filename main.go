package main

import (
	"encoding/json"
	"github.com/eldeal/skills/config"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

func main() {
	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Event(nil, "error retrieving the configuration", log.FATAL, log.Error(err))
	}

	r := mux.NewRouter()

	r.HandleFunc("/library", createBook).Methods("POST")
	r.HandleFunc("/library", listBooks).Methods("GET")
	r.HandleFunc("/library/{id}", getBook).Methods("GET")

	r.HandleFunc("/library/{id}/checkout", checkoutBook).Methods("POST")
	r.HandleFunc("/library/{id}/checkout", checkinBook).Methods("PUT")

	http.ListenAndServe(cfg.BindAddr, r)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		readFailed(w, err)
		return
	}

	var book Book
	err = json.Unmarshal(b, &book)
	if err != nil {
		unmarshalFailed(w, err)
		return
	}

	add(book)

	w.Header().Set("content-type", "application/json")
	w.Write(b)
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(getAll())
	if err != nil {
		marshalFailed(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	book := get(id)
	if book == nil {
		bookNotFound(w)
		return
	}

	b, err := json.Marshal(book)
	if err != nil {
		marshalFailed(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func checkoutBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	book := get(mux.Vars(r)["id"])
	if book == nil {
		bookNotFound(w)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		readFailed(w, err)
		return
	}

	var co Checkout
	err = json.Unmarshal(b, &co)
	if err != nil {
		unmarshalFailed(w, err)
		return
	}

	if err := checkout(book, co.Who); err != nil {
		log.Event(ctx, "could not check out book", log.ERROR, log.Error(err))
		http.Error(w, "invalid checkout details provided", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return
}

func checkinBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		readFailed(w, err)
		return
	}

	var co Checkout
	err = json.Unmarshal(b, &co)
	if err != nil {
		unmarshalFailed(w, err)
		return
	}

	book := get(mux.Vars(r)["id"])
	if book == nil {
		bookNotFound(w)
		return
	}

	if err := checkin(book, co.Review); err != nil {
		log.Event(ctx, "could not check in book", log.ERROR, log.Error(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return
}
