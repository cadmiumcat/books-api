package models

// A Review contains the fields that identify a review
type Review struct {
	ID    string `json:"id" bson:"_id"`
	Links *link  `json:"links,omitempty" bson:"links,omitempty"`
}

// link is the relationship between a Book and a Review
type link struct {
	Self string
	Book string
}

// Reviews contains all the items (Review) in the library and a total count of those items
type Reviews struct {
	Count int      `json:"totalCount"`
	Items []Review `json:"items"`
}
