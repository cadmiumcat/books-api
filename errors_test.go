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

		Convey("When a User checks out the book", func() {
			user := "otherUser"
			err := checkout(&book, user)

			Convey("An error should say it's already checked out", func() {
				So(err, ShouldBeError, "this book is currently checked out to: user")
			})
		})

	})
}