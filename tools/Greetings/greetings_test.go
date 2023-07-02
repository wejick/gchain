package greetings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreetingsTool(t *testing.T) {
	greetingTool := NewGreetingsTool()
	assert.NotNil(t, greetingTool)

	assert.NotNil(t, greetingTool.GetFunctionDefinition())

	// test run
	output, err := greetingTool.SimpleRun(context.Background(), `{"user_name":"John"}`)
	assert.Nil(t, err)
	assert.Equal(t, "Hello John welcome to the paradise of the world", output)

	// test run with invalid input
	output, err = greetingTool.SimpleRun(context.Background(), `{"user_name":123}`)
	assert.NotNil(t, err)
	assert.Equal(t, "", output)

	// test run
	outputRun, err := greetingTool.Run(context.Background(), map[string]string{"input": `{"user_name":"John"}`})
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"output": "Hello John welcome to the paradise of the world"}, outputRun)

	// get description
	description := `name = ` + greetingTool.functionDefinition.Name + `
description = ` + greetingTool.functionDefinition.Description + "\n"
	assert.Equal(t, description+greetingTool.GetFunctionDefinition().Parameters.String(), greetingTool.GetDefinitionString())
}
