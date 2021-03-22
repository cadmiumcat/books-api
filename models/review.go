package models

type Review struct {
	ID string `json:"id" bson:"_id"`
}

// Reviews contains all the items (Review) in the library and a total count of those items
type Reviews struct {
	Count int    `json:"totalCount"`
	Items []Review `json:"items"`
}