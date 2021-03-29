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
	"testing"
)

func TestReviews(t *testing.T) {
	hcMock := mock.HealthCheckerMock{}

	Convey("Given a existing book with at least one review (review_id=123)", t, func() {
		bookID := "1"
		reviewID := "123"
		mockDataStore := &mock.DataStoreMock{
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
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a existing book with no reviews", t, func() {
		bookID := "1"
		reviewID := "123"
		mockDataStore := &mock.DataStoreMock{
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
				So(mockDataStore.GetReviewCalls(), ShouldHaveLength, 1)
			})
		})
	})
}
