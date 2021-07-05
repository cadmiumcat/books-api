package api

import (
	"context"
	"encoding/json"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/interfaces/mock"
	"github.com/cadmiumcat/books-api/models"
	"github.com/cadmiumcat/books-api/mongo"
	"github.com/cadmiumcat/books-api/pagination"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	bookID1                   = "1"
	bookID2                   = "2"
	reviewID1                 = "123"
	reviewID2                 = "567"
	emptyID                   = ""
	bookIDNotInStore          = "notInStore"
	reviewInvalidMessage      = `{"message": ""}`
	reviewInvalidUpdate       = `{""`
	reviewValid               = `{"message": "my review", "user": {"forenames": "name", "surname": "surname"}}`
	internalSeverErrorMessage = "internal server error\n"
)

var bookReview1 = models.Review{
	ID: reviewID1,
	Links: &models.ReviewLink{
		Book: bookID1,
	},
	User: models.User{
		Forenames: "name",
		Surname:   "surname",
	},
}

var bookReview2 = models.Review{
	ID: reviewID2,
	Links: &models.ReviewLink{
		Book: bookID1,
	},
}

var reviewUpdated = models.Review{
	ID: reviewID1,
	User: models.User{
		Forenames: "new name",
		Surname:   "old surname",
	},
	Message:     "old review",
	LastUpdated: time.Now(),
}

var reviewUpdate = models.Review{
	User: models.User{
		Forenames: "new Name",
	},
	Message: "",
}

var errMongoDB = errors.New("unexpected error in MongoDB")

func marshalJSON(t *testing.T, data interface{}) string {
	t.Helper()
	out, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return string(out)
}

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
				return &bookReview1, nil
			},
		}

		api := &API{dataStore: mockDataStore}

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
				So(response.Body.String(), ShouldEqual, marshalJSON(t, bookReview1))
			})
		})
	})

	Convey("Given an existing book with no reviews", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: bookID1}, nil
			},
			GetReviewFunc: func(ctx context.Context, reviewID string) (*models.Review, error) {
				return nil, mongo.ErrReviewNotFound
			},
		}

		api := &API{dataStore: mockDataStore}

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
					return nil, mongo.ErrBookNotFound
				},
			}

			api := &API{dataStore: mockDataStore}

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

			api := &API{dataStore: mockDataStore}

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

			api := &API{dataStore: mockDataStore}

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

