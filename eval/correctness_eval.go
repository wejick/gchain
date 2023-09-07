package eval

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wejick/gchain/model"
	"github.com/wejick/gchain/prompt"
)

// CorrectnessEval is an evaluator to evaluate input againts expectation
// Example :
// -> NewCorrectnessEval(model, "The color of the sky is blue").Evaluate("When I go outside, I see the sky is blue")
// -> True
type CorrectnessEval struct {
	llmModel           model.LLMModel
	evaluationTemplate *prompt.PromptTemplate
	expectation        string
}

type llmEvalOutput struct {
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

// NewCorrectnessEval create new correctness evaluator
func NewCorrectnessEval(llmModel model.LLMModel, expectation string) *CorrectnessEval {
	evaluationTemplate, _ := prompt.NewPromptTemplate("evaluation", correctnessEvalInstruction)

	return &CorrectnessEval{
		llmModel:           llmModel,
		expectation:        expectation,
		evaluationTemplate: evaluationTemplate,
	}
}

// The prompt based on
// https://github.com/axilla-io/ax/blob/main/packages/axeval/src/prompt.ts
var correctnessEvalInstruction = `You are grading output according to a user-specified rubric. If the statement in the rubric is true, then the output passes the test. You respond with a JSON object with this structure: {pass: boolean; reason: string;}. Only return the JSON object.
Examples:

Input: Hello world
Rubric: Content contains a greeting
{"reason": "the content contains the word 'world'","pass": true}

Input: Avast ye swabs, repel the invaders!
Rubric: Does not speak like a pirate
{"reason": "'avast ye' is a common pirate term","pass": false}

Input: {{.input}}
Rubric: {{.rubric}}`

// Evaluate will return true if input match with llm evaluation
func (L *CorrectnessEval) Evaluate(input string) (bool, error) {
	data := make(map[string]string)
	data["input"] = input
	data["rubric"] = L.expectation
	prompt, err := L.evaluationTemplate.FormatPrompt(data)
	if err != nil {
		return false, err
	}
	output, err := L.llmModel.Call(context.Background(), prompt, model.WithMaxToken(100))
	if err != nil {
		return false, err
	}

	evalOutput := llmEvalOutput{}
	err = json.Unmarshal([]byte(output), &evalOutput)
	if err != nil {
		return false, err
	}

	if !evalOutput.Pass {
		return false, errors.New(evalOutput.Reason)
	}

	return true, nil
}
