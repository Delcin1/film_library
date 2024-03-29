// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UserSaver is an autogenerated mock type for the UserSaver type
type UserSaver struct {
	mock.Mock
}

// SaveUser provides a mock function with given fields: username, password
func (_m *UserSaver) SaveUser(username string, password string) error {
	ret := _m.Called(username, password)

	if len(ret) == 0 {
		panic("no return value specified for SaveUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserSaver creates a new instance of UserSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserSaver {
	mock := &UserSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
