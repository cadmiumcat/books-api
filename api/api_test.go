package api

import (
	"context"
	"fmt"
	"github.com/cadmiumcat/books-api/apierrors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleError(t *testing.T) {
	t.Parallel()

	cases := []struct{
		description string
		input error
		expected int
	}{
		{
			input: apierrors.ErrBookNotFound,
			expected: http.StatusNotFound,
		},
		{
			input: apierrors.ErrReviewNotFound,
			expected: http.StatusNotFound,
		},
		{
			input: apierrors.ErrRequiredFieldMissing,
			expected: http.StatusBadRequest,
		},
		{
			input: apierrors.ErrEmptyRequest,
			expected: http.StatusBadRequest,
		},
		{
			description: "unknown error",
			input: errMongoDB,
			expected: http.StatusInternalServerError,
		},

	}

	Convey("Given a specific error", t, func() {
		for _, test := range cases {
			ctx := context.Background()
			err := test.input
			Convey("When I pass the " + test.input.Error() + " error to the handleError function", func() {
				writer := httptest.NewRecorder()
				handleError(ctx, writer, err, nil)

				Convey(fmt.Sprintf("Then the status returned is %v",  test.expected), func() {
					So(writer.Code, ShouldEqual, test.expected)

				})
			})
		}
	})

}