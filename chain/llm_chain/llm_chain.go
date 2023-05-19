package llm_chain

import (
	"context"
	"errors"

	"github.com/wejick/gochain/chain"
	"github.com/wejick/gochain/model"
)

var _ chain.BaseChain = &LLMChain{}

type LLMChain struct {
	llmModel model.LLMModel
}

func NewLLMChain(llmModel model.LLMModel) (llmchain *LLMChain) {
	return &LLMChain{
		llmModel: llmModel,
	}
}

// Run expect input["input"] as input, and put the result to output["output"]
func (L *LLMChain) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := input["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	output = make(map[string]string)

	output["output"], err = L.llmModel.Call(ctx, input["input"], options...)

	return
}

// SimpleRun will run the input string agains llmchain
func (L *LLMChain) SimpleRun(ctx context.Context, input string, options ...func(*model.Option)) (output string, err error) {
	output, err = L.llmModel.Call(ctx, input, options...)
	return
}
