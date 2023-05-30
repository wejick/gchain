package summarization

import (
	"context"
	"errors"

	"github.com/wejick/gochain/chain"
	"github.com/wejick/gochain/chain/combine_document"
	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
)

const (
	promptSummarizeStuff = `Write a concise summary of the following:
"{{.text}}"
CONCISE SUMMARY:`
)

type StuffSummarizationChain struct {
	stuffCombineDocument *combine_document.StuffCombineDocument
}

var _ chain.BaseChain = &StuffSummarizationChain{}

func NewStuffSummarizationChain(llm_chain *llm_chain.LLMChain,
	promptTemplateString string, promptTemplateKey string) (s *StuffSummarizationChain, err error) {

	var promptTemplate *prompt.PromptTemplate

	if promptTemplateString == "" {
		promptTemplate, err = prompt.NewPromptTemplate("stuff", promptSummarizeStuff)
		if err != nil {
			return
		}
		promptTemplateKey = "text"
	}

	stuffCombineDocument := combine_document.NewStuffCombineDocument(promptTemplate, promptTemplateKey, llm_chain)
	s = &StuffSummarizationChain{
		stuffCombineDocument: stuffCombineDocument,
	}

	return
}

// Run all entries in input map will be treated as document to be combined
// output will be output["output"]
func (S *StuffSummarizationChain) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := input["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	output, err = S.stuffCombineDocument.Run(ctx, input, options...)
	return
}

// SimpleRun will run the input prompt string againts llmchain
func (S *StuffSummarizationChain) SimpleRun(ctx context.Context, input string, options ...func(*model.Option)) (output string, err error) {
	output, err = S.stuffCombineDocument.SimpleRun(ctx, input, options...)
	return
}
