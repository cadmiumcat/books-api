package models

import (
	"fmt"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/pagination"
	uuid "github.com/satori/go.uuid"
	"time"
)

// A Book contains the fields that identify a book and its status.
type Book struct {
	ID       string     `json:"id" bson:"_id"`
	Title    string     `json:"title" bson:"title"`
	Author   string     `json:"author" bson:"author"`
	Synopsis string     `json:"synopsis,omitempty" bson:"synopsis,omitempty"`
	Links    *Link      `json:"links,omitempty" bson:"links,omitempty"`
	History  []Checkout `json:"history,omitempty" bson:"history,omitempty"`
}

// Validate checks a Book for missing required fields.
// It returns an error when required fields (e.g. author/title) are not provided.
func (b *Book) Validate() error {

	if b.Title == "" || b.Author == "" {
		return apierrors.ErrRequiredFieldMissing
	}

	return nil
}

// Checkout stores the details of when a someone has borrowed/returned a Book, as well as their review.
// To be deprecated
type Checkout struct {
	Who    string
	Out    time.Time
	In     time.Time
	Review int
}

// Link stores the details of when a someone has borrowed/returned a Book, as well user reviews.
type Link struct {
	Self         string `json:"self" bson:"self"`
	Reservations string `json:"reservations" bson:"reservations"`
	Reviews      string `json:"reviews" bson:"reviews"`
}

// BooksResponse represents a paginated list of Books
type BooksResponse struct {
	Items []Book `json:"items"`
	pagination.Page
}

// NewBook returns a Book structure
func NewBook() *Book {
	bookID := uuid.NewV4().String()
	return &Book{
		ID: bookID,
		Links: &Link{
			Self:    fmt.Sprintf("/books/%s", bookID),
			Reviews: fmt.Sprintf("/books/%s/reviews", bookID),
		},
	}
}
