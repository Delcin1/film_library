// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	postgres "film_library/internal/storage/postgres"

	mock "github.com/stretchr/testify/mock"
)

// MovieSearcherByPart is an autogenerated mock type for the MovieSearcherByPart type
type MovieSearcherByPart struct {
	mock.Mock
}

// GetMoviesBySearchRequest provides a mock function with given fields: searchRequest
func (_m *MovieSearcherByPart) GetMoviesBySearchRequest(searchRequest string) ([]postgres.Movie, error) {
	ret := _m.Called(searchRequest)

	if len(ret) == 0 {
		panic("no return value specified for GetMoviesBySearchRequest")
	}

	var r0 []postgres.Movie
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]postgres.Movie, error)); ok {
		return rf(searchRequest)
	}
	if rf, ok := ret.Get(0).(func(string) []postgres.Movie); ok {
		r0 = rf(searchRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]postgres.Movie)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(searchRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMovieSearcherByPart creates a new instance of MovieSearcherByPart. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMovieSearcherByPart(t interface {
	mock.TestingT
	Cleanup(func())
}) *MovieSearcherByPart {
	mock := &MovieSearcherByPart{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
