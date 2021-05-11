package pagination

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPaginator(t *testing.T) {
	Convey("Given a set of pagination parameters", t, func() {
		expectedPaginator := &Paginator{
			DefaultLimit:        10,
			DefaultOffset:       2,
			DefaultMaximumLimit: 100,
		}
		Convey("When NewPaginator is called", func() {
			actualPaginator := NewPaginator(10, 2, 100)
			Convey("Then a Paginator structure is returned with the correct values", func() {
				So(actualPaginator, ShouldResemble, expectedPaginator)
			})
		})
	})
}

func TestPaginator_Paginate(t *testing.T) {
	Convey("Given a GET request and valid query parameters", t, func() {
		r := httptest.NewRequest("GET", "/test?offset=2&limit=1", nil)
		w := httptest.NewRecorder()
		handler := func(w http.ResponseWriter, r *http.Request, offset int, limit int) (interface{}, int, error) {
			return []int{offset, limit}, 10, nil
		}

		paginator := &Paginator{
			DefaultLimit:        10,
			DefaultOffset:       0,
			DefaultMaximumLimit: 100,
		}

		expectedPage := page{
			Items:      []int{2,1},
			Count:      2,
			Offset:     2,
			Limit:      1,
			TotalCount: 10,
		}


		Convey("When paginate is called", func() {
			paginatedHandler := paginator.Paginate(handler)
			paginatedHandler(w, r)
			Convey("Then the response code is 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
			Convey("Then the parameters are passed to the handler function", func() {
				content, _ := ioutil.ReadAll(w.Body)
				expectedContent, _ := json.Marshal(expectedPage)
				So(string(content), ShouldEqual, string(expectedContent))
			})
		})
	})
}
