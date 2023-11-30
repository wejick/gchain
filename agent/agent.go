package agent

import (
	"context"
	"errors"

	"github.com/wejick/gchain/tools"
)

var (
	ErrMaxIteration  = errors.New("max iteration reached")
	ErrParsingOutput = errors.New("error parsing output")
)

// Action define what to execute and also containing the result
type Action struct {
	Question      string
	Thought       string
	ToolName      string // tool name to run by executor
	ToolInputJson string // input for the tool in json
	ToolOutput    string // output from tool
	Message       string // message from agent llm prediction
	FinalAction   bool   // is it final action or not
}

// BaseAgent
type BaseAgent interface {
	// Plan is determine what action to text next
	Plan(ctx context.Context, userPrompt string, actionTaken []Action) (plan Action, err error)
	RegisterToolDefinition(toolDefinition string)
	RegisterToolName(toolName string)
}

// BaseExecutor
type BaseExecutor interface {
	Run(ctx context.Context, input map[string]string) (output map[string]string, err error)
	RegisterTool(tool *tools.BaseTool)
}

// Executor
type Executor struct {
	agent        BaseAgent
	tools        map[string]tools.BaseTool
	maxIteration int
}

// NewExecutor create new executor
func NewExecutor(agent BaseAgent, maxIteration int) *Executor {
	if maxIteration == 0 {
		maxIteration = 10
	}
	return &Executor{
		agent:        agent,
		tools:        map[string]tools.BaseTool{},
		maxIteration: maxIteration,
	}
}

// RegisterTool put tool to executor
// this action also register the tool definition to agent
func (E *Executor) RegisterTool(tool tools.BaseTool) {
	E.tools[tool.GetFunctionDefinition().Name] = tool
	E.agent.RegisterToolDefinition(tool.GetDefinitionString())
	E.agent.RegisterToolName(tool.GetFunctionDefinition().Name)
}

// Run the executor loop until it reach final answer/action
func (E *Executor) Run(ctx context.Context, input map[string]string) (output map[string]string, err error) {
	actionTaken := []Action{}
	iterationNumber := 0
	for {
		// get plan from agent
		plan, err := E.agent.Plan(ctx, input["input"], actionTaken)
		if err != nil {
			return nil, err
		}

		// if plan is the final result, return it
		if plan.FinalAction {
			output = map[string]string{"output": plan.Message}
			break
		}

		// run plan / execute tool
		if plan.ToolName != "" {
			toolsOutput, err := E.tools[plan.ToolName].SimpleRun(ctx, plan.ToolInputJson)
			if err != nil {
				return nil, err
			}
			plan.ToolOutput = toolsOutput
		}

		// append action taken
		actionTaken = append(actionTaken, plan)

		if iterationNumber == E.maxIteration {
			return map[string]string{"output": plan.ToolOutput}, ErrMaxIteration
		}
		iterationNumber++
	}
	return
}
