package api

import (
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) getReviews(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	bookID := mux.Vars(request)["id"]

	reviews, err := api.dataStore.GetReviews(ctx, bookID)

	reviews.Count = len(reviews.Items)

	bytes, err := json.Marshal(reviews)
	if err != nil {
		marshalFailed(ctx, writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(bytes)
}

func (api *API) getReview(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	reviewID := mux.Vars(request)["reviewID"]

	review, _ := api.dataStore.GetReview(ctx, reviewID)
	if review == nil {
		msg := fmt.Sprintf("review id %q not found", reviewID)
		log.Event(ctx, msg, log.INFO)
		http.Error(writer, msg, http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(review)
	if err != nil {
		marshalFailed(ctx, writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(bytes)
}
