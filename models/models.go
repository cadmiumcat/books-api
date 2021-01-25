package models

import (
	"errors"
	"time"
)

type Book struct {
	ID       string     `json:"id" bson:"_id"`
	Title    string     `json:"title" bson:"title"`
	Author   string     `json:"author" bson:"author"`
	Synopsis string     `json:"synopsis,omitempty" bson:"synopsis"`
	Links    *Link      `json:"links,omitempty" bson:"links"`
	History  []Checkout `json:"history,omitempty" bson:"history"`
}

func (b *Book) Validate() error {
	if b.Title == "" || b.Author == "" {
		return errors.New("invalid book. Missing required field")
	}

	return nil
}

type Checkout struct {
	Who    string
	Out    time.Time
	In     time.Time
	Review int
}

type Link struct {
	Self         string
	Reservations string
	Reviews      string
}

type Books struct {
	Count int    `json:"totalCount" bson:"totalCount"`
	Items []Book `json:"items" bson:"items"`
}