func TestGetReviewsHandler(t *testing.T) {
	t.Parallel()

	Convey("Given an HTTP GET request to the /books/{id}/reviews endpoint", t, func() {

		Convey("When the {id} is empty", func() {
			paginator := mockPaginator()

			api := &API{paginator: paginator}
			request := httptest.NewRequest("GET", "/books/"+emptyID+"/reviews", nil)

			expectedUrlVars := map[string]string{"id": emptyID}
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
			GetReviewsFunc: func(ctx context.Context, bookID string, offset int, limit int) ([]models.Review, int, error) {
				return []models.Review{
					bookReview1,
					bookReview2,
				}, 2, nil
			},
		}

		paginator := mockPaginator()

		api := &API{dataStore: mockDataStore, paginator: paginator}

		Convey("When a http get request is sent to /books/1/reviews", func() {

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{"id": bookID1}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewsHandler(response, request)
			Convey("Then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
			Convey("And the GetReviews function is called once, and the pagination parameters passed to it", func() {
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls()[0].Limit, ShouldEqual, limit)
				So(mockDataStore.GetReviewsCalls()[0].Offset, ShouldEqual, offset)
			})
			Convey("And the paginator is called to extract the pagination parameters ", func() {
				So(paginator.GetPaginationValuesCalls(), ShouldHaveLength, 1)
				So(paginator.GetPaginationValuesCalls()[0].R, ShouldEqual, request)
			})
			Convey("And the response body contains the right number of reviews", func() {
				payload, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)
				page := models.ReviewsResponse{}
				err = json.Unmarshal(payload, &page)
				So(err, ShouldBeNil)
				expectedPage := pagination.Page{Count: 2, Offset: offset, Limit: limit, TotalCount: 2}

				So(page.TotalCount, ShouldEqual, expectedPage.TotalCount)
				So(page.Count, ShouldEqual, len(page.Items))
				So(page.Offset, ShouldEqual, offset)
				So(page.Limit, ShouldEqual, limit)
				So(page.Page, ShouldResemble, expectedPage)
			})
		})
	})

	Convey("Given an existing book with no reviews", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: bookID1}, nil
			},
			GetReviewsFunc: func(ctx context.Context, bookID string, offset int, limit int) ([]models.Review, int, error) {
				return []models.Review{}, 0, nil
			},
		}

		paginator := mockPaginator()
		api := &API{dataStore: mockDataStore, paginator: paginator}

		Convey("When a HTTP GET request is sent to /books/1/reviews", func() {
			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{"id": bookID1}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewsHandler(response, request)
			Convey("Then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
			Convey("And the GetReviews function is called once, and the pagination parameters passed to it", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewsCalls()[0].Limit, ShouldEqual, limit)
				So(mockDataStore.GetReviewsCalls()[0].Offset, ShouldEqual, offset)
			})
			Convey("And the paginator is called to extract the pagination parameters ", func() {
				So(paginator.GetPaginationValuesCalls(), ShouldHaveLength, 1)
				So(paginator.GetPaginationValuesCalls()[0].R, ShouldEqual, request)
			})
			Convey("And the response contains a count of zero and no review items", func() {
				payload, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)
				page := models.ReviewsResponse{}
				err = json.Unmarshal(payload, &page)
				So(err, ShouldBeNil)
				expectedPage := pagination.Page{Count: 0, Offset: offset, Limit: limit, TotalCount: 0}

				So(page.TotalCount, ShouldEqual, expectedPage.TotalCount)
				So(page.Count, ShouldEqual, len(page.Items))
				So(page.Offset, ShouldEqual, offset)
				So(page.Limit, ShouldEqual, limit)
				So(page.Page, ShouldResemble, expectedPage)
			})
		})
	})

	Convey("Given a GET request for a list of reviews of a book that does not exist", t, func() {

		Convey("When a http get request is sent to /books/1/reviews", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, mongo.ErrBookNotFound
				},
			}

			paginator := mockPaginator()
			api := &API{dataStore: mockDataStore, paginator: paginator}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{"id": bookID1}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getReviewsHandler(response, request)
			Convey("Then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
				So(response.Body.String(), ShouldContainSubstring, mongo.ErrBookNotFound.Error())
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
				GetReviewsFunc: func(ctx context.Context, bookID string, offset int, limit int) ([]models.Review, int, error) {
					return []models.Review{}, 0, errors.Wrap(errMongoDB, "unexpected error when getting a review")
				},
			}

			paginator := mockPaginator()
			api := &API{dataStore: mockDataStore, paginator: paginator}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{"id": bookID1}
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
				So(mockDataStore.GetReviewsCalls()[0].Limit, ShouldEqual, limit)
				So(mockDataStore.GetReviewsCalls()[0].Offset, ShouldEqual, offset)
			})
			Convey("And the paginator is called to extract the pagination parameters ", func() {
				So(paginator.GetPaginationValuesCalls(), ShouldHaveLength, 1)
				So(paginator.GetPaginationValuesCalls()[0].R, ShouldEqual, request)
			})
		})

		Convey("When GetBook returns a database error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return &models.Book{}, errors.Wrap(errMongoDB, "unexpected error when getting a book")
				},
			}

			paginator := mockPaginator()
			api := &API{dataStore: mockDataStore, paginator: paginator}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{"id": bookID1}
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

