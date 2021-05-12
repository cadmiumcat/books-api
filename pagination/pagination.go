package pagination

import (
	"encoding/json"
	"errors"
	"github.com/ONSdigital/log.go/log"
	"net/http"
	"reflect"
	"strconv"
)

type page struct {
	Items      interface{} `json:"items"`
	Count      int         `json:"count"`
	Offset     int         `json:"offset"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"total_count"`
}

type Handler func(w http.ResponseWriter, r *http.Request, offset int, limit int) (list interface{}, totalCount int, err error)

type Paginator struct {
	DefaultLimit        int
	DefaultOffset       int
	DefaultMaximumLimit int
}

func NewPaginator(defaultLimit, defaultOffset, defaultMaximumLimit int) *Paginator {
	return &Paginator{
		DefaultLimit:        defaultLimit,
		DefaultOffset:       defaultOffset,
		DefaultMaximumLimit: defaultMaximumLimit,
	}
}

func (p *Paginator) Paginate(handler Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		offset, limit, err := p.validateQueryParameters(r)
		logData := log.Data{"offset": offset, "limit": limit}
		if err != nil {
			log.Event(r.Context(), "api endpoint failed to paginate results", log.ERROR, log.Error(err), logData)
			http.Error(w, "invalid query parameters", http.StatusBadRequest)
		}

		//offset := p.DefaultOffset
		//limit := p.DefaultLimit

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

func (p *Paginator) validateQueryParameters(r *http.Request) (offset int, limit int, err error) {
	logData := log.Data{}

	offsetParameter := r.URL.Query().Get("offset")
	limitParameter := r.URL.Query().Get("limit")

	offset = p.DefaultOffset
	limit = p.DefaultLimit

	if offsetParameter != "" {
		offset, err = strconv.Atoi(offsetParameter)
		if err != nil || offset < 0 {
			return 0, 0, errors.New("invalid query parameter: offset")
		}
	}

	if limitParameter != "" {
		limit, err = strconv.Atoi(limitParameter)
		if err != nil || limit < 0 {
			return 0, 0, errors.New("invalid query parameter: limit")
		}
	}

	logData["offset"] = offsetParameter
	logData["limit"] = limitParameter

	return
}
