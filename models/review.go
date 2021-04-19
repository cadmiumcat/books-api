package models

import "github.com/cadmiumcat/books-api/apierrors"

// A Review contains the fields that identify a review
type Review struct {
	ID      string      `json:"id" bson:"_id"`
	User    User        `json:"user,omitempty" bson:"user,omitempty"`
	Message string      `json:"message,omitempty" bson:"message,omitempty"`
	Links   *ReviewLink `json:"links,omitempty" bson:"links,omitempty"`
}

// ReviewLink is the relationship between a Book and a Review
type ReviewLink struct {
	Self string `json:"self" bson:"self"`
	Book string `json:"book" bson:"book"`
}

// Reviews contains all the items (Review) in the library and a total count of those items
type Reviews struct {
	Count int      `json:"totalCount"`
	Items []Review `json:"items"`
}

func (r Review) Validate() error {

	if r.Message == "" {
		return apierrors.ErrEmptyReviewMessage
	}

	if r.User.Forename == "" || r.User.Surname == "" {
		return apierrors.ErrEmptyReviewUser
	}

	if len(r.Message) > 200 {
		return apierrors.ErrLongReviewMessage
	}

	return nil
}

type User struct {
	Forename string `json:"forename,omitempty" bson:"forename,omitempty"`
	Surname  string `json:"surname,omitempty" bson:"surname,omitempty"`
}
