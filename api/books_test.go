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
)

var (
	offset = 0
	limit  = 3
)

var book1 = models.Book{
	ID:     bookID1,
	Title:  "Girl, Woman, Other",
	Author: "Bernardine Evaristo",
}

var book2 = models.Book{
	ID:     bookID2,
	Title:  "Girl, Woman, Other",
	Author: "Bernardine Evaristo",
}

func TestGetBookHandler(t *testing.T) {
	t.Parallel()

	Convey("Given an HTTP GET request to the /books/{id} endpoint", t, func() {

		api := &API{}

		Convey("When the {id} is empty", func() {
			request := httptest.NewRequest("GET", "/books/"+emptyID, nil)

			response := httptest.NewRecorder()

			api.getBookHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldContainSubstring, apierrors.ErrEmptyBookID.Error())
			})
		})
	})

	Convey("Given an existing book with book id=1", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: bookID1}, nil
			},
		}
		api := &API{dataStore: mockDataStore}
		Convey("When a http get request is sent to /books/1", func() {
			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1, nil)
			expectedUrlVars := map[string]string{
				"id": bookID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getBookHandler(response, request)
			Convey("Then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
			Convey("And the GetBook function is called once", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
			})
		})

	})

	Convey("Given a book that does not exist", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return nil, mongo.ErrBookNotFound
			},
		}
		api := &API{dataStore: mockDataStore}

		Convey("When a http get request is sent to /books/3", func() {

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookIDNotInStore, nil)
			expectedUrlVars := map[string]string{
				"id": bookIDNotInStore,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getBookHandler(response, request)

			Convey("then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
				So(response.Body.String(), ShouldContainSubstring, mongo.ErrBookNotFound.Error())
			})
			Convey("And the GetBook function is called once", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a GET request for a book", t, func() {
		Convey("When GetBook returns an unexpected database error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
					return nil, errors.Wrap(errMongoDB, "unexpected error when getting a book")
				},
			}
			api := &API{dataStore: mockDataStore}

			request := httptest.NewRequest(http.MethodGet, "/books/"+bookID1, nil)
			expectedUrlVars := map[string]string{
				"id": bookID1,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getBookHandler(response, request)

			Convey("Then 500 InternalServerError status code is returned", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
				So(response.Body.String(), ShouldEqual, internalSeverErrorMessage)
			})
		})
	})

}

