package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	mockTool "github.com/wejick/gchain/mocks/tools"
	"github.com/wejick/gchain/model"
)

type mockAgent struct{}

func (m *mockAgent) Plan(ctx context.Context, input string, actions []Action) (plan Action, err error) {
	if input == "final" {
		return Action{finalAction: true, message: "final message"}, nil
	}
	if input == "error" {
		return Action{}, errors.New("error")
	}
	if input == "tool_error" {
		return Action{toolName: "mockTool", toolInputJson: "error"}, nil
	}

	return Action{toolName: "mockTool", toolInputJson: "input"}, nil
}

func TestExecutor_Run(t *testing.T) {
	tool := mockTool.BaseTool{}
	agent := mockAgent{}

	// Create a new Executor
	tool.On("GetFunctionDefinition").Return(model.FunctionDefinition{Name: "mockTool"})
	executor := NewExecutor(&agent, 0)
	executor.RegisterTool(&tool)

	// Test final plan
	tool.On("SimpleRun", context.Background(), "input").Return("output", nil)
	output, err := executor.Run(context.Background(), map[string]string{"input": "final"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output["output"] != "final message" {
		t.Errorf("Expected output to be 'final message', got %s", output["output"])
	}

	// Test error from agent
	_, err = executor.Run(context.Background(), map[string]string{"input": "error"})
	assert.Error(t, err, "Expected error, got nil")

	// Test error from tool
	tool.On("SimpleRun", context.Background(), "error").Return("", errors.New("error"))
	_, err = executor.Run(context.Background(), map[string]string{"input": "tool_error"})
	assert.Error(t, err, "Expected error, got nil")

	// Test max iteration
	_, err = executor.Run(context.Background(), map[string]string{"input": "input"})
	if err != ErrMaxIteration {
		t.Errorf("Expected error to be ErrMaxIteration, got %v", err)
	}
}
