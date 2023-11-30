package mrkl

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wejick/gchain/agent"
	mockChain "github.com/wejick/gchain/mocks/chain"
)

func TestMRKLAgent_Plan(t *testing.T) {
	// Create a new MRKLAgent
	chain := mockChain.NewBaseChain(t)
	mrklAgent, err := NewMRKLAgent(chain)
	assert.NoError(t, err)

	// Prepare test data
	ctx := context.Background()
	actionTaken := []agent.Action{{Message: "mock message", FinalAction: true}}

	// Test 1 Plan Final Action
	// Call the Plan method
	promptInput, err := mrklAgent.formatPrompt("test 1", actionTaken)
	assert.NoError(t, err, "Test 1 Plan Final Action: error formatting agent prompt")

	chain.On("SimpleRun", ctx, promptInput).Return(`{"Message":"mock message","FinalAction":true}`, nil)
	plan, err := mrklAgent.Plan(ctx, "test 1", actionTaken)
	assert.NoError(t, err, "Test 1 Plan Final Action: error running chain")

	expectedPlan := agent.Action{Message: "mock message", FinalAction: true}
	assert.Equal(t, expectedPlan, plan, "Test 1 Plan Final Action: plan is not as expected")

	// Test 2 Plan Not Final Action
	// Call the Plan method
	promptInput, err = mrklAgent.formatPrompt("test 2", actionTaken)
	assert.NoError(t, err, "Test 2 Plan Not Final Action: error formatting agent prompt")

	chain.On("SimpleRun", ctx, promptInput).Return(`{"Message":"mock message","FinalAction":false}`, nil)
	plan, err = mrklAgent.Plan(ctx, "test 2", actionTaken)
	assert.NoError(t, err, "Test 2 Plan Not Final Action: error running chain")

	expectedPlan = agent.Action{Message: "mock message", FinalAction: false}
	assert.Equal(t, expectedPlan, plan, "Test 2 Plan Not Final Action: plan is not as expected")

	// Test 3 Plan output is not JSON
	// Call the Plan method
	promptInput, err = mrklAgent.formatPrompt("test 3", actionTaken)
	assert.NoError(t, err, "Test 3 Plan Error: error formatting agent prompt")

	chain.On("SimpleRun", ctx, promptInput).Return(`"Message":"mock message","FinalAction":false"`, nil)
	plan, err = mrklAgent.Plan(ctx, "test 3", actionTaken)
	assert.Error(t, err, "Test 3 Plan Error: expected error")

	expectedPlan = agent.Action{Message: `"Message":"mock message","FinalAction":false"`, FinalAction: true}
	assert.Equal(t, expectedPlan, plan, "Test 3 Plan Error: plan is not as expected")

}

func Test_parseInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    agent.Action
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   "",
			want:    agent.Action{},
			wantErr: true,
		},
		{
			name: "valid input",
			input: `Question: What is your name?
					Thought: My name is John.
					ToolName: tool1
					ToolInputJson: {"key": "value"}
					ToolOutput: output
					Message: success
					FinalAction: false`,
			want: agent.Action{
				Question:      "What is your name?",
				Thought:       "My name is John.",
				ToolName:      "tool1",
				ToolInputJson: `{"key": "value"}`,
				ToolOutput:    "output",
				Message:       "success",
				FinalAction:   false,
			},
			wantErr: false,
		},
		{
			name: "invalid input : tool name with no tool input json",
			input: `Question: What is your name?
					Thought: My name is John.
					ToolName: tool1
					ToolOutput: output
					Message: success
					FinalAction: false`,
			want: agent.Action{
				Question:    "What is your name?",
				Thought:     "My name is John.",
				ToolName:    "tool1",
				ToolOutput:  "output",
				Message:     "success",
				FinalAction: false,
			},
			wantErr: true,
		},
		{
			name: "invalid input : final action with no message",
			input: `Question: What is your name?
					Thought: My name is John.
					ToolName: tool1
					ToolInputJson: {"key": "value"}
					ToolOutput: output
					FinalAction: true`,
			want: agent.Action{
				Question:      "What is your name?",
				Thought:       "My name is John.",
				ToolName:      "tool1",
				ToolInputJson: `{"key": "value"}`,
				ToolOutput:    "output",
				FinalAction:   true,
			},
			wantErr: true,
		},
		{
			name: "invalid input : not final action with no tool",
			input: `Question: What is your name?
					Thought: My name is John.
					ToolName:
					FinalAction: false`,
			want: agent.Action{
				Question:    "What is your name?",
				Thought:     "My name is John.",
				FinalAction: false,
			},
			wantErr: true,
		},
		{
			name: "invalid input : final action with tool",
			input: `Question: What is your name?
					Thought: My name is John.
					ToolName: tool1
					ToolInputJson: {"key": "value"}
					Message: success
					FinalAction: true`,
			want: agent.Action{
				Question:      "What is your name?",
				Thought:       "My name is John.",
				ToolName:      "tool1",
				ToolInputJson: `{"key": "value"}`,
				Message:       "success",
				FinalAction:   true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
