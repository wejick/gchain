package tools

import (
	"fmt"

	"github.com/wejick/gchain/chain"
	"github.com/wejick/gchain/model"
)

// BaseTool is the interface for tool
// the idea is to keep it compatible with chain interface, so chain be used as tool as well
type BaseTool interface {
	chain.BaseChain
	GetFunctionDefinition() model.FunctionDefinition // Get tools definition in the form of function definition
	GetDefinitionString() string                     // Get tools definition in the form of text description
}

const toolDefinitionString = "name = %s\ndescription = %s\n%s"

// GetDefinitionString return tool definition in string format
func GetDefinitionString(t BaseTool) string {
	output := fmt.Sprintf(toolDefinitionString, t.GetFunctionDefinition().Name, t.GetFunctionDefinition().Description, t.GetFunctionDefinition().Parameters.String())
	return output
}
