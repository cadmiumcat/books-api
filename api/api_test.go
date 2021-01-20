package api

import (
	"github.com/cadmiumcat/books-api/interfaces/datastoretest"
	"github.com/cadmiumcat/books-api/models"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	host = "http://localhost:80"
)

func TestEndpoints(t *testing.T) {
	response := httptest.NewRecorder()

	mockDataStore := &datastoretest.DataStoreMock{AddBookFunc: func(book *models.Book) {}}

	api := Setup(host, mux.NewRouter(), mockDataStore)

	Convey("Given a POST request to add a book", t, func() {
		Convey("When the body does not contain a valid book", func() {
			body := strings.NewReader(`{}`)
			request, err := http.NewRequest(http.MethodPost, "/books", body)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)

			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When the body contains a valid book", func() {
			body := strings.NewReader(`{"title":"Girl, Woman, Other", "author":"Bernardine Evaristo" }`)
			request, err := http.NewRequest(http.MethodPost, "/books", body)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)

			Convey("Then the HTTP response code is 201", func() {
				So(response.Code, ShouldEqual, http.StatusCreated)
			})
		})
	})

	Convey("Given an existing book with book id=1", t, func() {
		id := "1"

		Convey("When I send an HTTP GET request to /books/1", func() {
			request, err := http.NewRequest(http.MethodGet, "/books/"+id, nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)

			Convey("Then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
		})

	})

	Convey("Given a book that does not exist with book id=3", t, func() {
		id := "3"
		Convey("When I send an HTTP GET request to /books/3", func() {
			request, err := http.NewRequest(http.MethodGet, "/books/"+id, nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)
			Convey("then the HTTP response code is 404", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})

	Convey("Given ", t, func() {
		Convey("When I send an HTTP GET request to /books", func() {
			request, err := http.NewRequest(http.MethodGet, "/books", nil)
			So(err, ShouldBeNil)

			api.router.ServeHTTP(response, request)

			Convey("then the HTTP response code is 200", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}
