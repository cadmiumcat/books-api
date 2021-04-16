package models

import "github.com/cadmiumcat/books-api/apierrors"

// A Review contains the fields that identify a review
type Review struct {
	ID      string      `json:"id" bson:"_id"`
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

	if len(r.Message) > 200 {
		return apierrors.ErrLongReviewMessage
	}

	return nil
}
