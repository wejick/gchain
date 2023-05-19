//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/model"
	_openai "github.com/wejick/gochain/model/openAI"
)

var llmModel model.LLMModel

func TestMain(m *testing.M) {
	fmt.Println("Running integration tests...")
	// Perform any setup or initialization here

	var authToken = os.Getenv("OPENAI_API_KEY")
	llmModel = _openai.NewOpenAIModel(authToken, "", "text-ada-001")

	exitCode := m.Run()

	// Perform any cleanup or teardown here

	// Exit with the appropriate exit code
	// (0 for success, non-zero for failure)
	os.Exit(exitCode)
}

func TestLlmChain(t *testing.T) {
	chain := llm_chain.NewLLMChain(llmModel)
	outputMap, err := chain.Run(context.Background(), map[string]string{"input": "Indonesia Capital is Jakarta\nJakarta is the capital of "})
	assert.NoError(t, err, "error Run")
	assert.Contains(t, outputMap["output"], "Indonesia", "unexpected result")
}
