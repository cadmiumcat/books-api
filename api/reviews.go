package api

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) addReviewHandler(writer http.ResponseWriter, request *http.Request)  {
	ctx := request.Context()
	bookID := mux.Vars(request)["id"]

	logData := log.Data{"book_id": bookID}

	if bookID == "" {
		handleError(ctx, writer, apierrors.ErrEmptyBookID, logData)
		return
	}

	// Confirm that book exists. If bookID not found, then a review cannot be added!
	_, err := api.dataStore.GetBook(ctx, bookID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

}

func (api *API) getReviewsHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	bookID := mux.Vars(request)["id"]

	logData := log.Data{"book_id": bookID}

	if bookID == "" {
		handleError(ctx, writer, apierrors.ErrEmptyBookID, logData)
		return
	}

	// Confirm that book exists. If bookID not found, then do not check for the reviews
	_, err := api.dataStore.GetBook(ctx, bookID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	reviews, err := api.dataStore.GetReviews(ctx, bookID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	reviews.Count = len(reviews.Items)

	if err := WriteJSONBody(reviews, writer, http.StatusOK); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	log.Event(ctx, "successfully retrieved review", log.INFO, logData)
}

func (api *API) getReviewHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	bookID := mux.Vars(request)["id"]
	reviewID := mux.Vars(request)["reviewID"]

	logData := log.Data{"book_id": bookID, "review_id": reviewID}

	if bookID == "" {
		handleError(ctx, writer, apierrors.ErrEmptyBookID, logData)
		return
	}

	if reviewID == "" {
		handleError(ctx, writer, apierrors.ErrEmptyReviewID, logData)
		return
	}

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
