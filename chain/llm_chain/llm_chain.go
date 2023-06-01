package llm_chain

import (
	"context"

	"github.com/wejick/gochain/chain"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
)

var _ chain.BaseChain = &LLMChain{}

const defaultTemplate = `{{.input}}`

type LLMChain struct {
	llmModel       model.LLMModel
	promptTemplate *prompt.PromptTemplate
}

// NewLLMChain create an LLMChain instance
// if nil promptTemplate provided, the default one will be used
// default promptTemplate expect prompt["input"] as template key
func NewLLMChain(llmModel model.LLMModel, promptTemplate *prompt.PromptTemplate) (llmchain *LLMChain, err error) {
	if promptTemplate == nil {
		promptTemplate, err = prompt.NewPromptTemplate("default", defaultTemplate)
		if err != nil {
			return
		}
	}
	llmchain = &LLMChain{
		llmModel:       llmModel,
		promptTemplate: promptTemplate,
	}

	return
}

// Run do completion
// the default template expect prompt["input"] as input, and put the result to output["output"]
func (L *LLMChain) Run(ctx context.Context, prompt map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	output = make(map[string]string)

	promptStr, err := L.promptTemplate.FormatPrompt(prompt)
	if err != nil {
		return
	}

	output["output"], err = L.llmModel.Call(ctx, promptStr, options...)

	return
}

// SimpleRun will run the prompt string agains llmchain
func (L *LLMChain) SimpleRun(ctx context.Context, prompt string, options ...func(*model.Option)) (output string, err error) {
	output, err = L.llmModel.Call(ctx, prompt, options...)
	return
}
