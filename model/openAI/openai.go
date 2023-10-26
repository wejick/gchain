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

type OpenAIOption struct {
	OrgID      string
	BaseURL    string
	APIVersion string

	Verbose bool
}

func WithOrgID(OrgID string) func(*OpenAIOption) {
	return func(o *OpenAIOption) {
		o.OrgID = OrgID
	}
}

func WithBaseURL(baseURL string) func(*OpenAIOption) {
	return func(o *OpenAIOption) {
		o.BaseURL = baseURL
	}
}

func WithAPIVersion(apiVersion string) func(*OpenAIOption) {
	return func(o *OpenAIOption) {
		o.APIVersion = apiVersion
	}
}

func WithVerbose(verbose bool) func(*OpenAIOption) {
	return func(o *OpenAIOption) {
		o.Verbose = verbose
	}
}

func newOpenAIClient(authToken string, opts OpenAIOption) (client *goopenai.Client) {
	var clientConfig goopenai.ClientConfig

	// check if it's azure or not
	if opts.BaseURL == "" {
		clientConfig = goopenai.DefaultConfig(authToken)
	} else {
		clientConfig = goopenai.DefaultAzureConfig(authToken, opts.BaseURL)
		clientConfig.OrgID = opts.OrgID
	}

	if opts.APIVersion != "" {
		clientConfig.APIVersion = opts.APIVersion
	}

	client = goopenai.NewClientWithConfig(clientConfig)
	return
}

// NewOpenAIModel return new openAI Model instance
func NewOpenAIModel(authToken string, modelName string, callbackManager *callback.Manager, options ...func(*OpenAIOption)) (llm *OpenAIModel) {
	opts := OpenAIOption{}
	for _, opt := range options {
		opt(&opts)
	}

	client := newOpenAIClient(authToken, opts)

	llm = &OpenAIModel{
		c:               client,
		modelName:       modelName,
		callbackManager: callbackManager,
	}

	if opts.Verbose {
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
