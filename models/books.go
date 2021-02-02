package models

import (
	"errors"
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
// It returns an error when
func (b *Book) Validate() error {
	if b.Title == "" || b.Author == "" {
		return errors.New("invalid book. Missing required field")
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
	Self         string
	Reservations string
	Reviews      string
}

// Books contains all the items (Book) in the library and a total count of those items
type Books struct {
	Count int    `json:"totalCount"`
	Items []Book `json:"items"`
}
