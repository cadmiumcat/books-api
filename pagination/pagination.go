package pagination

import (
	"encoding/json"
	"github.com/ONSdigital/log.go/log"
	"net/http"
)

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

		offset := p.DefaultOffset
		limit := p.DefaultLimit
		list, _, err := handler(w, r, offset, limit)

		logData := log.Data{"offset": offset, "limit": limit}

		b, err := json.Marshal(list)

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
