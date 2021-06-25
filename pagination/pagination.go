package pagination

import (
	"errors"
	"net/http"
	"strconv"
)

var (
	// ErrInvalidOffsetParameter represents an error case where an invalid offset value is provided
	ErrInvalidOffsetParameter = errors.New("invalid offset query parameter")

	// ErrInvalidLimitParameter represents an error case where an invalid limit value is provided
	ErrInvalidLimitParameter = errors.New("invalid limit query parameter")

	// ErrLimitOverMax represents an error case where the given limit value is larger than the maximum allowed
	ErrLimitOverMax = errors.New("limit query parameter is larger than the maximum allowed")
)

type Paginator struct {
	DefaultLimit int
	DefaultOffset int
	DefaultMaximumLimit int
}

// NewPaginator creates a new instance of Paginator
func NewPaginator(limit, offset, maximumLimit int) *Paginator {
	return &Paginator{
		DefaultLimit:        limit,
		DefaultOffset:       offset,
		DefaultMaximumLimit: maximumLimit,
	}
}

// SetPaginationValues returns pagination parameters based on a request, or the default values if the request does not specify them.
// It returns an error if the parameters are not valid
func (p *Paginator) SetPaginationValues(r *http.Request) (offset int, limit int, err error) {
	offsetParameter := r.URL.Query().Get("offset")
	limitParameter := r.URL.Query().Get("limit")

	if offsetParameter != "" {
		p.DefaultOffset, err = strconv.Atoi(offsetParameter)
		if err != nil {
			return 0, 0, ErrInvalidOffsetParameter
		}
	}

	if limitParameter != "" {
		p.DefaultLimit, err = strconv.Atoi(limitParameter)
		if err != nil {
			return 0, 0, ErrInvalidLimitParameter
		}
	}

	return p.DefaultOffset, p.DefaultLimit, nil
}
