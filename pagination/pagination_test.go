package pagination

import (
	. "github.com/smartystreets/goconvey/convey"
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
