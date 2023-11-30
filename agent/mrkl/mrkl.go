package mrkl

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/wejick/gchain/agent"
	"github.com/wejick/gchain/chain"
	"github.com/wejick/gchain/prompt"
)

// based on https://github.com/langchain-ai/langchain/tree/master/libs/langchain/langchain/agents/mrkl

type MRKLAgent struct {
	chain           chain.BaseChain
	agentPrompt     *prompt.PromptTemplate
	toolsName       []string
	toolsDefinition string
}

func NewMRKLAgent(chain chain.BaseChain) (*MRKLAgent, error) {
	agentPrompt, err := prompt.NewPromptTemplate("agentPrompt", prefix+formatInstruction+suffix)
	if err != nil {
		log.Println("error creating agent prompt template", err)
	}
	return &MRKLAgent{
		chain:       chain,
		agentPrompt: agentPrompt,
	}, err
}

func (M *MRKLAgent) Plan(ctx context.Context, userPrompt string, actionTaken []agent.Action) (plan agent.Action, err error) {
	agentPrompt, err := M.formatPrompt(userPrompt, actionTaken)
	if err != nil {
		log.Println("error formatting agent prompt", err)
		return
	}
	chainOutput, err := M.chain.SimpleRun(ctx, agentPrompt)
	if err != nil {
		log.Println("error running chain", err)
		return
	}

	plan, err = parseInput(chainOutput)
	if err != nil {
		return plan, errors.New("chain output is not JSON : " + err.Error())
	}

	return
}

func (M *MRKLAgent) formatPrompt(userPrompt string, actionTaken []agent.Action) (agentPrompt string, err error) {
	dataPrompt := map[string]string{
		"input":            userPrompt,
		"tool_definition":  M.toolsDefinition,
		"tool_names":       fmt.Sprintf("%v", M.toolsName),
		"agent_scratchpad": fmt.Sprintf("%v", actionTaken),
	}
	agentPrompt, err = M.agentPrompt.FormatPrompt(dataPrompt)
	if err != nil {
		log.Println("error formatting agent prompt", err)
	}

	return
}

func (M *MRKLAgent) RegisterToolDefinition(toolDefinition string) {
	M.toolsDefinition += "\n" + toolDefinition
}

func (M *MRKLAgent) RegisterToolName(toolName string) {
	M.toolsName = append(M.toolsName, toolName)
}
