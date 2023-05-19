package model

import "context"

//go:generate moq -out model_moq.go . LLMModel
type LLMModel interface {
	Call(ctx context.Context, prompt string, options ...func(*Option)) (string, error)
}

type Option struct {
	Temperature float32
}

func WithTemperature(temp float32) func(*Option) {
	return func(o *Option) {
		o.Temperature = temp
	}
}
