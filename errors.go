package main

import (
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

func readFailed(w http.ResponseWriter, err error) {
	log.Event(nil, "error reading request body", log.ERROR, log.Error(err))
	http.Error(w, "cannot read request body", http.StatusInternalServerError)
}

func bookNotFound(w http.ResponseWriter) {
	log.Event(nil, "book not found in list", log.INFO)
	http.Error(w, "book not found", http.StatusNotFound)
}

func unmarshalFailed(w http.ResponseWriter, err error) {
	log.Event(nil, "error returned from json unmarshal", log.ERROR, log.Error(err))
	http.Error(w, "cannot unmarshal json body", http.StatusInternalServerError)
}

func marshalFailed(w http.ResponseWriter, err error) {
	log.Event(nil, "error returned from json marshal", log.ERROR, log.Error(err))
	http.Error(w, "cannot marshal content to json", http.StatusInternalServerError)
}
