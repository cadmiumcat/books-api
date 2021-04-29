package models

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBook_Validate(t *testing.T) {
	Convey("Given a book with a title and an author", t, func() {
		book := Book{
			Title:  "Kindred",
			Author: "Octavia E. Butler",
		}
		Convey("When the book is validated", func() {
			err := book.Validate()
			Convey("Then no errors are returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a book with an author but no title", t, func() {
		book := Book{
			Author: "Octavia E. Butler",
		}
		Convey("When the book is validated", func() {
			err := book.Validate()
			Convey("Then an invalid book error is returned", func() {
				So(err, ShouldBeError, "invalid book. Missing required field")
			})
		})
	})

	Convey("Given a book with a title but no author", t, func() {
		book := Book{
			Title: "Kindred",
		}
		Convey("When the book is validated", func() {
			err := book.Validate()
			Convey("Then an invalid book error is returned", func() {
				So(err, ShouldBeError, "invalid book. Missing required field")
			})
		})
	})

	Convey("Given a book without any required fields", t, func() {
		book := Book{}
		Convey("When the book is validated", func() {
			err := book.Validate()
			Convey("Then an invalid book error is returned", func() {
				So(err, ShouldBeError, "invalid book. Missing required field")
			})
		})
	})
}

func TestNewBook(t *testing.T) {
	Convey("Given a new book is required", t, func() {
		Convey("When NewBook() is called", func() {
			book := NewBook()
			Convey("Then the book ID should not be empty", func() {
				So(book.ID, ShouldNotBeEmpty)
			})
			Convey("And the book's reviews link should have the correct structure", func() {
				So(book.Links.Reviews, ShouldStartWith, "/books/")
				So(book.Links.Reviews, ShouldEndWith, "/reviews")
			})
			Convey("And the book's self link should have the correct structure", func() {
				So(book.Links.Self, ShouldStartWith, "/books/")
			})
		})
	})
}
