package _openai

import (
	"reflect"
	"testing"

	"github.com/sashabaranov/go-openai"
	goopenai "github.com/sashabaranov/go-openai"
	model "github.com/wejick/gchain/model"
)

func TestConvertMessageToOai(t *testing.T) {
	message := model.ChatMessage{Role: "system", Content: "Welcome to our system"}
	expected := openai.ChatCompletionMessage{Role: "system", Content: "Welcome to our system"}

	result := convertMessageToOai(message)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Conversion was incorrect, got: %v, want: %v.", result, expected)
	}
}

func TestConvertMessagesToOai(t *testing.T) {
	messages := []model.ChatMessage{
		{Role: "system", Content: "Welcome to our system"},
		{Role: "AI", Content: "Heiho, I am AI"},
		{Role: "user", Content: "Hello, I need assistance"},
		{Role: "AI", FunctionName: "assist", Content: "I will assist you"},
	}

	expected := []openai.ChatCompletionMessage{
		{Role: "system", Content: "Welcome to our system"},
		{Role: "AI", Content: "Heiho, I am AI"},
		{Role: "user", Content: "Hello, I need assistance"},
		{Role: "AI", Name: "assist", Content: "I will assist you"},
	}

	result := convertMessagesToOai(messages)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Conversion was incorrect, got: %v, want: %v.", result, expected)
	}
}

func Test_convertOaiMessageToChat(t *testing.T) {
	type args struct {
		chatMessage goopenai.ChatCompletionMessage
	}
	tests := []struct {
		name string
		args args
		want model.ChatMessage
	}{
		{
			name: "empty message",
		},
		{
			name: "no function system",
			args: args{
				chatMessage: goopenai.ChatCompletionMessage{
					Role:    "system",
					Content: "Hello, I need assistance",
				},
			},
			want: model.ChatMessage{
				Role:    model.ChatMessageRoleSystem,
				Content: "Hello, I need assistance",
			},
		},
		{
			name: "no function Ai",
			args: args{
				chatMessage: goopenai.ChatCompletionMessage{
					Role:    "assistant",
					Content: "Hello, I need assistance",
				},
			},
			want: model.ChatMessage{
				Role:    model.ChatMessageRoleAssistant,
				Content: "Hello, I need assistance",
			},
		},
		{
			name: "need function",
			args: args{
				chatMessage: goopenai.ChatCompletionMessage{
					Role: "assistant",
					FunctionCall: &goopenai.FunctionCall{
						Name:      "test",
						Arguments: "argument here",
					},
				},
			},
			want: model.ChatMessage{
				Role:          model.ChatMessageRoleAssistant,
				Content:       "",
				FunctionName:  "test",
				ParameterJson: "argument here",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertOaiMessageToChat(tt.args.chatMessage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertOaiMessageToChat() = %v, want %v", got, tt.want)
			}
		})
	}
}
