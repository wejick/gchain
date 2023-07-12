package greeting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreetingsTool(t *testing.T) {
	greetingTool := NewGreetingTool()
	assert.NotNil(t, greetingTool)

	assert.NotNil(t, greetingTool.GetFunctionDefinition())

	// test simple run
	output, err := greetingTool.SimpleRun(context.Background(), `{"user_name":"John"}`)
	assert.Nil(t, err)
	assert.Equal(t, "Hello John welcome to the paradise of the world", output)

	// test simple run with invalid input
	output, err = greetingTool.SimpleRun(context.Background(), `123`)
	assert.NotNil(t, err)
	assert.Equal(t, "", output)

	// test run
	outputRun, err := greetingTool.Run(context.Background(), map[string]string{"user_name": "John"})
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"output": "Hello John welcome to the paradise of the world"}, outputRun)

	// test run with invalid input
	outputRun, err = greetingTool.Run(context.Background(), nil)
	assert.NotNil(t, err)
	assert.Nil(t, outputRun)

	// get description
	description := `name = ` + greetingTool.functionDefinition.Name + `
description = ` + greetingTool.functionDefinition.Description + "\n"
	assert.Equal(t, description+greetingTool.GetFunctionDefinition().Parameters.String(), greetingTool.GetDefinitionString())
}
