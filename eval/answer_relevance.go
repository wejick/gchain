package eval

import (
	"github.com/wejick/gchain/model"
	"github.com/wejick/gchain/prompt"
)

// AnswerRelevanceEval is an evaluator to evaluate input againts question expectation
// From given answer, LLM will to generate question. This evaluator will evaluate whether the question is relevant to the expectation.
// Example :
// -> NewAnswerRelevanceEval(model, "The color of the sky is blue").Evaluate("When I go outside, I see the sky is blue")
// -> True
type AnswerRelevanceEval struct {
	llmModel            model.LLMModel
	evaluationTemplate  *prompt.PromptTemplate
	questionExpectation string
}
