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

	// ErrLimitOverMax represents an error case where the given limit value is larger than the allowed maximum
	ErrLimitOverMax = errors.New("limit query parameter is larger than the allowed maximum")
)

type Paginator struct {
	DefaultLimit        int
	DefaultOffset       int
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

// A Page is a section of paginated items, as well as the parameters used to determine the items that belong to the it
type Page struct {
	Count      int `json:"count"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
}

// GetPaginationValues returns pagination parameters based on a request, or the default values if the request does not specify them.
// It returns an error if the parameters are not valid
func (p *Paginator) GetPaginationValues(r *http.Request) (offset int, limit int, err error) {
	offsetParameter := r.URL.Query().Get("offset")
	limitParameter := r.URL.Query().Get("limit")

	offset = p.DefaultOffset
	limit = p.DefaultLimit

	if offsetParameter != "" {
		offset, err = strconv.Atoi(offsetParameter)
		if err != nil || offset < 0 {
			return 0, 0, ErrInvalidOffsetParameter
		}
	}

	if limitParameter != "" {
		limit, err = strconv.Atoi(limitParameter)
		if err != nil || limit < 0 {
			return 0, 0, ErrInvalidLimitParameter
		}
	}

	if limit > p.DefaultMaximumLimit {
		return 0, 0, ErrLimitOverMax
	}

	return
}
