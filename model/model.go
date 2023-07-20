package model

import (
	"context"
)

//go:generate moq -out model_moq.go . LLMModel
type LLMModel interface {
	Call(ctx context.Context, prompt string, options ...func(*Option)) (string, error)
}

//go:generate moq -out model_chat_moq.go . ChatModel
type ChatModel interface {
	LLMModel
	Chat(ctx context.Context, messages []ChatMessage, options ...func(*Option)) (ChatMessage, error)
}

//go:generate moq -out model_embedding_moq.go . EmbeddingModel

type EmbeddingModel interface {
	EmbedQuery(input string) ([]float32, error)
	EmbedDocuments(documents []string) ([][]float32, error)
}

type Option struct {
	Temperature              float32
	StreamingChannel         chan ChatMessage // non chat model can also use this
	Functions                []FunctionDefinition
	AdditionalMetadataFields []string
	MaxToken                 int
	IsStreaming              bool
	Verbose                  bool
}

func WithTemperature(temp float32) func(*Option) {
	return func(o *Option) {
		o.Temperature = temp
	}
}

func WithMaxToken(maxToken int) func(*Option) {
	return func(o *Option) {
		o.MaxToken = maxToken
	}
}

func WithStreamingChannel(streamingChannel chan ChatMessage) func(*Option) {
	return func(o *Option) {
		o.StreamingChannel = streamingChannel
	}
}

func WithReturnMetadataFields(fields []string) func(*Option) {
	return func(o *Option) {
		o.AdditionalMetadataFields = fields
	}
}

func WithStreaming(isStreaming bool) func(*Option) {
	return func(o *Option) {
		o.IsStreaming = isStreaming
	}
}

func WithFunctions(function []FunctionDefinition) func(*Option) {
	return func(o *Option) {
		o.Functions = function
	}
}
