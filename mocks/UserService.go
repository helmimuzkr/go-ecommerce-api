// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	users "e-commerce-api/feature/users"

	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// Delete provides a mock function with given fields: token
func (_m *UserService) Delete(token interface{}) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Login provides a mock function with given fields: username, password
func (_m *UserService) Login(username string, password string) (string, users.Core, error) {
	ret := _m.Called(username, password)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 users.Core
	if rf, ok := ret.Get(1).(func(string, string) users.Core); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Get(1).(users.Core)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(username, password)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Profile provides a mock function with given fields: token
func (_m *UserService) Profile(token interface{}) (interface{}, error) {
	ret := _m.Called(token)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(interface{}) interface{}); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: newUser
func (_m *UserService) Register(newUser users.Core) (users.Core, error) {
	ret := _m.Called(newUser)

	var r0 users.Core
	if rf, ok := ret.Get(0).(func(users.Core) users.Core); ok {
		r0 = rf(newUser)
	} else {
		r0 = ret.Get(0).(users.Core)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(users.Core) error); ok {
		r1 = rf(newUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: token, updateData
func (_m *UserService) Update(token interface{}, updateData users.Core) (users.Core, error) {
	ret := _m.Called(token, updateData)

	var r0 users.Core
	if rf, ok := ret.Get(0).(func(interface{}, users.Core) users.Core); ok {
		r0 = rf(token, updateData)
	} else {
		r0 = ret.Get(0).(users.Core)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}, users.Core) error); ok {
		r1 = rf(token, updateData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserService interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserService(t mockConstructorTestingTNewUserService) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
