package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) getReview(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	reviewID := mux.Vars(request)["reviewID"]

	review, _ := api.dataStore.GetReview(reviewID)

	bytes, err := json.Marshal(review)
	if err != nil {
		marshalFailed(ctx, writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(bytes)
}
