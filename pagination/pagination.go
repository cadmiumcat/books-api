package pagination

import (
	"encoding/json"
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

		offsetParameter := r.URL.Query().Get("offset")
		limitParameter := r.URL.Query().Get("limit")

		offset, err := strconv.Atoi(offsetParameter)
		limit, err := strconv.Atoi(limitParameter)
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

		logData := log.Data{"offset": offset, "limit": limit}

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
