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

			Convey("Then an error shows that the book is already checked out", func() {
				So(err, ShouldBeError, "this book is currently checked out")
			})
		})

	})
}