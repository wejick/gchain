package greetings

import (
	"context"
	"encoding/json"

	"github.com/wejick/gchain/model"
)

type GreetingsTool struct {
	functionDefinition model.FunctionDefinition
}

func NewGreetingsTool() *GreetingsTool {
	return &GreetingsTool{
		functionDefinition: model.FunctionDefinition{
			Name:        "greetings_tool",
			Description: "This tool is used to greet user with hello",
			Parameters: model.FunctionJsonSchema{
				Type: model.FunctionDataTypeObject,
				Properties: map[string]model.FunctionJsonSchema{
					"user_name": {
						Type:        model.FunctionDataTypeString,
						Description: "User name",
					},
				},
				Required: []string{"user_name"},
			},
		},
	}
}

// Run give greeting to user, this is to demonstrate the simples form of tool
// expect input inside prompt["input"] and return output inside output["output"]
// the input is a non nested json that will be marshalled into map[string]string and will be passed as parameter to the real greeting function
func (G *GreetingsTool) Run(ctx context.Context, prompt map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	stringOutput, err := G.SimpleRun(ctx, prompt["input"], options...)
	output = map[string]string{"output": stringOutput}
	return
}

// SimpleRun give greeting to user, this is to demonstrate the simples form of tool
func (G *GreetingsTool) SimpleRun(ctx context.Context, prompt string, options ...func(*model.Option)) (output string, err error) {
	var parameter map[string]string
	err = json.Unmarshal([]byte(prompt), &parameter)
	if err != nil {
		return
	}
	output = G.greetings(parameter["user_name"])
	return
}

func (G *GreetingsTool) greetings(username string) string {
	return "Hello " + username + " welcome to the paradise of the world"
}

// GetFunctionDefinition return function definition of the tool
func (G *GreetingsTool) GetFunctionDefinition() model.FunctionDefinition {
	return G.functionDefinition
}

// GetDefinitionString tool definition in string format
func (G *GreetingsTool) GetDefinitionString() string {
	description := `name = ` + G.functionDefinition.Name + `
description = ` + G.functionDefinition.Description + `
` + G.functionDefinition.Parameters.String()

	return description
}
