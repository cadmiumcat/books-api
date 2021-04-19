package apierrors

import (
	"github.com/pkg/errors"
)

// Error messages for the books-api
var (
	ErrInvalidBook          = errors.New("invalid book. Missing required field")
	ErrEmptyRequestBody     = errors.New("empty request body")
	ErrEmptyBookID          = errors.New("empty book ID in request")
	ErrEmptyReviewID        = errors.New("empty review ID in request")
	ErrBookNotFound         = errors.New("book not found")
	ErrReviewNotFound       = errors.New("review not found")
	ErrReviewMissing        = errors.New("a review between 1 and 5 must be provided")
	ErrBookCheckedOut       = errors.New("this book is currently checked out")
	ErrNameMissing          = errors.New("a name must be provided for checkout")
	ErrBookNotCheckedOut    = errors.New("this book is not currently checked out")
	ErrUnableToReadMessage  = errors.New("failed to read request body")
	ErrUnableToParseJSON    = errors.New("failed to parse json body")
	ErrRequiredFieldMissing = errors.New("invalid book. Missing required field")
	ErrInternalServer       = errors.New("internal server error")
)
