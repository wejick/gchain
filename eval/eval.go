/*
Package eval provide evaluator for gchain
especially to evaluate LLM output using LLM output

Use this package for integration test or evaluating llm output
*/
package eval

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wejick/gchain/model"
	"github.com/wejick/gchain/prompt"
)

type Evaluator interface {
	Evaluate(input string) (bool, error)
}

// ValidJson will check if input is valid json
type ValidJson struct {
}

// NewValidJson return valid json evaluator
func NewValidJson() *ValidJson {
	return &ValidJson{}
}

// Evaluate will return true if input is valid json
func (V *ValidJson) Evaluate(input string) (bool, error) {
	var v interface{}
	err := json.Unmarshal([]byte(input), v)
	if err != nil {
		return false, err
	}
	return true, nil
}

type LLMEval struct {
	llmModel           model.LLMModel
	evaluationTemplate *prompt.PromptTemplate
	rubric             string
}

type LLMEvalOutput struct {
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

func NewLLMEval(llmModel model.LLMModel, rubric string) *LLMEval {
	evaluationTemplate, _ := prompt.NewPromptTemplate("evaluation", instruction)

	return &LLMEval{
		llmModel:           llmModel,
		rubric:             rubric,
		evaluationTemplate: evaluationTemplate,
	}
}

// Evaluate will return true if input match with llm evaluation
func (L *LLMEval) Evaluate(input string) (bool, error) {
	data := make(map[string]string)
	data["input"] = input
	data["rubric"] = L.rubric
	prompt, err := L.evaluationTemplate.FormatPrompt(data)
	if err != nil {
		return false, err
	}
	output, err := L.llmModel.Call(context.Background(), prompt, model.WithMaxToken(100))
	if err != nil {
		return false, err
	}

	evalOutput := LLMEvalOutput{}
	err = json.Unmarshal([]byte(output), &evalOutput)
	if err != nil {
		return false, err
	}

	if !evalOutput.Pass {
		return false, errors.New(evalOutput.Reason)
	}

	return true, nil
}
