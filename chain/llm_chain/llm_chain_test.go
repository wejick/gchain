package llm_chain

import (
	"context"
	"reflect"
	"testing"

	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
)

func TestLLMChain_SimpleRun(t *testing.T) {
	type fields struct {
		llmModel model.LLMModel
	}
	type args struct {
		ctx     context.Context
		input   string
		options []func(*model.Option)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput string
		wantErr    bool
	}{
		{
			name: "empty",
			fields: fields{
				llmModel: &model.LLMModelMock{
					CallFunc: func(ctx context.Context, prompt string, options ...func(*model.Option)) (string, error) {
						return "", nil
					},
				},
			},
		},
		{
			name: "echo prompt",
			fields: fields{
				llmModel: &model.LLMModelMock{
					CallFunc: func(ctx context.Context, prompt string, options ...func(*model.Option)) (string, error) {
						return prompt, nil
					},
				},
			},
			args: args{
				input: "echo prompt",
			},
			wantOutput: "echo prompt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := &LLMChain{
				llmModel: tt.fields.llmModel,
			}
			gotOutput, err := L.SimpleRun(tt.args.ctx, tt.args.input, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("LLMChain.SimpleRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput != tt.wantOutput {
				t.Errorf("LLMChain.SimpleRun() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestLLMChain_Run(t *testing.T) {
	type fields struct {
		llmModel model.LLMModel
	}
	type args struct {
		ctx     context.Context
		input   map[string]string
		options []func(*model.Option)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput map[string]string
		wantErr    bool
	}{
		{
			name: "empty",
			fields: fields{
				llmModel: &model.LLMModelMock{
					CallFunc: func(ctx context.Context, prompt string, options ...func(*model.Option)) (string, error) {
						return "", nil
					},
				},
			},
			wantOutput: map[string]string{"output": ""},
		},
		{
			name: "echo input",
			fields: fields{
				llmModel: &model.LLMModelMock{
					CallFunc: func(ctx context.Context, prompt string, options ...func(*model.Option)) (string, error) {
						return prompt, nil
					},
				},
			},
			args: args{
				input: map[string]string{"input": "echo input"},
			},
			wantOutput: map[string]string{"output": "echo input"},
		},
	}
	customPrompt, _ := prompt.NewPromptTemplate("customPrompt", "{{.input}}")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := &LLMChain{
				llmModel:       tt.fields.llmModel,
				promptTemplate: customPrompt,
			}
			gotOutput, err := L.Run(tt.args.ctx, tt.args.input, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("LLMChain.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("LLMChain.Run() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
