//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/wejick/gchain/callback"
	"github.com/wejick/gchain/model"
	_openai "github.com/wejick/gchain/model/openAI"
)

var llmModel model.LLMModel

func TestMain(m *testing.M) {
	fmt.Println("Running integration tests...")

	var authToken = os.Getenv("OPENAI_API_KEY")
	llmModel = _openai.NewOpenAIModel(authToken, "text-ada-001", callback.NewManager())

	exitCode := m.Run()

	os.Exit(exitCode)
}

var authToken = os.Getenv("OPENAI_API_KEY")

func TestOpenAICall(t *testing.T) {
	var testModel = _openai.NewOpenAIModel(authToken, "text-ada-001", callback.NewManager())
	output, err := testModel.Call(context.Background(), "we are us, we are us, we are ", model.WithTemperature(0))
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}

func TestOpenAIChat(t *testing.T) {
	var testModel = _openai.NewOpenAIChatModel(authToken, _openai.GPT3Dot5Turbo0301, callback.NewManager())

	testMessages := []model.ChatMessage{
		{Role: model.ChatMessageRoleUser, Content: "Answer in short and directly, Jakarta is capital city of what ?"},
	}
	output, err := testModel.Chat(context.Background(), testMessages)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}

func TestOpenAIChatCall(t *testing.T) {
	var testModel = _openai.NewOpenAIChatModel(authToken, _openai.GPT3Dot5Turbo0301, callback.NewManager())

	output, err := testModel.Call(context.Background(), "Answer in short and directly, Jakarta is capital city of what ?")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("output : ", output)
	}
}

func TestOpenAIEmbedding(t *testing.T) {
	embeddingModel := _openai.NewOpenAIEmbedModel(authToken, openai.AdaEmbeddingV2)

	embedding, err := embeddingModel.EmbedQuery("answer in short and direct")
	if err != nil {
		t.Error(err)
	}
	if embedding == nil {
		t.Error("embedding is nil")
	}
}
