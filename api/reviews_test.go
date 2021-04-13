package api

import (
	"context"
	"encoding/json"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/interfaces/mock"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	bookID1   = "1"
	reviewID1 = "123"
	reviewID2 = "567"
	emptyID   = ""
)

var bookReview1 = models.Review{
	ID: reviewID1,
	Links: &models.ReviewLink{
		Book: bookID1,
	},
}

var bookReview2 = models.Review{
	ID: reviewID2,
	Links: &models.ReviewLink{
		Book: bookID1,
	},
}

var emptyReviews = models.Reviews{
	Count: 0,
	Items: nil,
}

var errMongoDB = errors.New("unexpected error in MongoDB")

func TestReviewEndpoints(t *testing.T) {
	t.Parallel()
	hcMock := mock.HealthCheckerMock{}

	Convey("Given an existing book with at least one review (review_id=123)", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: bookID1}, nil
			},
			GetReviewFunc: func(ctx context.Context, id string) (*models.Review, error) {
				return &models.Review{ID: reviewID1}, nil
			},
		}

		ctx := context.Background()

		api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
		Convey("When I send an HTTP GET request to /books/1/reviews/123", func() {
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews/"+reviewID1, nil)
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
				return &models.Book{ID: bookID1}, nil
			},
			GetReviewFunc: func(ctx context.Context, id string) (*models.Review, error) {
				return nil, apierrors.ErrReviewNotFound
			},
			GetReviewsFunc: func(ctx context.Context, bookID string) (models.Reviews, error) {
				return emptyReviews, nil
			},
		}

		ctx := context.Background()

		api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
		Convey("When I send an HTTP GET request to /books/1/reviews/123", func() {
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews/"+reviewID1, nil)
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

		Convey("When I send a HTTP GET request to /books/1/reviews", func() {
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("Then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
			Convey("And the GetReviews function is called once", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 1)
			})
			Convey("And the response contains a count of zero and no review items", func() {
				payload, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)
				reviews := models.Reviews{}
				err = json.Unmarshal(payload, &reviews)
				So(reviews.Count, ShouldEqual, 0)
				So(reviews.Items, ShouldBeNil)
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

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews/"+reviewID1, nil)
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

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews/"+reviewID1, nil)
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
			request, err := http.NewRequest("GET", "/books/"+bookID1+"/reviews/"+emptyID, nil)
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
			request, err := http.NewRequest("GET", "/books/"+emptyID+"/reviews/"+reviewID1, nil)
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

func TestReviewsEndpoint(t *testing.T) {
	t.Parallel()
	hcMock := mock.HealthCheckerMock{}

	Convey("Given a book with one or more reviews", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: bookID1}, nil
			},
			GetReviewsFunc: func(ctx context.Context, bookID string) (models.Reviews, error) {
				return models.Reviews{
					Count: 2,
					Items: []models.Review{
						bookReview1,
						bookReview2,
					}}, nil
			},
		}

		ctx := context.Background()

		api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
		Convey("When I send an HTTP GET request to /books/1/reviews", func() {
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("Then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
			Convey("And the GetReviews function is called once", func() {
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 1)
			})
			Convey("And the response body contains the right number of reviews", func() {
				payload, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)
				reviews := models.Reviews{}
				err = json.Unmarshal(payload, &reviews)
				So(reviews.Count, ShouldEqual, 2)
			})
		})
	})

	Convey("Given a GET request for a list of reviews of a book", t, func() {
		Convey("When the book does not exist", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, apierrors.ErrBookNotFound
				},
			}

			ctx := context.Background()

			api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("Then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
			})
			Convey("And the GetReviews function is not called", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 0)
			})
		})

		Convey("When GetReviews returns a database error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return &models.Book{ID: bookID1}, nil
				},
				GetReviewsFunc: func(ctx context.Context, bookID string) (models.Reviews, error) {
					return models.Reviews{}, errMongoDB
				},
			}

			ctx := context.Background()

			api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("Then the HTTP response code is 500", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
			})
			Convey("And the GetBook and GetReviews functions are called", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 1)
			})
		})

		Convey("When GetBook returns a database error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return &models.Book{}, errMongoDB
				},
			}

			ctx := context.Background()

			api := Setup(ctx, host, mux.NewRouter(), mockDataStore, &hcMock)
			response := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("Then the HTTP response code is 500", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
			})
			Convey("And the GetReviews function is not called", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 0)
			})
		})
	})
}
