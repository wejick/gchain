package _openai

import (
	"context"
	"errors"
	"io"
	"log"

	goopenai "github.com/sashabaranov/go-openai"
	"github.com/wejick/gchain/callback"
	model "github.com/wejick/gchain/model"
)

var _ model.LLMModel = &OpenAIChatModel{}

type OpenAIChatModel struct {
	c               *goopenai.Client
	modelName       string
	callbackManager *callback.Manager
}

// NewOpenAIChatModel return new openAI Model instance
func NewOpenAIChatModel(authToken string, orgID string, baseURL string, apiVersion string, modelName string, callbackManager *callback.Manager, verbose bool) (llm *OpenAIChatModel) {
	var client *goopenai.Client
	if baseURL == "" {
		client = goopenai.NewClient(authToken)
	} else {
		config := goopenai.DefaultAzureConfig(authToken, baseURL)
		config.OrgID = orgID
		if apiVersion != "" {
			config.APIVersion = apiVersion
		}
		client = goopenai.NewClientWithConfig(config)
	}

	llm = &OpenAIChatModel{
		c:               client,
		modelName:       modelName,
		callbackManager: callbackManager,
	}

	if verbose {
		llm.callbackManager.RegisterCallback(model.CallbackModelEnd, callback.VerboseCallback)
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

	// Trigger start callback
	flattenMessages := model.FlattenChatMessages(messages)
	O.callbackManager.TriggerEvent(ctx, model.CallbackModelStart, callback.CallbackData{
		EventName:    model.CallbackModelStart,
		FunctionName: "OpenAIChatModel.Chat",
		Input:        map[string]string{"input": flattenMessages},
		Output:       map[string]string{"output": output.String()},
	})

	// call chatStreaming if it's streaming chat
	if opts.IsStreaming && opts.StreamingChannel != nil {
		output, err = O.chatStreaming(ctx, messages, options...)
		return
	}

	RequestFunctions := []goopenai.FunctionDefinition{}
	for _, f := range opts.Functions {
		RequestFunctions = append(RequestFunctions, goopenai.FunctionDefinition{
			Name:        f.Name,
			Description: f.Description,
			Parameters:  f.Parameters,
		})
	}
	request := goopenai.ChatCompletionRequest{
		Model:       goopenai.GPT3Dot5Turbo,
		MaxTokens:   opts.MaxToken,
		Temperature: opts.Temperature,
		Messages:    convertMessagesToOai(messages),
		Functions:   RequestFunctions,
		Stream:      false,
	}
	if opts.FunctionCall != nil {
		request.FunctionCall = opts.FunctionCall
	}

	response, err := O.c.CreateChatCompletion(ctx, request)
	if err != nil {
		return
	} else if len(response.Choices) > 0 {
		output = convertOaiMessageToChat(response.Choices[0].Message)
	}
	promptUsage := model.PromptUsage{
		PromptTokens:     response.Usage.PromptTokens,
		CompletionTokens: response.Usage.CompletionTokens,
		TotalTokens:      response.Usage.TotalTokens,
	}
	output.PromptUsage = promptUsage
	// Trigger end callback
	O.callbackManager.TriggerEvent(ctx, model.CallbackModelEnd, callback.CallbackData{
		EventName:    model.CallbackModelEnd,
		FunctionName: "OpenAIChatModel.Chat",
		Input:        map[string]string{"input": flattenMessages},
		Output:       map[string]string{"output": output.String()},
	})

	return
}

// chatStreaming call chat completion in streaming mode
// this is a blocking function that will return after all response is completely streamed
// to get the streaming loop streamingCallbackFunc until it's finished
func (O *OpenAIChatModel) chatStreaming(ctx context.Context, messages []model.ChatMessage, options ...func(*model.Option)) (output model.ChatMessage, err error) {
	opts := model.Option{}
	for _, opt := range options {
		opt(&opts)
	}

	request := goopenai.ChatCompletionRequest{
		Model:       goopenai.GPT3Dot5Turbo,
		MaxTokens:   opts.MaxToken,
		Temperature: opts.Temperature,
		Messages:    convertMessagesToOai(messages),
		Stream:      true,
	}

	// reset the channel
	stream, err := O.c.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return
	}
	defer stream.Close()
	for {
		response, errStrea := stream.Recv()
		if errors.Is(errStrea, io.EOF) {
			opts.StreamingChannel <- model.ChatMessage{Role: "signal", Content: "finished"}
			break
		}

		if errStrea != nil {
			log.Printf("Stream error: %v\n", err)
			return
		}

		if len(response.Choices) > 0 {
			content := response.Choices[0].Delta.Content
			opts.StreamingChannel <- model.ChatMessage{Role: goopenai.ChatMessageRoleAssistant, Content: content}

			// update final output
			output.Content += content
			output.Role = goopenai.ChatMessageRoleAssistant
		}
	}

	// Trigger end callback
	O.callbackManager.TriggerEvent(ctx, model.CallbackModelEnd, callback.CallbackData{
		EventName:    model.CallbackModelEnd,
		FunctionName: "OpenAIChatModel.Chat",
		Input:        map[string]string{"input": model.FlattenChatMessages(messages)},
		Output:       map[string]string{"output": output.String()},
	})

	return
}

func convertMessageToOai(chatMessage model.ChatMessage) goopenai.ChatCompletionMessage {
	return goopenai.ChatCompletionMessage{
		Role:    chatMessage.Role,
		Name:    chatMessage.Name,
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

func convertOaiMessageToChat(chatMessage goopenai.ChatCompletionMessage) (message model.ChatMessage) {
	message = model.ChatMessage{
		Role:    chatMessage.Role,
		Content: chatMessage.Content,
	}

	if chatMessage.FunctionCall != nil {
		message.Name = chatMessage.FunctionCall.Name
		message.ParameterJson = chatMessage.FunctionCall.Arguments
	}

	return
}
