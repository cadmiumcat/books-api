package api

import (
	"context"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleError(t *testing.T) {
	Convey("Given an unknown error", t, func() {
		ctx := context.Background()
		err := errors.New("Unknown error")

		Convey("When I pass the error to the handleError function", func() {
			writer := httptest.NewRecorder()
			handleError(ctx, writer, err, nil)

			Convey("Then the status returned is internal server error", func() {
				So(writer.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})

}