package api

import (
	"context"
	"fmt"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/interfaces/mock"
	"github.com/cadmiumcat/books-api/mongo"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func hasRoute(t *testing.T, r *mux.Router, path, method string) bool {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		api := Setup(ctx, "", r, &mock.PaginatorMock{}, &mock.DataStoreMock{}, &mock.HealthCheckerMock{})

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(t, api.router, "/books", "GET"), ShouldBeTrue)
			So(hasRoute(t, api.router, "/books", "POST"), ShouldBeTrue)
			So(hasRoute(t, api.router, "/books/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(t, api.router, "/books/{id}/reviews", "GET"), ShouldBeTrue)
			So(hasRoute(t, api.router, "/books/{id}/reviews", "POST"), ShouldBeTrue)
			So(hasRoute(t, api.router, "/books/{id}/reviews/{review_id}", "GET"), ShouldBeTrue)
		})
	})
}

func TestHandleError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		input       error
		expected    int
	}{
		{
			input:    mongo.ErrBookNotFound,
			expected: http.StatusNotFound,
		},
		{
			input:    mongo.ErrReviewNotFound,
			expected: http.StatusNotFound,
		},
		{
			input:    apierrors.ErrRequiredFieldMissing,
			expected: http.StatusBadRequest,
		},
		{
			input:    apierrors.ErrEmptyRequestBody,
			expected: http.StatusBadRequest,
		},
		{
			input:    apierrors.ErrEmptyBookID,
			expected: http.StatusBadRequest,
		},
		{
			input:    apierrors.ErrEmptyReviewID,
			expected: http.StatusBadRequest,
		},
		{
			description: "unknown error",
			input:       errMongoDB,
			expected:    http.StatusInternalServerError,
		},
	}

	Convey("Given a specific error", t, func() {
		for _, test := range cases {
			ctx := context.Background()
			err := test.input
			Convey("When the "+test.input.Error()+" error is passed to the handleError function", func() {
				writer := httptest.NewRecorder()
				handleError(ctx, writer, err, nil)

				Convey(fmt.Sprintf("Then the status returned is %v", test.expected), func() {
					So(writer.Code, ShouldEqual, test.expected)

				})
			})
		}
	})

}

type errReader int

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("test error")
}

func TestReadJSONBody(t *testing.T) {
	type fakeBook struct {
		Title string
	}

	Convey("Given a request body that can be unmarshalled as JSON", t, func() {
		request := httptest.NewRequest(http.MethodPost, "/something", strings.NewReader(`{"Title":"fakeBook"}`))

		fakeLibrary := &fakeBook{}
		Convey("When the ReadJSONBody function is called", func() {
			err := ReadJSONBody(nil, request.Body, fakeLibrary)
			Convey("Then the body is successfully unmarshalled", func() {
				So(err, ShouldBeNil)
				So(fakeLibrary.Title, ShouldEqual, "fakeBook")
			})
		})
	})

	Convey("Given a request body with an error", t, func() {
		request := httptest.NewRequest(http.MethodPost, "/something", errReader(0))

		Convey("When the ReadJSONBody function is called", func() {
			err := ReadJSONBody(nil, request.Body, nil)
			Convey("Then an error is returned saying it was unable to read the message", func() {
				So(err, ShouldBeError, apierrors.ErrUnableToReadMessage)
			})
		})
	})

	Convey("Given a request with a body that cannot be unmarshalled as JSON", t, func() {
		request := httptest.NewRequest(http.MethodPost, "/something", strings.NewReader(`"Title":"fakeBook"`))

		Convey("When the ReadJSONBody function is called", func() {
			err := ReadJSONBody(nil, request.Body, &fakeBook{})
			Convey("Then an error is returned saying it was unable to read the message", func() {
				So(err, ShouldBeError, apierrors.ErrUnableToParseJSON)
			})
		})
	})
}

func TestWriteJSONBody(t *testing.T) {
	Convey("Given an interface that cannot be marshalled into JSON", t, func() {
		badPayload := make(chan int)

		Convey("When the WriteJSONBody function is called", func() {
			err := WriteJSONBody(badPayload, httptest.NewRecorder(), http.StatusOK)

			Convey("An error is returned", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
