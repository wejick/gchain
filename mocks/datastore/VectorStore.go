// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"

	datastore "github.com/wejick/gchain/datastore"
	document "github.com/wejick/gchain/document"

	mock "github.com/stretchr/testify/mock"
)

// VectorStore is an autogenerated mock type for the VectorStore type
type VectorStore struct {
	mock.Mock
}

// AddDocuments provides a mock function with given fields: ctx, indexName, documents
func (_m *VectorStore) AddDocuments(ctx context.Context, indexName string, documents []document.Document) ([]error, error) {
	ret := _m.Called(ctx, indexName, documents)

	var r0 []error
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []document.Document) ([]error, error)); ok {
		return rf(ctx, indexName, documents)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []document.Document) []error); ok {
		r0 = rf(ctx, indexName, documents)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]error)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []document.Document) error); ok {
		r1 = rf(ctx, indexName, documents)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddText provides a mock function with given fields: ctx, indexName, input
func (_m *VectorStore) AddText(ctx context.Context, indexName string, input string) error {
	ret := _m.Called(ctx, indexName, input)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, indexName, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteIndex provides a mock function with given fields: ctx, indexName
func (_m *VectorStore) DeleteIndex(ctx context.Context, indexName string) error {
	ret := _m.Called(ctx, indexName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, indexName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Search provides a mock function with given fields: ctx, indexName, query, options
func (_m *VectorStore) Search(ctx context.Context, indexName string, query string, options ...func(*datastore.Option)) ([]document.Document, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, indexName, query)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []document.Document
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...func(*datastore.Option)) ([]document.Document, error)); ok {
		return rf(ctx, indexName, query, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...func(*datastore.Option)) []document.Document); ok {
		r0 = rf(ctx, indexName, query, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]document.Document)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, ...func(*datastore.Option)) error); ok {
		r1 = rf(ctx, indexName, query, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchVector provides a mock function with given fields: ctx, indexName, vector, options
func (_m *VectorStore) SearchVector(ctx context.Context, indexName string, vector []float32, options ...func(*datastore.Option)) ([]document.Document, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, indexName, vector)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []document.Document
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []float32, ...func(*datastore.Option)) ([]document.Document, error)); ok {
		return rf(ctx, indexName, vector, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []float32, ...func(*datastore.Option)) []document.Document); ok {
		r0 = rf(ctx, indexName, vector, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]document.Document)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []float32, ...func(*datastore.Option)) error); ok {
		r1 = rf(ctx, indexName, vector, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewVectorStore creates a new instance of VectorStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewVectorStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *VectorStore {
	mock := &VectorStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