func TestAddReviewHandler(t *testing.T) {
	Convey("Given an HTTP POST request to the /books/{id}/reviews endpoint", t, func() {

		Convey("When the book exists and the review is valid", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return &book1, nil
				},
				AddReviewFunc: func(ctx context.Context, review *models.Review) error {
					return nil
				},
			}

			api := &API{dataStore: mockDataStore}
			body := strings.NewReader(reviewValid)
			request := httptest.NewRequest("POST", "/books/"+bookID1+"/reviews", body)

			expectedUrlVars := map[string]string{"id": bookID1}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.addReviewHandler(response, request)
			Convey("Then the HTTP response code is 201", func() {
				So(response.Code, ShouldEqual, http.StatusCreated)
			})
			Convey("And the AddReview function is called", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.AddReviewCalls(), ShouldHaveLength, 1)
			})
		})

		Convey("When the book exist, but the review is not valid (empty message)", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return &book1, nil
				},
				AddReviewFunc: func(ctx context.Context, review *models.Review) error {
					return nil
				},
			}

			api := &API{dataStore: mockDataStore}
			body := strings.NewReader(reviewInvalidMessage)
			request := httptest.NewRequest("POST", "/books/"+bookID1+"/reviews", body)

			expectedUrlVars := map[string]string{"id": bookID1}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.addReviewHandler(response, request)
			Convey("Then the HTTP response is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, "empty review provided. Please enter a message\n")
			})
		})

		Convey("When the {id} is empty", func() {
			api := &API{}
			request := httptest.NewRequest("POST", "/books/"+emptyID+"/reviews", nil)

			expectedUrlVars := map[string]string{"id": emptyID}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.addReviewHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, "empty book ID in request\n")
			})
		})

		Convey("When there is no request body", func() {
			api := &API{}
			request := httptest.NewRequest("POST", "/books/"+bookID1+"/reviews", nil)

			expectedUrlVars := map[string]string{"id": bookID1}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.addReviewHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, "empty request body\n")
			})
		})

		Convey("When the book does not exits in the datastore", func() {
			mockDataStore := mock.DataStoreMock{GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{}, mongo.ErrBookNotFound
			}}

			api := &API{dataStore: &mockDataStore}
			body := strings.NewReader(reviewValid)
			request := httptest.NewRequest("POST", "/books/"+bookID1+"/reviews", body)

			expectedUrlVars := map[string]string{"id": bookID1}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.addReviewHandler(response, request)
			Convey("Then the HTTP response is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
				So(response.Body.String(), ShouldEqual, "book not found\n")
			})
		})

		Convey("When the request body is invalid", func() {
			mockDataStore := mock.DataStoreMock{GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{}, nil
			}}

			api := &API{dataStore: &mockDataStore}
			body := strings.NewReader("invalidReviewText")
			request := httptest.NewRequest("POST", "/books/"+bookID1+"/reviews", body)

			expectedUrlVars := map[string]string{"id": bookID1}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.addReviewHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, "invalid review\n")
			})
		})
	})
}

