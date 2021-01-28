package models

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBook_Validate(t *testing.T) {
	Convey("Given a book without any fields", t, func() {
		book := Book{}
		Convey("When I validate the book", func() {
			err := book.Validate()
			Convey("Then I get an error that tells me the book is invalid", func() {
				So(err, ShouldBeError, "invalid book. Missing required field")
			})
		})
	})

}
