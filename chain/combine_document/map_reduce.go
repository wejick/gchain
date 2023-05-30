package combine_document

import (
	"context"
	"errors"
	"strings"

	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
)

// MapReduceCombineDocument Combining documents by mapping a chain over them first, then combining results.
type MapReduceCombineDocument struct {
	mapPrompt         *prompt.PromptTemplate
	reducePrompt      *prompt.PromptTemplate
	llmChain          *llm_chain.LLMChain
	promptTemplateKey string
	maxDocLength      int
}

// NewMapReduceCombineDocument creates new instance of MapReduceCombineDocument
func NewMapReduceCombineDocument(mapPrompt *prompt.PromptTemplate, reducePrompt *prompt.PromptTemplate,
	promptTemplateKey string, llmChain *llm_chain.LLMChain, maxDocLength int) *MapReduceCombineDocument {
	if maxDocLength == 0 {
		maxDocLength = 1000
	}
	return &MapReduceCombineDocument{
		mapPrompt:         mapPrompt,
		reducePrompt:      reducePrompt,
		llmChain:          llmChain,
		promptTemplateKey: promptTemplateKey,
		maxDocLength:      maxDocLength,
	}
}

// Combine run the mapreduce process
func (S *MapReduceCombineDocument) Combine(ctx context.Context, docs []string, options ...func(*model.Option)) (output string, err error) {
	/*
		Map
		for each doc in Docs
		We're creating batch of map process, the size of batch <= maxToken. The batch contains S.mapPrompt + doc
		If Batch > maxToken. Create a new batch of the remaining text in the doc
		For each batch, we run S.llmChain.SimpleRun(ctx, prompt, options...), where prompt = SmapPrompt+doc

		Reduce
		when every batch finished, we compile all the result into 1 doc
		if s.reducePrompt + doc > maxToken. Create a new batch of the remaining text in the doc
		For each batch, we run S.llmChain.SimpleRun(ctx, prompt, options...), where prompt = SmapPrompt+doc
		do the reduce step again
	*/
	// Store intermediate results

	var intermediateResults []string

	// Map
	for _, doc := range docs {
		// split document into batches based on maxToken limit
		batches := splitIntoBatches(doc, S.maxDocLength)

		mapResults, err := S.runOperation(ctx, batches, S.mapPrompt, options...)
		if err != nil {
			return "", err
		}
		intermediateResults = append(intermediateResults, mapResults...)
	}

	// Combine intermediate results with newline character
	intermediateString := strings.Join(intermediateResults, "\n")

	// Split combined result into batches again
	batches := splitIntoBatches(intermediateString, S.maxDocLength)

	// Reduce
	finalResults, err := S.runOperation(ctx, batches, S.reducePrompt, options...)
	if err != nil {
		return "", err
	}

	output = strings.Join(finalResults, "")
	return
}

func (S *MapReduceCombineDocument) runOperation(ctx context.Context, batches []string, promptTemplate *prompt.PromptTemplate, options ...func(*model.Option)) ([]string, error) {
	var results []string
	for _, batch := range batches {
		// prepare data for prompt formatting
		data := map[string]string{"text": batch}
		prompt, err := promptTemplate.FormatPrompt(data)
		if err != nil {
			return nil, err
		}

		// run operation
		result, err := S.llmChain.SimpleRun(ctx, prompt, options...)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}
	return results, nil
}

// splitIntoBatches creates word batches where length's doesn't exceed maxToken.
func splitIntoBatches(input string, maxToken int) []string {
	var batches []string

	words := strings.Fields(input)
	var batch []string
	var lenCounter int

	for _, word := range words {
		// +1 is for a possible space character
		if lenCounter+len(word)+1 > maxToken {
			batches = append(batches, strings.Join(batch, " "))
			batch = []string{}
			lenCounter = 0
		}

		batch = append(batch, word)
		lenCounter += len(word) + 1
	}

	if len(batch) > 0 {
		batches = append(batches, strings.Join(batch, " "))
	}

	return batches
}

// Run all entries in input map will be treated as document to be combined
// output will be output["output"]
func (S *MapReduceCombineDocument) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := input["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	output = make(map[string]string)
	output["output"], err = S.Combine(ctx, []string{input["input"]})

	return output, err
}

// SimpleRun will run the input string agains llmchain
func (S *MapReduceCombineDocument) SimpleRun(ctx context.Context, input string, options ...func(*model.Option)) (output string, err error) {
	err = errors.New("MapReduceCombineDocument.SimpleRun Method is Not Implemented")
	return
}
