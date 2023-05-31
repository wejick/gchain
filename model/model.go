package model

import "context"

//go:generate moq -out model_moq.go . LLMModel
type LLMModel interface {
	Call(ctx context.Context, prompt string, options ...func(*Option)) (string, error)
}

type Option struct {
	Temperature float32
	MaxToken    int
}

func WithTemperature(temp float32) func(*Option) {
	return func(o *Option) {
		o.Temperature = temp
	}
}

func MaxToken(maxToken int) func(*Option) {
	return func(o *Option) {
		o.MaxToken = maxToken
	}
}

type ChatModel interface {
	LLMModel
	Chat(ctx context.Context, messages []ChatMessage, options ...func(*Option)) (ChatMessage, error)
	ChatStreaming(ctx context.Context, prompt string, options ...func(*Option)) (string, error)
}

type ChatMessage struct {
	Role    string
	Content string
}

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)
