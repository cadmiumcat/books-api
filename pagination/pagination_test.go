package pagination

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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
			Items:      []int{2, 1},
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

	//Convey("Given a GET request with invalid limit parameter", t, func() {
	//	r := httptest.NewRequest("GET", "/test?limit=-1&offset=1", nil)
	//	w := httptest.NewRecorder()
	//
	//	handler := func(w http.ResponseWriter, r *http.Request, offset int, limit int) (interface{}, int, error) {
	//		return []int{offset, limit}, 0, errors.New("invalid query parameter: limit")
	//	}
	//
	//	paginator := Paginator{
	//		DefaultLimit:        10,
	//		DefaultOffset:       0,
	//		DefaultMaximumLimit: 100,
	//	}
	//
	//	Convey("When paginate is called", func() {
	//		paginateHandler := paginator.Paginate(handler)
	//		paginateHandler(w, r)
	//		Convey("Then the response code is 400", func() {
	//			So(w.Code, ShouldEqual, http.StatusBadRequest)
	//		})
	//	})
	//})
	//
	//Convey("Given a GET request with invalid limit parameter", t, func() {
	//	r := httptest.NewRequest("GET", "/test?limit=1&offset=-1", nil)
	//	w := httptest.NewRecorder()
	//
	//	handler := func(w http.ResponseWriter, r *http.Request, offset int, limit int) (interface{}, int, error) {
	//		return []int{offset, limit}, 0, errors.New("invalid query parameter: offset")
	//	}
	//
	//	paginator := Paginator{
	//		DefaultLimit:        10,
	//		DefaultOffset:       0,
	//		DefaultMaximumLimit: 100,
	//	}
	//
	//	Convey("When paginate is called", func() {
	//		paginateHandler := paginator.Paginate(handler)
	//		paginateHandler(w, r)
	//		Convey("Then the response code is 400", func() {
	//			So(w.Code, ShouldEqual, http.StatusBadRequest)
	//		})
	//	})
	//})
}

func Test_validateQueryParameters(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name           string
		args           args
		wantOffset     int
		wantLimit      int
		wantErr        bool
		wantErrMessage error
	}{
		{
			name:           "an invalid offset (negative)",
			args:           args{httptest.NewRequest("GET", "/test?limit=1&offset=-1", nil)},
			wantOffset:     0,
			wantLimit:      0,
			wantErr:        true,
			wantErrMessage: errors.New("invalid query parameter: offset"),
		},
		{
			name:           "an invalid limit (negative)",
			args:           args{httptest.NewRequest("GET", "/test?limit=-1&offset=1", nil)},
			wantOffset:     0,
			wantLimit:      0,
			wantErr:        true,
			wantErrMessage: errors.New("invalid query parameter: limit"),
		},
		{
			name:           "an invalid limit (not an integer)",
			args:           args{httptest.NewRequest("GET", "/test?limit=words&offset=1", nil)},
			wantOffset:     0,
			wantLimit:      0,
			wantErr:        true,
			wantErrMessage: errors.New("invalid query parameter: limit"),
		},
		{
			name:           "an invalid offset (not an integer)",
			args:           args{httptest.NewRequest("GET", "/test?limit=1&offset=words", nil)},
			wantOffset:     0,
			wantLimit:      0,
			wantErr:        true,
			wantErrMessage: errors.New("invalid query parameter: offset"),
		},
		{
			name:           "a limit which exceeds the DefaultMaximumLimit",
			args:           args{httptest.NewRequest("GET", "/test?limit=101&offset=1", nil)},
			wantOffset:     0,
			wantLimit:      0,
			wantErr:        true,
			wantErrMessage: errors.New("invalid query parameter: limit exceeds maximum limit allowed"),
		},
		{
			name:           "no limit/offset",
			args:           args{httptest.NewRequest("GET", "/test1", nil)},
			wantOffset:     1,
			wantLimit:      10,
			wantErr:        false,
			wantErrMessage: nil,
		},
	}

	Convey("Given a GET request and a Paginator", t, func() {
		paginator := &Paginator{
			DefaultLimit:        10,
			DefaultOffset:       1,
			DefaultMaximumLimit: 100,
		}
		for _, tt := range tests {
			Convey(fmt.Sprintf("When the query parameters are validated, and they contain %q", tt.name), func() {
				gotOffset, gotLimit, err := paginator.validateQueryParameters(tt.args.r)
				if tt.wantErr {
					Convey(fmt.Sprintf("Then the error matches %q", tt.wantErrMessage), func() {
						So(err, ShouldBeError, tt.wantErrMessage)
					})
				}
				Convey(fmt.Sprintf("And the offset is set to %v", tt.wantOffset), func() {
					So(gotOffset, ShouldEqual, tt.wantOffset)
				})
				Convey(fmt.Sprintf("And the limit is set to %v", tt.wantLimit), func() {
					So(gotLimit, ShouldEqual, tt.wantLimit)
				})
			})
		}
	})
}
