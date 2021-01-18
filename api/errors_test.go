package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	Convey("Given a checked out book", t, func() {
		book := Book{
			History: []Checkout{{Who: "user"}},
		}
		Convey("When a User tries to check out the book", func() {
			user := "otherUser"
			err := checkout(&book, user)

			Convey("Then an error message shows that the book is already checked out", func() {
				So(err, ShouldBeError, ErrBookCheckedOut)
			})
		})

		Convey("When a user tries to check in a book without providing a valid review", func() {
			err := checkin(&book, 100)
			Convey("Then an error message shows that a review must be provided", func() {
				So(err, ShouldBeError, ErrReviewMissing)
			})
		})

	})

	Convey("Given a book in the book store that has never been checked out", t, func() {
		book := Book{}
		Convey("When a user tries to check out a book without providing a name", func() {
			err := checkout(&book, "")
			Convey("Then an error message shows that a name must be provided", func() {
				So(err, ShouldBeError, ErrNameMissing)
			})
		})

		Convey("When a user tries to check in the book", func() {
			err := checkin(&book, 4)
			Convey("Then an error message shows that the book has not been checked out", func() {
				So(err, ShouldBeError, ErrBookNotCheckedOut)
			})
		})

	})

	Convey("Given an invalid book", t, func() {
		book := Book{}
		Convey("When the book is checked for validation", func() {
			err := book.validate()
			Convey("Then an error message shows that the book is invalid", func() {
				So(err, ShouldBeError, ErrInvalidBook)
			})
		})
	})

}
