package combine_document

import (
	"context"
	"errors"
	"strings"

	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
	"github.com/wejick/gochain/textsplitter"
)

// MapReduceCombineDocument Combining documents by mapping a chain over them first, then combining results.
type MapReduceCombineDocument struct {
	mapPrompt         *prompt.PromptTemplate
	reducePrompt      *prompt.PromptTemplate
	llmChain          *llm_chain.LLMChain
	splitter          textsplitter.TextSplitter
	promptTemplateKey string
	maxDocLength      int
}

// NewMapReduceCombineDocument creates new instance of MapReduceCombineDocument
func NewMapReduceCombineDocument(mapPrompt *prompt.PromptTemplate, reducePrompt *prompt.PromptTemplate,
	promptTemplateKey string, llmChain *llm_chain.LLMChain, splitter textsplitter.TextSplitter, maxDocLength int) *MapReduceCombineDocument {
	if maxDocLength == 0 {
		maxDocLength = 1000
	}
	return &MapReduceCombineDocument{
		mapPrompt:         mapPrompt,
		reducePrompt:      reducePrompt,
		llmChain:          llmChain,
		promptTemplateKey: promptTemplateKey,
		splitter:          splitter,
		maxDocLength:      maxDocLength,
	}
}

// Combine run the mapreduce process
func (M *MapReduceCombineDocument) Combine(ctx context.Context, docs []string, options ...func(*model.Option)) (output string, err error) {
	/*
		Map
		for each doc in Docs
		We're creating batch of map process, the size of batch <= maxToken. The batch contains M.mapPrompt + doc
		If Batch > maxToken. Create a new batch of the remaining text in the doc
		For each batch, we run M.llmChain.SimpleRun(ctx, prompt, options...), where prompt = SmapPrompt+doc

		Reduce
		when every batch finished, we compile all the result into 1 doc
		if M.reducePrompt + doc > maxToken. Create a new batch of the remaining text in the doc
		For each batch, we run M.llmChain.SimpleRun(ctx, prompt, options...), where prompt = SmapPrompt+doc
		do the reduce step again
	*/
	// Store intermediate results

	var intermediateResults []string

	// Map
	for _, doc := range docs {
		// split document into batches based on maxToken limit
		batches := M.splitter.SplitText(doc, M.maxDocLength, 0)

		mapResults, err := M.runOperation(ctx, batches, M.mapPrompt, options...)
		if err != nil {
			return "", err
		}
		intermediateResults = append(intermediateResults, mapResults...)
	}

	// Combine intermediate results with newline character
	intermediateString := strings.Join(intermediateResults, "\n")

	// Split combined result into batches again
	batches := M.splitter.SplitText(intermediateString, M.maxDocLength, 0)

	// Reduce
	finalResults, err := M.runOperation(ctx, batches, M.reducePrompt, options...)
	if err != nil {
		return "", err
	}

	output = strings.Join(finalResults, "")
	return
}

func (M *MapReduceCombineDocument) runOperation(ctx context.Context, batches []string, promptTemplate *prompt.PromptTemplate, options ...func(*model.Option)) ([]string, error) {
	var results []string
	for _, batch := range batches {
		// prepare data for prompt formatting
		data := map[string]string{"text": batch}
		prompt, err := promptTemplate.FormatPrompt(data)
		if err != nil {
			return nil, err
		}

		// run operation
		result, err := M.llmChain.SimpleRun(ctx, prompt, options...)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}
	return results, nil
}

// Run all entries in input map will be treated as document to be combined
// output will be output["output"]
func (M *MapReduceCombineDocument) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := input["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	output = make(map[string]string)
	output["output"], err = M.Combine(ctx, []string{input["input"]})

	return output, err
}

// SimpleRun will run the input string agains llmchain
func (M *MapReduceCombineDocument) SimpleRun(ctx context.Context, input string, options ...func(*model.Option)) (output string, err error) {
	err = errors.New("MapReduceCombineDocument.SimpleRun Method is Not Implemented")
	return
}
