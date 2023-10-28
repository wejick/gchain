// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "github.com/wejick/gchain/model"
)

// BaseTool is an autogenerated mock type for the BaseTool type
type BaseTool struct {
	mock.Mock
}

// GetFunctionDefinition provides a mock function with given fields:
func (_m *BaseTool) GetFunctionDefinition() model.FunctionDefinition {
	ret := _m.Called()

	var r0 model.FunctionDefinition
	if rf, ok := ret.Get(0).(func() model.FunctionDefinition); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.FunctionDefinition)
	}

	return r0
}

// GetToolDescription provides a mock function with given fields:
func (_m *BaseTool) GetToolDescription() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Run provides a mock function with given fields: ctx, prompt, options
func (_m *BaseTool) Run(ctx context.Context, prompt map[string]string, options ...func(*model.Option)) (map[string]string, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, prompt)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string, ...func(*model.Option)) (map[string]string, error)); ok {
		return rf(ctx, prompt, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string, ...func(*model.Option)) map[string]string); ok {
		r0 = rf(ctx, prompt, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string]string, ...func(*model.Option)) error); ok {
		r1 = rf(ctx, prompt, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SimpleRun provides a mock function with given fields: ctx, prompt, options
func (_m *BaseTool) SimpleRun(ctx context.Context, prompt string, options ...func(*model.Option)) (string, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, prompt)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...func(*model.Option)) (string, error)); ok {
		return rf(ctx, prompt, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...func(*model.Option)) string); ok {
		r0 = rf(ctx, prompt, options...)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...func(*model.Option)) error); ok {
		r1 = rf(ctx, prompt, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBaseTool creates a new instance of BaseTool. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBaseTool(t interface {
	mock.TestingT
	Cleanup(func())
}) *BaseTool {
	mock := &BaseTool{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
