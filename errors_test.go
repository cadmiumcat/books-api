package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	Convey("Given a checked out book", t, func() {
		book := Book{
			History: []Checkout{{Who: "user"}},
		}

		Convey("When a User tries checks out the book", func() {
			user := "otherUser"
			err := checkout(&book, user)

			Convey("Then an error message shows that the book is already checked out", func() {
				So(err, ShouldBeError, "this book is currently checked out")
			})
		})

	})

	Convey("Given a book in the book store", t, func() {
		book := Book{}
		Convey("When a user tries to check out a book without providing a name", func() {
			err := checkout(&book, "")
			Convey("Then an error message shows that a name must be provided", func() {
				So(err, ShouldBeError, "a name must be provided for checkout")
			})
		})
	})

}
