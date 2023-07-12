package _openai

import (
	"context"

	goopenai "github.com/sashabaranov/go-openai"
	"github.com/wejick/gchain/callback"
	model "github.com/wejick/gchain/model"
)

type OpenAIModel struct {
	c               *goopenai.Client
	modelName       string
	callbackManager *callback.Manager
}

// NewOpenAIModel return new openAI Model instance
func NewOpenAIModel(authToken string, orgID string, baseURL string, modelName string, callbackManager *callback.Manager, verbose bool) (llm *OpenAIModel) {
	var client *goopenai.Client
	if baseURL == "" {
		client = goopenai.NewClient(authToken)
	} else {
		config := goopenai.DefaultAzureConfig(authToken, baseURL)
		config.OrgID = orgID
		client = goopenai.NewClientWithConfig(config)
	}

	llm = &OpenAIModel{
		c:               client,
		modelName:       modelName,
		callbackManager: callbackManager,
	}

	if verbose {
		llm.callbackManager.RegisterCallback(model.CallbackModelEnd, callback.VerboseCallback)
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

	// Trigger start and end callback
	O.callbackManager.TriggerEvent(ctx, model.CallbackModelStart, callback.CallbackData{
		EventName:    model.CallbackModelStart,
		FunctionName: "OpenAIModel.Call",
		Input:        map[string]string{"input": prompt},
		Output:       map[string]string{"output": output},
	})

	response, err := O.c.CreateCompletion(ctx, request)
	if err != nil {
		return
	} else if len(response.Choices) > 0 {
		output = response.Choices[0].Text
	}

	O.callbackManager.TriggerEvent(ctx, model.CallbackModelEnd, callback.CallbackData{
		EventName:    model.CallbackModelEnd,
		FunctionName: "OpenAIModel.Call",
		Input:        map[string]string{"input": prompt},
		Output:       map[string]string{"output": output},
	})

	return
}
