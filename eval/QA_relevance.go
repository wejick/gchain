package eval

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wejick/gchain/model"
	"github.com/wejick/gchain/prompt"
)

// QARelevanceEval is grading whether the the answer relevant to the question according to the given fact.
// Example :
// -> NewQARelevanceEval(model, "The color of the sky is blue","What is the color of the sky?").Evaluate("When I go outside, I see the sky is blue")
// -> True
type QARelevanceEval struct {
	llmModel           model.LLMModel
	evaluationTemplate *prompt.PromptTemplate
	fact               string
	question           string
}

func NewQARelevanceEval(model model.LLMModel, fact, question string) (evaluator *QARelevanceEval) {
	evaluationTemplate, _ := prompt.NewPromptTemplate("evaluation", QARelevanceInstruction)

	evaluator = &QARelevanceEval{
		llmModel:           model,
		question:           question,
		fact:               fact,
		evaluationTemplate: evaluationTemplate,
	}

	return
}

func (A *QARelevanceEval) Evaluate(answer string) (bool, error) {
	data := make(map[string]string)
	data["fact"] = A.fact
	data["question"] = A.question
	data["answer"] = answer

	prompt, err := A.evaluationTemplate.FormatPrompt(data)
	if err != nil {
		return false, err
	}
	output, err := A.llmModel.Call(context.Background(), prompt, model.WithMaxToken(1000))
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

const QARelevanceInstruction = `You are grading whether the provided answer is answering the question. If the answer is relevant, then the output passes the test. You respond with a JSON object with this structure: {pass: boolean; reason: string;}. Only return the JSON object.
Examples:
Fact: Today the weather is nice and the sky is blue
Question: What is the color of the sky?
Answer: The sky is blue
{"reason": "the sky is blue","pass": true}

Fact: Today the weather is gloomy and the sky is grey
Question: Do you love the weather?
Answer: The sky is blue
{"reason": "Question is not related to the answer","pass": false}

Fact: {{.fact}}
Question: {{.question}}
Answer: {{.answer}}`
