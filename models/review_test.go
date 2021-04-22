package models

import (
	"fmt"
	"github.com/cadmiumcat/books-api/apierrors"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestReview_Validate(t *testing.T) {

	tests := []struct {
		name     string
		input    Review
		expected error
	}{
		{
			name:     "Empty review",
			input:    Review{Message: ""},
			expected: apierrors.ErrEmptyReviewMessage,
		},
		{
			name: "Long review",
			input: Review{
				Message: RandomString(t,201),
				User:    User{Forename: "Avid", Surname: "Reader"},
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
				Forename: "Reviewer", Surname: "OfBooks",
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

func RandomString(t *testing.T, n int) string {
	t.Helper()
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
