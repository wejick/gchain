// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "github.com/wejick/gochain/model"
)

// BaseChain is an autogenerated mock type for the BaseChain type
type BaseChain struct {
	mock.Mock
}

// Run provides a mock function with given fields: ctx, prompt, options
func (_m *BaseChain) Run(ctx context.Context, prompt map[string]string, options ...func(*model.Option)) (map[string]string, error) {
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
func (_m *BaseChain) SimpleRun(ctx context.Context, prompt string, options ...func(*model.Option)) (string, error) {
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

type mockConstructorTestingTNewBaseChain interface {
	mock.TestingT
	Cleanup(func())
}

// NewBaseChain creates a new instance of BaseChain. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBaseChain(t mockConstructorTestingTNewBaseChain) *BaseChain {
	mock := &BaseChain{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
