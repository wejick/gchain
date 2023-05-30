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
	promptSummarizeMapReduce = `Write a concise summary of the following:
"{{.text}}""
CONCISE SUMMARY:
	`
)

type MapReduceSummarizationChain struct {
	mapReduceCombineDocument *combine_document.MapReduceCombineDocument
}

var _ chain.BaseChain = &MapReduceSummarizationChain{}

// NewMapReduceSummarizationChain create new map reduce summarization chain instance
// put empty "" string to use default prompt
// put 0 to use default maxToken
func NewMapReduceSummarizationChain(llmChain *llm_chain.LLMChain, mapPromptString string, reducePromptString string,
	promptTemplateKey string, maxToken int) (m *MapReduceSummarizationChain, err error) {

	var promptTemplateMap, promptTemplateReduce *prompt.PromptTemplate

	if mapPromptString == "" {
		promptTemplateMap, err = prompt.NewPromptTemplate("map", promptSummarizeMapReduce)
		if err != nil {
			return
		}
		promptTemplateKey = "text"
	}

	if reducePromptString == "" {
		promptTemplateReduce, err = prompt.NewPromptTemplate("map", promptSummarizeMapReduce)
		if err != nil {
			return
		}
	}

	if maxToken == 0 {
		maxToken = 1000
	}

	mapReduceCombineDocument := combine_document.NewMapReduceCombineDocument(promptTemplateMap,
		promptTemplateReduce, promptTemplateKey, llmChain, maxToken)
	m = &MapReduceSummarizationChain{
		mapReduceCombineDocument: mapReduceCombineDocument,
	}

	return
}

// Run expect input["input"] as input, and put the result to output["output"]
func (M *MapReduceSummarizationChain) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := input["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	output, err = M.mapReduceCombineDocument.Run(ctx, input, options...)
	return
}

// SimpleRun will run the input prompt string againts llmchain
func (M *MapReduceSummarizationChain) SimpleRun(ctx context.Context, input string, options ...func(*model.Option)) (output string, err error) {
	output, err = M.mapReduceCombineDocument.SimpleRun(ctx, input, options...)
	return
}
