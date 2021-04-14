package api

import (
	"context"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/interfaces/mock"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const host = "localhost:8080"

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
			})
		})
	})

	Convey("Given an existing book with book id=1", t, func() {
		id := "1"
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return &models.Book{ID: "1"}, nil
			},
		}
		api := &API{dataStore: mockDataStore}
		Convey("When I send an HTTP GET request to /books/1", func() {
			request := httptest.NewRequest(http.MethodGet, "/books/"+id, nil)
			expectedUrlVars := map[string]string{
				"id": id,
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

	Convey("Given a book that does not exist with book id=3", t, func() {
		id := "3"
		mockDataStore := &mock.DataStoreMock{
			GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
				return nil, apierrors.ErrBookNotFound
			},
		}
		api := &API{dataStore: mockDataStore}

		Convey("When I send an HTTP GET request to /books/3", func() {

			request := httptest.NewRequest(http.MethodGet, "/books/"+id, nil)
			expectedUrlVars := map[string]string{
				"id": id,
			}
			request = mux.SetURLVars(request, expectedUrlVars)
			response := httptest.NewRecorder()

			api.getBookHandler(response, request)

			Convey("then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
			})
			Convey("And the GetBook function is called once", func() {
				So(mockDataStore.GetBookCalls(), ShouldHaveLength, 1)
			})
		})
	})

}

func TestGetBooksHandler(t *testing.T) {
	t.Parallel()
	Convey("Given ", t, func() {

		mockDataStore := &mock.DataStoreMock{
			GetBooksFunc: func(ctx context.Context) (models.Books, error) {
				return models.Books{}, nil
			},
		}

		api := &API{dataStore: mockDataStore}

		Convey("When I send an HTTP GET request to /books", func() {

			request := httptest.NewRequest(http.MethodGet, "/books", nil)
			response := httptest.NewRecorder()

			api.getBooksHandler(response, request)
			Convey("then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
			Convey("And the GetBooks function is called once", func() {
				So(mockDataStore.GetBooksCalls(), ShouldHaveLength, 1)
			})
		})
	})
}

func TestAddBookHandler(t *testing.T) {
	t.Parallel()

	Convey("Given a POST request to add a book", t, func() {
		mockDataStore := &mock.DataStoreMock{
			AddBookFunc: func(book *models.Book) {},
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
