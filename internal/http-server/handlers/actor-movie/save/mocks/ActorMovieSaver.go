// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ActorMovieSaver is an autogenerated mock type for the ActorMovieSaver type
type ActorMovieSaver struct {
	mock.Mock
}

// SaveActorMovie provides a mock function with given fields: movieId, actorsIds
func (_m *ActorMovieSaver) SaveActorMovie(movieId int, actorsIds []int) error {
	ret := _m.Called(movieId, actorsIds)

	if len(ret) == 0 {
		panic("no return value specified for SaveActorMovie")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []int) error); ok {
		r0 = rf(movieId, actorsIds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewActorMovieSaver creates a new instance of ActorMovieSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewActorMovieSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *ActorMovieSaver {
	mock := &ActorMovieSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
