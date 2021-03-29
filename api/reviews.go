package api

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) getReview(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	bookID := mux.Vars(request)["id"]
	reviewID := mux.Vars(request)["reviewID"]

	logData := log.Data{"book_id": bookID, "review_id": reviewID}

	review, err := api.dataStore.GetReview(ctx, reviewID)
	if review == nil {
		handleError(ctx, writer, err, logData)
	}

	if err := WriteJSONBody(review, writer, http.StatusOK); err != nil {
		handleError(ctx, writer, err, logData)
		return
	}

	log.Event(ctx, "successfully retrieved review", log.INFO, logData)
}
