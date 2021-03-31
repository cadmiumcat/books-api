package api

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) getReviews(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	bookID := mux.Vars(request)["id"]

	logData := log.Data{"book_id": bookID}

	reviews, err := api.dataStore.GetReviews(ctx, bookID)
	if err != nil {
		handleError(ctx, writer, err, logData)
	}

	reviews.Count = len(reviews.Items)

	if err := WriteJSONBody(reviews, writer, http.StatusOK); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	log.Event(ctx, "successfully retrieved review", log.INFO, logData)
}

func (api *API) getReview(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	bookID := mux.Vars(request)["id"]
	reviewID := mux.Vars(request)["reviewID"]

	logData := log.Data{"book_id": bookID, "review_id": reviewID}

	// Confirm that book exists. If bookID not found, then do not check for the review
	_, err := api.dataStore.GetBook(ctx, bookID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	review, err := api.dataStore.GetReview(ctx, reviewID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	if err := WriteJSONBody(review, writer, http.StatusOK); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	log.Event(ctx, "successfully retrieved review", log.INFO, logData)
}
