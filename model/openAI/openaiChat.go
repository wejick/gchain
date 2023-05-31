package _openai

import (
	"context"

	goopenai "github.com/sashabaranov/go-openai"
	model "github.com/wejick/gochain/model"
)

var _ model.LLMModel = &OpenAIChatModel{}

type OpenAIChatModel struct {
	c         *goopenai.Client
	modelName string
}

// NewOpenAIChatModel return new openAI Model instance
func NewOpenAIChatModel(authToken string, orgID string, modelName string) (model *OpenAIChatModel) {
	var client *goopenai.Client
	if orgID != "" {
		client = goopenai.NewClient(authToken)
	} else {
		config := goopenai.DefaultConfig(authToken)
		config.OrgID = orgID
		client = goopenai.NewClientWithConfig(config)
	}

	model = &OpenAIChatModel{
		c:         client,
		modelName: modelName,
	}

	return
}

// Call runs completion on chat model, the prompt will be put as user chat
func (O *OpenAIChatModel) Call(ctx context.Context, prompt string, options ...func(*model.Option)) (output string, err error) {
	messages := []model.ChatMessage{
		{Role: model.ChatMessageRoleUser, Content: prompt},
	}
	responds, err := O.Chat(ctx, messages, options...)
	if err != nil {
		return
	} else {
		output = responds.Content
	}

	return
}

// Chat call chat completion
func (O *OpenAIChatModel) Chat(ctx context.Context, messages []model.ChatMessage, options ...func(*model.Option)) (output model.ChatMessage, err error) {
	opts := model.Option{}
	for _, opt := range options {
		opt(&opts)
	}

	request := goopenai.ChatCompletionRequest{
		Model:       goopenai.GPT3Dot5Turbo,
		MaxTokens:   opts.MaxToken,
		Temperature: opts.Temperature,
		Messages:    convertMessagesToOai(messages),
		Stream:      false,
	}

	response, err := O.c.CreateChatCompletion(ctx, request)
	if err != nil {
		return
	} else if len(response.Choices) > 0 {
		output = model.ChatMessage{
			Role:    response.Choices[0].Message.Role,
			Content: response.Choices[0].Message.Content,
		}
	}

	return
}

// Chat call chat completion in streaming mode
func (O *OpenAIChatModel) ChatStreaming(ctx context.Context, prompt string, options ...func(*model.Option)) (output string, err error) {
	opts := model.Option{}
	for _, opt := range options {
		opt(&opts)
	}

	return
}

func convertMessageToOai(chatMessage model.ChatMessage) goopenai.ChatCompletionMessage {
	return goopenai.ChatCompletionMessage{
		Role:    chatMessage.Role,
		Content: chatMessage.Content,
	}
}

func convertMessagesToOai(chatMessages []model.ChatMessage) []goopenai.ChatCompletionMessage {
	var completionMessages []goopenai.ChatCompletionMessage
	for _, message := range chatMessages {
		completionMessages = append(completionMessages, convertMessageToOai(message))
	}
	return completionMessages
}
