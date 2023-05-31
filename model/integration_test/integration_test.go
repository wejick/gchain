//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/wejick/gochain/model"
	_openai "github.com/wejick/gochain/model/openAI"
)

var llmModel model.LLMModel

func TestMain(m *testing.M) {
	fmt.Println("Running integration tests...")

	var authToken = os.Getenv("OPENAI_API_KEY")
	llmModel = _openai.NewOpenAIModel(authToken, "", "text-ada-001")

	exitCode := m.Run()

	os.Exit(exitCode)
}

var authToken = os.Getenv("OPENAI_API_KEY")

func TestOpenAICall(t *testing.T) {
	var testModel = _openai.NewOpenAIModel(authToken, "", "text-ada-001")
	output, err := testModel.Call(context.Background(), "we are us, we are us, we are ", model.WithTemperature(0))
	if err != nil {
		t.Error(err)
	} else {
		t.Log("output : ", output)
	}
}

func TestOpenAIChat(t *testing.T) {
	var testModel = _openai.NewOpenAIChatModel(authToken, "", _openai.GPT3Dot5Turbo0301)

	testMessages := []model.ChatMessage{
		{Role: model.ChatMessageRoleUser, Content: "Answer in short and directly, Jakarta is capital city of what ?"},
	}
	output, err := testModel.Chat(context.Background(), testMessages)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("output : ", output)
	}
}

func TestOpenAIChatCall(t *testing.T) {
	var testModel = _openai.NewOpenAIChatModel(authToken, "", _openai.GPT3Dot5Turbo0301)

	output, err := testModel.Call(context.Background(), "Answer in short and directly, Jakarta is capital city of what ?")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("output : ", output)
	}
}
