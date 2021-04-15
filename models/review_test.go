package models

import (
	"fmt"
	"github.com/cadmiumcat/books-api/apierrors"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestReview_Validate(t *testing.T) {
	type review struct {
		Message string
	}
	tests := []struct {
		name     string
		input    review
		expected error
	}{
		{
			name:     "Empty review",
			input:    review{Message: ""},
			expected: apierrors.ErrEmptyReviewMessage,
		},
		{
			name:     "Long review",
			input:    review{Message: RandomString(201)},
			expected: apierrors.ErrLongReviewMessage,
		},
	}

	Convey("Given a review", t, func() {
		for _, tt := range tests {
			r := Review{
				Message: tt.input.Message,
			}
			Convey(fmt.Sprintf("When I validate the review with the message %q", tt.input.Message), func() {
				err := r.Validate()
				Convey(fmt.Sprintf("Then the error matches expected %s", tt.expected), func() {
					So(err, ShouldBeError, tt.expected)
				})
			})
		}
	})

	Convey("Given a review with a valid message", t, func() {
		review := &Review{Message: "A perfect review. 10/10. Would read again"}
		Convey("When I validate it", func() {
			err := review.Validate()
			Convey("Then I get no errors", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
