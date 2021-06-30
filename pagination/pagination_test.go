package pagination

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

const (
	defaultLimit        = 10
	defaultOffset       = 1
	defaultMaximumLimit = 100
)

func TestReadPaginationValues(t *testing.T) {
	defaultPaginator := NewPaginator(defaultLimit, defaultOffset, defaultMaximumLimit)

	Convey("Given a request without pagination parameters", t, func() {
		r := httptest.NewRequest("GET", "/endpoint_to_paginate", nil)
		Convey("When GetPaginationValues is called", func() {

			offset, limit, err := defaultPaginator.GetPaginationValues(r)
			Convey("Then the default values are returned", func() {
				So(err, ShouldBeNil)
				So(limit, ShouldEqual, defaultLimit)
				So(offset, ShouldEqual, defaultOffset)
			})
		})
	})

	Convey("Given a request with valid pagination parameters", t, func() {
		r := httptest.NewRequest("GET", "/endpoint_to_paginate?limit=13&offset=2", nil)
		Convey("When GetPaginationValues is called", func() {
			offset, limit, err := defaultPaginator.GetPaginationValues(r)
			Convey("Then the default values are overwritten by the ones in the request", func() {
				So(err, ShouldBeNil)
				So(limit, ShouldEqual, 13)
				So(offset, ShouldEqual, 2)
			})
		})
	})

	Convey("Given a request with a negative limit value", t, func() {
		r := httptest.NewRequest("GET", "/endpoint_to_paginate?limit=-13&offset=2", nil)
		Convey("When GetPaginationValues is called", func() {
			offset, limit, err := defaultPaginator.GetPaginationValues(r)
			Convey("Then the an error is returned saying the limit value is invalid", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, ErrInvalidLimitParameter)
				So(limit, ShouldEqual, 0)
				So(offset, ShouldEqual, 0)
			})
		})
	})

	Convey("Given a request with a non-numerical limit value", t, func() {
		r := httptest.NewRequest("GET", "/endpoint_to_paginate?limit=two&offset=2", nil)
		Convey("When GetPaginationValues is called", func() {
			offset, limit, err := defaultPaginator.GetPaginationValues(r)
			Convey("Then the an error is returned saying the limit value is invalid", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, ErrInvalidLimitParameter)
				So(limit, ShouldEqual, 0)
				So(offset, ShouldEqual, 0)
			})
		})
	})

	Convey("Given a request with a negative offset value", t, func() {
		r := httptest.NewRequest("GET", "/endpoint_to_paginate?limit=13&offset=-2", nil)
		Convey("When GetPaginationValues is called", func() {
			offset, limit, err := defaultPaginator.GetPaginationValues(r)
			Convey("Then the an error is returned saying the offset value is invalid", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, ErrInvalidOffsetParameter)
				So(limit, ShouldEqual, 0)
				So(offset, ShouldEqual, 0)
			})
		})
	})

	Convey("Given a request with a non-numerical offset value", t, func() {
		r := httptest.NewRequest("GET", "/endpoint_to_paginate?limit=13&offset=two", nil)
		Convey("When GetPaginationValues is called", func() {
			offset, limit, err := defaultPaginator.GetPaginationValues(r)
			Convey("Then the an error is returned saying the offset value is invalid", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, ErrInvalidOffsetParameter)
				So(limit, ShouldEqual, 0)
				So(offset, ShouldEqual, 0)
			})
		})
	})

	Convey("Given a request with a limit value that exceeds the default maximum limit", t, func() {
		r := httptest.NewRequest("GET", "/endpoint_to_paginate?limit=101&offset=2", nil)
		Convey("When GetPaginationValues is called", func() {
			offset, limit, err := defaultPaginator.GetPaginationValues(r)
			Convey("Then the an error is returned saying the limit value is invalid and is over the maximum value", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, ErrLimitOverMax)
				So(limit, ShouldEqual, 0)
				So(offset, ShouldEqual, 0)
			})
		})
	})

}

func TestNewPaginator(t *testing.T) {
	Convey("Given an expectedPaginator with a set of values", t, func() {
		expectedPaginator := &Paginator{
			DefaultLimit:        defaultLimit,
			DefaultOffset:       defaultOffset,
			DefaultMaximumLimit: defaultMaximumLimit,
		}

		Convey("When NewPaginator is called using the same values", func() {
			actualPaginator := NewPaginator(defaultLimit, defaultOffset, defaultMaximumLimit)
			Convey("Then the Paginator returned resembles the expectedPaginator", func() {
				So(actualPaginator, ShouldResemble, expectedPaginator)
			})
		})
	})
}
