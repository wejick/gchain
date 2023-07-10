package greeting

import (
	"context"
	"encoding/json"

	"github.com/wejick/gchain/model"
	"github.com/wejick/gchain/tools"
)

type GreetingTool struct {
	functionDefinition model.FunctionDefinition
}

func NewGreetingTool() *GreetingTool {
	return &GreetingTool{
		functionDefinition: model.FunctionDefinition{
			Name:        "greeting_tool",
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
func (G *GreetingTool) Run(ctx context.Context, prompt map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	stringOutput, err := G.SimpleRun(ctx, prompt["user_name"], options...)
	output = map[string]string{"output": stringOutput}
	return
}

// SimpleRun give greeting to user, this is to demonstrate the simples form of tool
func (G *GreetingTool) SimpleRun(ctx context.Context, prompt string, options ...func(*model.Option)) (output string, err error) {
	var parameter map[string]string
	err = json.Unmarshal([]byte(prompt), &parameter)
	if err != nil {
		return
	}
	output = G.greetings(parameter["user_name"])
	return
}

func (G *GreetingTool) greetings(username string) string {
	return "Hello " + username + " welcome to the paradise of the world"
}

// GetFunctionDefinition return function definition of the tool
func (G *GreetingTool) GetFunctionDefinition() model.FunctionDefinition {
	return G.functionDefinition
}

// GetDefinitionString tool definition in string format
func (G *GreetingTool) GetDefinitionString() string {
	description := tools.GetDefinitionString(G)

	return description
}
