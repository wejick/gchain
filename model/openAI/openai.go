package _openai

import (
	"context"

	goopenai "github.com/sashabaranov/go-openai"
	model "github.com/wejick/gochain/model"
)

type OpenAIModel struct {
	c         *goopenai.Client
	modelName string
}

// NewOpenAIModel return new openAI Model instance
func NewOpenAIModel(authToken string, orgID string, modelName string) (model *OpenAIModel) {
	var client *goopenai.Client
	if orgID != "" {
		client = goopenai.NewClient(authToken)
	} else {
		config := goopenai.DefaultConfig(authToken)
		config.OrgID = orgID
		client = goopenai.NewClientWithConfig(config)
	}

	model = &OpenAIModel{
		c:         client,
		modelName: modelName,
	}

	return
}

// Call runs completion request to the specified model
func (O *OpenAIModel) Call(ctx context.Context, prompt string, options ...func(*model.Option)) (output string, err error) {
	opts := model.Option{}
	for _, opt := range options {
		opt(&opts)
	}

	request := goopenai.CompletionRequest{
		Model:       O.modelName,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxToken,
		Prompt:      prompt,
	}

	response, err := O.c.CreateCompletion(ctx, request)
	if err != nil {
		return
	} else if len(response.Choices) > 0 {
		output = response.Choices[0].Text
	}

	return
}
