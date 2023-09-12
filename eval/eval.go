/*
Package eval provide evaluator for gchain
especially to evaluate LLM output using LLM output

Use this package for integration test or evaluating llm output
*/
package eval

type Evaluator interface {
	Evaluate(input string) (bool, error)
}

type llmEvalOutput struct {
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}
