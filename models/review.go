package models

import (
	"fmt"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/pagination"
	uuid "github.com/satori/go.uuid"
	"time"
)

// A Review contains the fields that identify a review
type Review struct {
	ID          string      `json:"id" bson:"_id"`
	User        User        `json:"user,omitempty" bson:"user,omitempty"`
	Message     string      `json:"message,omitempty" bson:"message,omitempty"`
	BookID      string      `json:"book_id" bson:"book_id"`
	Links       *ReviewLink `json:"links,omitempty" bson:"links,omitempty"`
	LastUpdated time.Time   `json:"last_updated" bson:"last_updated"`
}

// ReviewLink is the relationship between a Book and a Review
type ReviewLink struct {
	Self string `json:"self" bson:"self"`
	Book string `json:"book" bson:"book"`
}

func (r Review) Validate() error {
	if r.Message == "" {
		return apierrors.ErrEmptyReviewMessage
	}

	if r.User.Forenames == "" || r.User.Surname == "" {
		return apierrors.ErrEmptyReviewUser
	}

	if len(r.Message) > 200 {
		return apierrors.ErrLongReviewMessage
	}

	return nil
}

type User struct {
	Forenames string `json:"forenames,omitempty" bson:"forenames,omitempty"`
	Surname   string `json:"surname,omitempty" bson:"surname,omitempty"`
}

// ReviewsResponse represents a paginated list of Books
type ReviewsResponse struct {
	Items []Review `json:"items"`
	pagination.Page
}

// NewReview returns a Review structure based on a bookID
func NewReview(bookID string) *Review {
	reviewID := uuid.NewV4().String()

	return &Review{
		ID:     reviewID,
		BookID: bookID,
		Links: &ReviewLink{
			Self: fmt.Sprintf("/books/%s/reviews/%s", bookID, reviewID),
			Book: fmt.Sprintf("/books/%s", bookID),
		},
		LastUpdated: time.Now().UTC(),
	}
}
