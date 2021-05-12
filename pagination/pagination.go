package pagination

import (
	"encoding/json"
	"errors"
	"github.com/ONSdigital/log.go/log"
	"net/http"
	"reflect"
	"strconv"
)

var (
	errInvalidQueryOffset        = errors.New("invalid query parameter: offset")
	errInvalidQueryLimit         = errors.New("invalid query parameter: limit")
	errInvalidQueryLimitTooLarge = errors.New("invalid query parameter: limit exceeds maximum limit allowed")
)

type page struct {
	Items      interface{} `json:"items"`
	Count      int         `json:"count"`
	Offset     int         `json:"offset"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"total_count"`
}

// Handler is an interface for an endpoint that returns a list of values to be paginated
type Handler func(w http.ResponseWriter, r *http.Request, offset int, limit int) (list interface{}, totalCount int, err error)

type Paginator struct {
	DefaultLimit        int
	DefaultOffset       int
	DefaultMaximumLimit int
}

// NewPaginator creates a new instance of a Paginator with the specified default values
func NewPaginator(defaultLimit, defaultOffset, defaultMaximumLimit int) *Paginator {
	return &Paginator{
		DefaultLimit:        defaultLimit,
		DefaultOffset:       defaultOffset,
		DefaultMaximumLimit: defaultMaximumLimit,
	}
}

// Paginate wraps an HTTP endpoint to return a paginated list from the list returned by the provided Handler
func (p *Paginator) Paginate(handler Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		offset, limit, err := p.validateQueryParameters(r)
		logData := log.Data{"offset": offset, "limit": limit}
		if err != nil {
			log.Event(r.Context(), "api endpoint failed to paginate results", log.ERROR, log.Error(err), logData)
			http.Error(w, "invalid query parameters", http.StatusBadRequest)
		}

		list, totalCount, err := handler(w, r, offset, limit)

		page := &page{
			Items:      list,
			Count:      reflect.ValueOf(list).Len(),
			Offset:     offset,
			Limit:      limit,
			TotalCount: totalCount,
		}

		b, err := json.Marshal(page)

		if err != nil {
			log.Event(r.Context(), "api endpoint failed to marshal resource into bytes", log.ERROR, log.Error(err), logData)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if _, err = w.Write(b); err != nil {
			log.Event(r.Context(), "api endpoint error writing response body", log.ERROR, log.Error(err), logData)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Event(r.Context(), "api endpoint request successful", log.INFO, logData)

	}
}

// validateQueryParameters retrieves the offset and limit parameters in a query and returns them when they are valid
// If no parameters are provided, they are set to the default value in the Paginator
// An error is returned if the offset/limit is not a valid positive integer,
// or if the limit exceeds the DefaultMaximumLimit set by the Paginator
func (p *Paginator) validateQueryParameters(r *http.Request) (offset int, limit int, err error) {
	logData := log.Data{}

	offsetParameter := r.URL.Query().Get("offset")
	limitParameter := r.URL.Query().Get("limit")

	offset = p.DefaultOffset
	limit = p.DefaultLimit

	if offsetParameter != "" {
		logData["offset"] = offsetParameter
		offset, err = strconv.Atoi(offsetParameter)
		if err != nil || offset < 0 {
			log.Event(r.Context(), errInvalidQueryOffset.Error(), log.ERROR, logData)
			return 0, 0, errInvalidQueryOffset
		}
	}

	if limitParameter != "" {
		logData["limit"] = limitParameter
		limit, err = strconv.Atoi(limitParameter)
		if err != nil || limit < 0 {
			log.Event(r.Context(), errInvalidQueryLimit.Error(), log.ERROR, logData)
			return 0, 0, errInvalidQueryLimit
		}
	}

	if limit > p.DefaultMaximumLimit {
		log.Event(r.Context(), errInvalidQueryLimitTooLarge.Error(), log.ERROR, logData)
		return 0, 0, errInvalidQueryLimitTooLarge
	}

	return
}
