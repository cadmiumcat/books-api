package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEndpoints(t *testing.T) {
	Convey("Given a POST request to add a book", t, func() {
		Convey("When the body does not contain a valid book", func() {
			body := strings.NewReader(`{}`)
			request, err := http.NewRequest(http.MethodPost, "/books", body)
			So(err, ShouldBeNil)

			response := httptest.NewRecorder()
			router := SetupRoutes()
			router.ServeHTTP(response, request)

			Convey("Then the HTTP response code is 400", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When the body contains a valid book", func() {
			body := strings.NewReader(`{ "title":"Girl, Woman, Other", "author":"Bernardine Evaristo" }`)
			request, err := http.NewRequest(http.MethodPost, "/books", body)
			So(err, ShouldBeNil)

			response := httptest.NewRecorder()
			router := SetupRoutes()
			router.ServeHTTP(response, request)

			Convey("Then the HTTP response code is 201", func() {
				So(response.Code, ShouldEqual, http.StatusCreated)
			})
		})
	})
}
