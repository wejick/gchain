package summarization

import (
	"context"
	"testing"

	"github.com/wejick/gochain/callback"
	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/model"
)

var echoLlmChain, _ = llm_chain.NewLLMChain(&model.LLMModelMock{
	CallFunc: func(ctx context.Context, prompt string, options ...func(*model.Option)) (string, error) {
		return prompt, nil
	},
}, callback.NewManager(), nil, false)
var testChain, _ = NewStuffSummarizationChain(echoLlmChain, "", "text")

func TestStuffSummarizationChain_SimpleRun(t *testing.T) {
	type args struct {
		ctx     context.Context
		input   string
		options []func(*model.Option)
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    bool
	}{
		{
			name: "empty",
			wantOutput: `Write a concise summary of the following:
"
"
CONCISE SUMMARY:`,
		},
		{
			name: "jakarta ramai sekali hari ini mamen",
			args: args{
				input: "jakarta ramai sekali hari ini mamen",
			},
			wantOutput: `Write a concise summary of the following:
"jakarta ramai sekali hari ini mamen
"
CONCISE SUMMARY:`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := testChain
			gotOutput, err := S.SimpleRun(tt.args.ctx, tt.args.input, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("StuffSummarizationChain.SimpleRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput != tt.wantOutput {
				t.Errorf("StuffSummarizationChain.SimpleRun() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
