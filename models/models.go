package models

import (
	"errors"
	"time"
)

type Book struct {
	Id       string     `json:"id"`
	Title    string     `json:"title"`
	Author   string     `json:"author"`
	Synopsis string     `json:"synopsis"`
	Links    *Link      `json:"links"`
	History  []Checkout `json:"history"`
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
	Count int    `json:"total_count"`
	Items []Book `json:"items"`
}
