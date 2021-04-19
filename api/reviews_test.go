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
	bookID1                   = "1"
	bookID2                   = "2"
	reviewID1                 = "123"
	reviewID2                 = "567"
	emptyID                   = ""
	bookIDNotInStore          = "notInStore"
	internalSeverErrorMessage = "internal server error\n"
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

func TestGetReviewHandler(t *testing.T) {
	t.Parallel()

	Convey("Given an HTTP GET request to the /books/{id}/reviews/{review_id} endpoint", t, func() {

		Convey("When the {review_id} is empty", func() {
			api := &API{}
			request := httptest.NewRequest("GET", "/books/"+bookID1+"/reviews/"+emptyID, nil)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": emptyID,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, "empty review ID in request\n")
			})
		})

		Convey("When the {id} is empty", func() {
			api := &API{}
			request := httptest.NewRequest("GET", "/books/"+emptyID+"/reviews/"+reviewID1, nil)

			expectedUrlVars := map[string]string{
				"id":       emptyID,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, "empty book ID in request\n")
			})
		})
	})

	Convey("Given an existing book with at least one review (review_id=123)", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: bookID1}, nil
			},
			GetReviewFunc: func(ctx context.Context, id string) (*models.Review, error) {
				return &models.Review{ID: reviewID1}, nil
			},
		}

		api := &API{
			dataStore: mockDataStore,
		}

		Convey("When a http get request is sent to /books/1/reviews/123", func() {
			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews/"+reviewID1, nil)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewHandler(response, request)
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
			GetReviewFunc: func(ctx context.Context, reviewID string) (*models.Review, error) {
				return nil, apierrors.ErrReviewNotFound
			},
		}

		api := &API{
			dataStore: mockDataStore,
		}

		Convey("When a http get request is sent to /books/1/reviews/123", func() {
			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews/"+reviewID1, nil)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewHandler(response, request)
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

	Convey("Given a GET request for a review of a book that does not exist", t, func() {
		Convey("When a http get request is sent to /books/1/reviews/123", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, apierrors.ErrBookNotFound
				},
			}

			api := &API{
				dataStore: mockDataStore,
			}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews/"+reviewID1, nil)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewHandler(response, request)
			Convey("Then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
			})
			Convey("And the GetReview function is not called", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 0)
				So(response.Body.String(), ShouldEqual, "book not found\n")
			})
		})
	})

	Convey("Given a GET request a review of a book", t, func() {
		Convey("When GetBook returns an unexpected database error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, errors.Wrap(errMongoDB, "unexpected error when getting a book")
				},
			}

			api := &API{
				dataStore: mockDataStore,
			}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews/"+reviewID1, nil)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewHandler(response, request)
			Convey("Then 500 InternalServerError status code is returned", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
				So(response.Body.String(), ShouldEqual, internalSeverErrorMessage)
			})
		})

		Convey("When GetReview returns an unexpected database error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return &models.Book{ID: bookID1}, nil
				},
				GetReviewFunc: func(ctx context.Context, reviewID string) (*models.Review, error) {
					return nil, errors.Wrap(errMongoDB, "unexpected error when getting a review")
				},
			}

			api := &API{
				dataStore: mockDataStore,
			}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews/"+reviewID1, nil)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewHandler(response, request)
			Convey("Then 500 InternalServerError status code is returned", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
				So(response.Body.String(), ShouldEqual, internalSeverErrorMessage)
			})
		})
	})
}

func TestReviewsHandler(t *testing.T) {
	t.Parallel()

	Convey("Given an HTTP GET request to the /books/{id}/reviews endpoint", t, func() {

		Convey("When the {id} is empty", func() {
			api := &API{}
			request := httptest.NewRequest("GET", "/books/"+emptyID+"/reviews", nil)

			expectedUrlVars := map[string]string{
				"id":       emptyID,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewsHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, "empty book ID in request\n")
			})
		})
	})

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

		api := &API{
			dataStore: mockDataStore,
		}

		Convey("When a http get request is sent to /books/1/reviews", func() {

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{
				"id": bookID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewsHandler(response, request)
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

	Convey("Given an existing book with no reviews", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: bookID1}, nil
			},
			GetReviewsFunc: func(ctx context.Context, bookID string) (models.Reviews, error) {
				return emptyReviews, nil
			},
		}

		api := &API{
			dataStore: mockDataStore,
		}

		Convey("When a HTTP GET request is sent to /books/1/reviews", func() {
			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{
				"id": bookID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewsHandler(response, request)
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

	Convey("Given a GET request for a list of reviews of a book that does not exist", t, func() {

		Convey("When a http get request is sent to /books/1/reviews", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, apierrors.ErrBookNotFound
				},
			}

			api := &API{
				dataStore: mockDataStore,
			}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{
				"id": bookID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewsHandler(response, request)
			Convey("Then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
				So(response.Body.String(), ShouldContainSubstring, apierrors.ErrBookNotFound.Error())
			})
			Convey("And the GetReviews function is not called", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a GET request for a list of reviews of a book", t, func() {
		Convey("When GetReviews returns a database error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return &models.Book{ID: bookID1}, nil
				},
				GetReviewsFunc: func(ctx context.Context, bookID string) (models.Reviews, error) {
					return models.Reviews{}, errors.Wrap(errMongoDB, "unexpected error when getting a review")
				},
			}

			api := &API{
				dataStore: mockDataStore,
			}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{
				"id": bookID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewsHandler(response, request)
			Convey("Then the HTTP response code is 500", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
				So(response.Body.String(), ShouldEqual, internalSeverErrorMessage)
			})
			Convey("And the GetBook and GetReviews functions are called", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 1)
			})
		})

		Convey("When GetBook returns a database error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return &models.Book{}, errors.Wrap(errMongoDB, "unexpected error when getting a book")
				},
			}

			api := &API{
				dataStore: mockDataStore,
			}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{
				"id": bookID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewsHandler(response, request)
			Convey("Then the HTTP response code is 500", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
				So(response.Body.String(), ShouldEqual, internalSeverErrorMessage)
			})
			Convey("And the GetReviews function is not called", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 0)
			})
		})
	})
}