func TestGetBooksHandler(t *testing.T) {
	t.Parallel()
	Convey("Given a datastore with no books", t, func() {

		mockDataStore := &mock.DataStoreMock{
			GetBooksFunc: func(ctx context.Context, offset int, limit int) ([]models.Book, int, error) {
				return []models.Book{}, 0, nil
			},
		}
		paginator := mockPaginator()

		api := &API{dataStore: mockDataStore, paginator: paginator}

		Convey("When a http get request is sent to /books", func() {

			request := httptest.NewRequest(http.MethodGet, "/books", nil)
			response := httptest.NewRecorder()

			api.getBooksHandler(response, request)
			Convey("then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
			Convey("And the GetBooks function is called once, and the pagination parameters passed to it", func() {
				So(mockDataStore.GetBooksCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetBooksCalls()[0].Limit, ShouldEqual, limit)
				So(mockDataStore.GetBooksCalls()[0].Offset, ShouldEqual, offset)
			})

			Convey("And the paginator is called to extract the pagination parameters ", func() {
				So(paginator.GetPaginationValuesCalls(), ShouldHaveLength, 1)
				So(paginator.GetPaginationValuesCalls()[0].R, ShouldEqual, request)
			})

			Convey("And the response contains a count of zero and no book items", func() {
				payload, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)
				page := models.BooksResponse{}
				err = json.Unmarshal(payload, &page)
				So(err, ShouldBeNil)
				expectedPage := pagination.Page{Count: 0, Offset: offset, Limit: limit, TotalCount: 0}

				So(page.TotalCount, ShouldBeZeroValue)
				So(page.Count, ShouldEqual, len(page.Items))
				So(page.Offset, ShouldEqual, offset)
				So(page.Limit, ShouldEqual, limit)
				So(page.Page, ShouldResemble, expectedPage)
			})
		})
	})

	Convey("Given a datastore with 2 books", t, func() {
		mockDataStore := &mock.DataStoreMock{
			GetBooksFunc: func(ctx context.Context, offset int, limit int) ([]models.Book, int, error) {
				return []models.Book{book1, book2}, 2, nil
			},
		}

		paginator := mockPaginator()
		api := &API{dataStore: mockDataStore, paginator: paginator}

		Convey("When a http get request is sent to /books", func() {
			request := httptest.NewRequest(http.MethodGet, "/books", nil)
			response := httptest.NewRecorder()

			api.getBooksHandler(response, request)
			Convey("then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
			Convey("And the GetBooks function is called once", func() {
				So(mockDataStore.GetBooksCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetBooksCalls()[0].Limit, ShouldEqual, limit)
				So(mockDataStore.GetBooksCalls()[0].Offset, ShouldEqual, offset)
			})

			Convey("And the paginator is called to extract the pagination parameters ", func() {
				So(paginator.GetPaginationValuesCalls(), ShouldHaveLength, 1)
				So(paginator.GetPaginationValuesCalls()[0].R, ShouldEqual, request)
			})

			Convey("And the response contains the paginated response", func() {
				payload, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)
				page := models.BooksResponse{}
				err = json.Unmarshal(payload, &page)
				So(err, ShouldBeNil)
				expectedPage := pagination.Page{Count: 2, Offset: offset, Limit: limit, TotalCount: 2}

				So(page.TotalCount, ShouldEqual, 2)
				So(page.Count, ShouldEqual, len(page.Items))
				So(page.Offset, ShouldEqual, offset)
				So(page.Limit, ShouldEqual, limit)
				So(page.Page, ShouldResemble, expectedPage)
			})
		})
	})

	Convey("Given a GET request for a list of books", t, func() {
		Convey("When GetBooks returns an unexpected database error", func() {
			mockDataStore := &mock.DataStoreMock{
				GetBooksFunc: func(ctx context.Context, offset int, limit int) ([]models.Book, int, error) {
					return []models.Book{}, 0, errors.Wrap(errMongoDB, "unexpected error when getting books")
				},
			}

			paginator := mockPaginator()
			api := &API{dataStore: mockDataStore, paginator: paginator}

			request := httptest.NewRequest(http.MethodGet, "/books", nil)
			response := httptest.NewRecorder()

			api.getBooksHandler(response, request)

			Convey("Then the GetBooks function is called once", func() {
				So(mockDataStore.GetBooksCalls(), ShouldHaveLength, 1)
				So(mockDataStore.GetBooksCalls()[0].Limit, ShouldEqual, limit)
				So(mockDataStore.GetBooksCalls()[0].Offset, ShouldEqual, offset)
			})

			Convey("And the paginator is called to extract the pagination parameters ", func() {
				So(paginator.GetPaginationValuesCalls(), ShouldHaveLength, 1)
				So(paginator.GetPaginationValuesCalls()[0].R, ShouldEqual, request)
			})

			Convey("And a 500 InternalServerError status code is returned", func() {
				So(response.Code, ShouldEqual, http.StatusInternalServerError)
				So(response.Body.String(), ShouldEqual, internalSeverErrorMessage)
			})
		})
	})

}

func TestAddBookHandler(t *testing.T) {
	t.Parallel()

	Convey("Given a POST request to add a book", t, func() {
		mockDataStore := &mock.DataStoreMock{
			AddBookFunc: func(ctx context.Context, book *models.Book) error {
				return nil
			},
		}

		Convey("When there is no request body", func() {
			api := &API{}

			body := strings.NewReader(``)
			request := httptest.NewRequest(http.MethodPost, "/books", body)

			response := httptest.NewRecorder()

			api.addBookHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("And there AddBook function is not called", func() {
				So(mockDataStore.AddBookCalls(), ShouldHaveLength, 0)
			})
			Convey("And the response says the request body is missing", func() {
				So(response.Body.String(), ShouldContainSubstring, apierrors.ErrEmptyRequestBody.Error())
			})
		})

		Convey("When the body is empty", func() {
			api := &API{}

			body := strings.NewReader(`{}`)
			request := httptest.NewRequest(http.MethodPost, "/books", body)

			response := httptest.NewRecorder()

			api.addBookHandler(response, request)
			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("And there AddBook function is not called", func() {
				So(mockDataStore.AddBookCalls(), ShouldHaveLength, 0)
			})
			Convey("And the response says the request is empty", func() {
				So(response.Body.String(), ShouldContainSubstring, apierrors.ErrRequiredFieldMissing.Error())
			})
		})

		Convey("When the body contains a valid book", func() {
			api := &API{dataStore: mockDataStore}

			body := strings.NewReader(`{"title":"Girl, Woman, Other", "author":"Bernardine Evaristo" }`)
			request := httptest.NewRequest(http.MethodPost, "/books", body)

			response := httptest.NewRecorder()

			api.addBookHandler(response, request)
			Convey("Then the HTTP response code is 201", func() {
				So(response.Code, ShouldEqual, http.StatusCreated)
			})
			Convey("And the AddBook function is called once", func() {
				So(mockDataStore.AddBookCalls(), ShouldHaveLength, 1)
			})
		})
	})
}

func mockPaginator() *mock.PaginatorMock {
	paginator := &mock.PaginatorMock{
		GetPaginationValuesFunc: func(r *http.Request) (int, int, error) {
			return offset, limit, nil
		},
	}
	return paginator
}
