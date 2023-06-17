package model

import (
	"context"
)

//go:generate moq -out model_moq.go . LLMModel
type LLMModel interface {
	Call(ctx context.Context, prompt string, options ...func(*Option)) (string, error)
}

type Option struct {
	Temperature      float32
	StreamingChannel chan ChatMessage // non chat model can also use this
	MaxToken         int
	IsStreaming      bool
	Verbose          bool
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

func WithIsStreaming(isStreaming bool) func(*Option) {
	return func(o *Option) {
		o.IsStreaming = isStreaming
	}
}

//go:generate moq -out model_chat_moq.go . ChatModel
type ChatModel interface {
	LLMModel
	Chat(ctx context.Context, messages []ChatMessage, options ...func(*Option)) (ChatMessage, error)
}

type ChatMessage struct {
	Role    string
	Content string
}

func (C *ChatMessage) String() string {
	return C.Role + ": " + C.Content
}

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

type EmbeddingModel interface {
	EmbedQuery(input string) ([]float32, error)
	EmbedDocuments(documents []string) ([][]float32, error)
}