func TestUpdateReviewHandler(t *testing.T) {
	Convey("Given an HTTP PUT request to the /books/{id}/reviews/{review_id} endpoint", t, func() {

		Convey("When the book and review exist, and the review update is valid", func() {
			mockDataStore := mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, nil
				},
				GetReviewFunc: func(ctx context.Context, reviewID string) (*models.Review, error) {
					return &reviewUpdated, nil
				},
				UpdateReviewFunc: func(ctx context.Context, reviewID string, review *models.Review) error {
					return nil
				},
			}
			api := API{dataStore: &mockDataStore}

			body := strings.NewReader(marshalJSON(t, reviewUpdate))
			request := httptest.NewRequest("PUT", "/books/"+bookID1+"/reviews"+reviewID1, body)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()
			api.updateReviewHandler(response, request)
			Convey("Then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 1)
				So(mockDataStore.UpdateReviewCalls(), ShouldHaveLength, 1)
			})
		})

		Convey("When the book and review exist, but the review update is not valid", func() {
			mockDataStore := mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, nil
				},
				GetReviewFunc: func(ctx context.Context, reviewID string) (*models.Review, error) {
					return &models.Review{}, nil
				},
			}
			api := API{dataStore: &mockDataStore}

			body := strings.NewReader(reviewInvalidUpdate)
			request := httptest.NewRequest("PUT", "/books/"+bookID1+"/reviews"+reviewID1, body)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.updateReviewHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("And it returns an error saying the review is invalid", func() {
				So(response.Body.String(), ShouldEqual, "invalid review\n")
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 1)
				So(mockDataStore.UpdateReviewCalls(), ShouldHaveLength, 0)
			})
		})

		Convey("When the book does not exist", func() {
			mockDataStore := mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, mongo.ErrBookNotFound
				},
			}
			api := API{dataStore: &mockDataStore}

			body := strings.NewReader(reviewInvalidUpdate)
			request := httptest.NewRequest("PUT", "/books/"+bookID1+"/reviews"+reviewID1, body)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.updateReviewHandler(response, request)
			Convey("Then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
			})
			Convey("And it returns an error saying the book was not found", func() {
				So(response.Body.String(), ShouldEqual, "book not found\n")
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 0)
				So(mockDataStore.UpdateReviewCalls(), ShouldHaveLength, 0)
			})
		})

		Convey("When the review does not exist", func() {
			mockDataStore := mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, nil
				},
				GetReviewFunc: func(ctx context.Context, reviewID string) (*models.Review, error) {
					return nil, mongo.ErrReviewNotFound
				},
			}
			api := API{dataStore: &mockDataStore}

			body := strings.NewReader(reviewInvalidUpdate)
			request := httptest.NewRequest("PUT", "/books/"+bookID1+"/reviews"+reviewID1, body)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.updateReviewHandler(response, request)
			Convey("Then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
			})
			Convey("And it returns an error saying the review was not found", func() {
				So(response.Body.String(), ShouldEqual, "review not found\n")
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 1)
				So(mockDataStore.UpdateReviewCalls(), ShouldHaveLength, 0)
			})
		})

		Convey("When the {id} is empty", func() {
			mockDataStore := &mock.DataStoreMock{}
			api := API{dataStore: mockDataStore}

			body := strings.NewReader(reviewInvalidUpdate)
			request := httptest.NewRequest("PUT", "/books/"+emptyID+"/reviews"+reviewID1, body)

			expectedUrlVars := map[string]string{
				"id":       emptyID,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.updateReviewHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("And it returns an error saying the book ID is empty", func() {
				So(response.Body.String(), ShouldEqual, "empty book ID in request\n")
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 0)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 0)
				So(mockDataStore.UpdateReviewCalls(), ShouldHaveLength, 0)
			})
		})

		Convey("When the {review_id} is empty", func() {
			mockDataStore := &mock.DataStoreMock{}
			api := API{dataStore: mockDataStore}

			body := strings.NewReader(reviewInvalidUpdate)
			request := httptest.NewRequest("PUT", "/books/"+bookID1+"/reviews"+emptyID, body)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": emptyID,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.updateReviewHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("And it returns an error saying the review ID is empty", func() {
				So(response.Body.String(), ShouldEqual, "empty review ID in request\n")
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 0)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 0)
				So(mockDataStore.UpdateReviewCalls(), ShouldHaveLength, 0)
			})
		})

		Convey("When the GetBook returns an error", func() {
			mockDataStore := mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, errMongoDB
				},
			}
			api := API{dataStore: &mockDataStore}

			body := strings.NewReader(marshalJSON(t, reviewUpdate))
			request := httptest.NewRequest("PUT", "/books/"+bookID1+"/reviews"+reviewID1, body)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.updateReviewHandler(response, request)
			Convey("Then the HTTP response code is 500", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
			})
			Convey("And it returns an internal server error", func() {
				So(response.Body.String(), ShouldEqual, "internal server error\n")
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 0)
				So(mockDataStore.UpdateReviewCalls(), ShouldHaveLength, 0)
			})
		})

		Convey("When the GetReview returns an error", func() {
			mockDataStore := mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, nil
				},
				GetReviewFunc: func(ctx context.Context, reviewID string) (*models.Review, error) {
					return nil, apierrors.ErrInternalServer
				},
			}
			api := API{dataStore: &mockDataStore}

			body := strings.NewReader(marshalJSON(t, reviewUpdate))
			request := httptest.NewRequest("PUT", "/books/"+bookID1+"/reviews"+reviewID1, body)

			expectedUrlVars := map[string]string{
				"id":       bookID1,
				"reviewID": reviewID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.updateReviewHandler(response, request)
			Convey("Then the HTTP response code is 500", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
			})
			Convey("And it returns an internal server error", func() {
				So(response.Body.String(), ShouldEqual, "internal server error\n")
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 1)
				So(mockDataStore.UpdateReviewCalls(), ShouldHaveLength, 0)
			})
		})
	})
}
