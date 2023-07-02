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
		{Role: "user", Content: "Hello, I need assistance"},
	}

	expected := []openai.ChatCompletionMessage{
		{Role: "system", Content: "Welcome to our system"},
		{Role: "user", Content: "Hello, I need assistance"},
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertOaiMessageToChat(tt.args.chatMessage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertOaiMessageToChat() = %v, want %v", got, tt.want)
			}
		})
	}
}
