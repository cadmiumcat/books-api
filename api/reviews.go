package api

import (
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

func (api *API) addReviewHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	bookID := mux.Vars(request)["id"]

	logData := log.Data{"book_id": bookID}

	if bookID == "" {
		handleError(ctx, writer, apierrors.ErrEmptyBookID, logData)
		return
	}

	if request.ContentLength == 0 {
		handleError(ctx, writer, apierrors.ErrEmptyRequestBody, logData)
		return
	}

	// Confirm that book exists. If bookID not found, then a review cannot be added!
	_, err := api.dataStore.GetBook(ctx, bookID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	review := &models.Review{
		User:  models.User{},
		Links: &models.ReviewLink{}}
	if err := ReadJSONBody(ctx, request.Body, review); err != nil {
		handleError(ctx, writer, apierrors.ErrInvalidReview, logData)
		return
	}

	logData["review"] = review

	err = review.Validate()
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	review.ID = uuid.NewV4().String()
	review.BookID = bookID
	review.Links.Self = fmt.Sprintf("/books/%s/reviews/%s", bookID, review.ID)
	review.Links.Book = fmt.Sprintf("/books/%s", bookID)
	review.LastUpdated = time.Now().UTC()

	api.dataStore.AddReview(review)

	if err := WriteJSONBody(review, writer, http.StatusCreated); err != nil {
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

func (api *API) updateReviewHandler(writer http.ResponseWriter, request *http.Request) {
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

	// Confirm that the review exists
	_, err = api.dataStore.GetReview(ctx, reviewID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	review := &models.Review{User: models.User{}}
	if err := ReadJSONBody(ctx, request.Body, review); err != nil {
		handleError(ctx, writer, apierrors.ErrInvalidReview, logData)
		return
	}

	logData["review"] = review

	err = api.dataStore.UpdateReview(ctx, reviewID, review)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	review, err = api.dataStore.GetReview(ctx, reviewID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	if err := WriteJSONBody(review, writer, http.StatusOK); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}
}
