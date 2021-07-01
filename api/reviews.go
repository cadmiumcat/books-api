package api

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/models"
	"github.com/cadmiumcat/books-api/pagination"
	"github.com/gorilla/mux"
	"net/http"
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

	review := models.NewReview(bookID)

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

	api.dataStore.AddReview(ctx, review)

	if err := WriteJSONBody(review, writer, http.StatusCreated); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

}

func (api *API) getReviewsHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	bookID := mux.Vars(request)["id"]

	logData := log.Data{"book_id": bookID}

	offset, limit, err := api.paginator.GetPaginationValues(request)
	logData["offset"] = offset
	logData["limit"] = limit
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	if bookID == "" {
		handleError(ctx, writer, apierrors.ErrEmptyBookID, logData)
		return
	}

	// Confirm that book exists. If bookID not found, then do not check for the reviews
	_, err = api.dataStore.GetBook(ctx, bookID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	reviews, totalCount, err := api.dataStore.GetReviews(ctx, bookID, offset, limit)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	response := models.ReviewsResponse{
		Items: reviews,
		Page: pagination.Page{
			Count:      len(reviews),
			Offset:     offset,
			Limit:      limit,
			TotalCount: totalCount,
		},
	}

	if err := WriteJSONBody(response, writer, http.StatusOK); err != nil {
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

	// Confirm that book exists. If bookID not found, or there's another error, then return
	_, err := api.dataStore.GetBook(ctx, bookID)
	if err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	// Confirm that the review exists. If reviewID not found, or there's another error, then return
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

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

}
