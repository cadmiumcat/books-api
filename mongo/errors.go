package mongo

import "github.com/pkg/errors"

var (
	ErrBookNotFound   = errors.New("book not found")
	ErrReviewNotFound = errors.New("review not found")
)
