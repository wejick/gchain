// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"

	agent "github.com/wejick/gchain/agent"

	mock "github.com/stretchr/testify/mock"
)

// BaseAgent is an autogenerated mock type for the BaseAgent type
type BaseAgent struct {
	mock.Mock
}

// Plan provides a mock function with given fields: ctx, userPrompt, actionTaken
func (_m *BaseAgent) Plan(ctx context.Context, userPrompt string, actionTaken []agent.Action) (agent.Action, error) {
	ret := _m.Called(ctx, userPrompt, actionTaken)

	var r0 agent.Action
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []agent.Action) (agent.Action, error)); ok {
		return rf(ctx, userPrompt, actionTaken)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []agent.Action) agent.Action); ok {
		r0 = rf(ctx, userPrompt, actionTaken)
	} else {
		r0 = ret.Get(0).(agent.Action)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []agent.Action) error); ok {
		r1 = rf(ctx, userPrompt, actionTaken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterToolDefinition provides a mock function with given fields: toolDefinition
func (_m *BaseAgent) RegisterToolDefinition(toolDefinition string) {
	_m.Called(toolDefinition)
}

// NewBaseAgent creates a new instance of BaseAgent. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBaseAgent(t interface {
	mock.TestingT
	Cleanup(func())
}) *BaseAgent {
	mock := &BaseAgent{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}