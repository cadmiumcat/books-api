package main

import (
	"time"
)

type Book struct {
	Title   string
	Author  string
	Self    *Link
	History []Checkout
}

func (b Book) validate() error {
	if b.Title == "" || b.Author == "" {
		return ErrInvalidBook
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
	HRef string
	ID   string
}

var books []Book

func init() {
	books = append(books, Book{
		Title: "Book 1",
		Self: &Link{
			HRef: "amazon.com",
			ID:   "1",
		},
	})
}

func get(id string) (book *Book) {
	for i, l := range books {
		if l.Self.ID == id {
			book = &books[i]
			break
		}
	}
	return
}

func getAll() []Book {
	return books
}

func add(b Book) {
	books = append(books, b)
}

func checkout(b *Book, name string) error {
	h := len(b.History)
	if h != 0 {
		lastCheckout := b.History[h-1]
		if lastCheckout.In.IsZero() {
			return ErrBookCheckedOut
		}
	}

	if len(name) == 0 {
		return ErrNameMissing
	}

	b.History = append(b.History, Checkout{
		Who: name,
		Out: time.Now(),
	})

	return nil
}

func checkin(b *Book, review int) error {
	h := len(b.History)
	if h == 0 {
		return ErrBookNotCheckedOut
	}

	if review < 1 || review > 5 {
		return ErrReviewMissing
	}

	lastCheckout := b.History[h-1]
	if !lastCheckout.In.IsZero() {
		return ErrBookNotCheckedOut
	}

	b.History[h-1] = Checkout{
		Who:    lastCheckout.Who,
		Out:    lastCheckout.Out,
		In:     time.Now(),
		Review: review,
	}

	return nil
}
