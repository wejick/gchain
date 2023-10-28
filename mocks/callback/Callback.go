// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"

	callback "github.com/wejick/gchain/callback"

	mock "github.com/stretchr/testify/mock"
)

// Callback is an autogenerated mock type for the Callback type
type Callback struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, data
func (_m *Callback) Execute(_a0 context.Context, data callback.CallbackData) {
	_m.Called(_a0, data)
}

// NewCallback creates a new instance of Callback. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCallback(t interface {
	mock.TestingT
	Cleanup(func())
}) *Callback {
	mock := &Callback{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
