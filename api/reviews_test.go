package api

import (
	"context"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/interfaces/mock"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	bookID   = "1"
	reviewID = "123"
	emptyID  = ""
)

var errMongoDB = errors.New("unexpected error in MongoDB")

func TestReviewEndpoints(t *testing.T) {
	hcMock := mock.HealthCheckerMock{}

	Convey("Given an existing book with at least one review (review_id=123)", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: bookID}, nil
			},
			GetReviewFunc: func(ctx context.Context, id string) (*models.Review, error) {
				return &models.Review{ID: reviewID}, nil
			},
		}

		ctx := context.Background()

		api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
		Convey("When I send an HTTP GET request to /books/1/reviews/123", func() {
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID+"/reviews/"+reviewID, nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("Then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
			Convey("And the GetReview function is called once", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 1)
				So(response.Body.String(), ShouldEqual, `{"id":"123"}`)
			})
		})
	})

	Convey("Given an existing book with no reviews", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: bookID}, nil
			},
			GetReviewFunc: func(ctx context.Context, id string) (*models.Review, error) {
				return nil, apierrors.ErrReviewNotFound
			},
		}

		ctx := context.Background()

		api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
		Convey("When I send an HTTP GET request to /books/1/reviews/123", func() {
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID+"/reviews/"+reviewID, nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("Then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
			})
			Convey("And the GetReview function is called once", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 1)
				So(response.Body.String(), ShouldEqual, "review not found\n")
			})
		})
	})

	Convey("Given a GET request for a review of a book", t, func() {
		Convey("When the book does not exist", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, apierrors.ErrBookNotFound
				},
			}

			ctx := context.Background()

			api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID+"/reviews/"+reviewID, nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("Then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
			})
			Convey("And the GetReview function is not called", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 0)
				So(response.Body.String(), ShouldEqual, "book not found\n")
			})
		})

		Convey("When the database returns an unexpected error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, errMongoDB
				},
			}

			ctx := context.Background()

			api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID+"/reviews/"+reviewID, nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("Then 500 InternalServerError status code is returned", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
				So(response.Body.String(), ShouldEqual, "unexpected error in MongoDB\n")
			})
		})
	})

	Convey("Given an HTTP GET request to the /books/{id}/reviews/{review_id} endpoint", t, func() {

		api := &API{
			host:      host,
			router:    mux.NewRouter(),
			dataStore: &mock.DataStoreMock{},
			hc:        &hcMock,
		}

		Convey("When the {review_id} is empty", func() {
			request, err := http.NewRequest("GET", "/books/"+bookID+"/reviews/"+emptyID, nil)
			So(err, ShouldBeNil)

			response := httptest.NewRecorder()

			api.getReview(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, "empty review ID in request\n")
			})
		})
	})

	Convey("Given an HTTP GET request to the /books/{id}/reviews/{review_id} endpoint", t, func() {

		api := &API{
			host:      host,
			router:    mux.NewRouter(),
			dataStore: &mock.DataStoreMock{},
			hc:        &hcMock,
		}

		Convey("When the {id} is empty", func() {
			request, err := http.NewRequest("GET", "/books/"+emptyID+"/reviews/"+reviewID, nil)
			So(err, ShouldBeNil)

			response := httptest.NewRecorder()

			api.getReview(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, "empty review ID in request\n")
			})
		})
	})

}
