package tools

import (
	"github.com/wejick/gchain/chain"
	"github.com/wejick/gchain/model"
)

// BaseTool is the interface for tool
// the idea is to keep it compatible with chain interface, so chain be used as tool as well
type BaseTool interface {
	chain.BaseChain
	GetFunctionDefinition() model.FunctionDefinition // Get tools definition in the form of function definition
	GetToolDescription() string                      // Get tools definition in the form of text description
}
