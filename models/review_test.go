package models

import (
	"fmt"
	"github.com/cadmiumcat/books-api/apierrors"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

const bookID = "123"

func TestReview_Validate(t *testing.T) {

	tests := []struct {
		name     string
		input    Review
		expected error
	}{
		{
			name:     "Empty review",
			input:    Review{},
			expected: apierrors.ErrEmptyReviewMessage,
		},
		{
			name: "Long review",
			input: Review{
				Message: RandomString(t, 201),
				User:    User{Forenames: "Avid", Surname: "Reader"},
			},
			expected: apierrors.ErrLongReviewMessage,
		},
		{
			name:     "Review with no user",
			input:    Review{Message: "my review"},
			expected: apierrors.ErrEmptyReviewUser,
		},
	}

	Convey("Given a review", t, func() {
		for _, tt := range tests {
			Convey(fmt.Sprintf("When I validate the review: %v", tt.input), func() {
				err := tt.input.Validate()
				Convey(fmt.Sprintf("Then the error matches: %s", tt.expected), func() {
					So(err, ShouldBeError, tt.expected)
				})
			})
		}
	})

	Convey("Given a review with a valid message", t, func() {
		review := &Review{
			Message: "A perfect review. 10/10. Would read again",
			User: User{
				Forenames: "Reviewer", Surname: "OfBooks",
			},
		}
		Convey("When I validate it", func() {
			err := review.Validate()
			Convey("Then I get no errors", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNewReview(t *testing.T) {
	Convey("Given a bookID", t, func() {
		Convey("When a new review is created for that book", func() {
			review := NewReview(bookID)
			Convey("Then the review ID should not be empty", func() {
				So(review.ID, ShouldNotBeEmpty)
			})
			Convey("And the review's BookID should match the given bookID", func() {
				So(review.BookID, ShouldEqual, bookID)
			})
			Convey("And the review's book link is correct", func() {
				So(review.Links.Book, ShouldEqual, fmt.Sprintf("/books/%s", bookID))
			})
			Convey("And the review's self link should have the correct structure", func() {
				So(review.Links.Self, ShouldStartWith, fmt.Sprintf("/books/%s/reviews/", bookID))
			})
		})
	})

}

func RandomString(t *testing.T, n int) string {
	t.Helper()
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
