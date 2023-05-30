package combine_document

import (
	"context"
	"errors"

	"github.com/wejick/gochain/chain"
	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
)

var _ CombinedDocument = &StuffCombineDocument{}
var _ chain.BaseChain = &StuffCombineDocument{}

// StuffCombineDocument chain to feed text document to LLM with specified prompt
type StuffCombineDocument struct {
	prompt            *prompt.PromptTemplate
	llmChain          *llm_chain.LLMChain
	promptTemplateKey string
}

// NewStuffCombineDocument creates new instance of StuffCombineDocument
func NewStuffCombineDocument(prompt *prompt.PromptTemplate,
	templateKey string, llmChain *llm_chain.LLMChain) *StuffCombineDocument {
	return &StuffCombineDocument{
		prompt:            prompt,
		llmChain:          llmChain,
		promptTemplateKey: templateKey,
	}
}

// Combine concatenate the given document and then feed to LLM
func (S *StuffCombineDocument) Combine(ctx context.Context, docs []string, options ...func(*model.Option)) (output string, err error) {
	//concat all docs into 1 string
	var doc string
	for _, item := range docs {
		doc += item + "\n"
	}
	templateData := map[string]string{S.promptTemplateKey: doc}

	prompt, err := S.prompt.FormatPrompt(templateData)
	if err != nil {
		return
	}
	output, err = S.llmChain.SimpleRun(ctx, prompt)

	return
}

// Run expect input["input"] as input, and put the result to output["output"]
func (S *StuffCombineDocument) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := input["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	output = make(map[string]string)
	output["output"], err = S.Combine(ctx, []string{input["input"]})

	return
}

// SimpleRun will run the input string agains llmchain
func (S *StuffCombineDocument) SimpleRun(ctx context.Context, input string, options ...func(*model.Option)) (output string, err error) {
	output, err = S.Combine(ctx, []string{input})
	return
}
