package mrkl

import (
	"strings"

	"github.com/wejick/gchain/agent"
)

const (
	prefix = `You can use tools to get new information.
Answer the following questions as best you can using the following tools:

{{.tool_definition}}
`

	formatInstruction = `answer  strictly using the following structure:
Question : <string, the input question you must answer>
Thought : <string, you should always think about what to do>
ToolName : <string, the action to take, should be one of [ {{.tool_names}} ]>
ToolInputJson : <string, the input to the action as stringified json>
Observation: <string, the result of the action>
... (this Thought/Action/Action Input/Observation can repeat N times)
Thought: <string, I now know the final answer>
FinalAnswer: <string, the final answer to the original input question>
`

	suffix = `Begin!

Question: {{.input}}
Thought:{{.agent_scratchpad}}`
)

func parseInput(input string) (action agent.Action, err error) {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		split := strings.SplitN(line, ":", 2)
		// trim space
		for i, s := range split {
			split[i] = strings.TrimSpace(s)
		}
		if len(split) != 2 {
			continue
		}
		switch split[0] {
		case "Question":
			action.Question = split[1]
		case "Thought":
			action.Thought = split[1]
		case "ToolName":
			action.ToolName = split[1]
		case "ToolInputJson":
			action.ToolInputJson = split[1]
		case "Observation":
			action.ToolOutput = split[1]
		case "FinalAnswer":
			action.Message = split[1]
			if action.Message != "" {
				action.FinalAction = true
			}
		}
	}

	if !action.FinalAction && action.ToolName == "" {
		return action, agent.ErrParsingOutput
	} else if action.ToolName != "" && action.ToolInputJson == "" {
		return action, agent.ErrParsingOutput
	} else if action.FinalAction && action.Message == "" {
		return action, agent.ErrParsingOutput
	} else if action.FinalAction && action.ToolName != "" {
		return action, agent.ErrParsingOutput
	}

	return
}
