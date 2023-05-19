package _openai

import (
	"context"
	"os"
	"testing"

	model "github.com/wejick/gochain/model"
)

var authToken = os.Getenv("OPENAI_API_KEY")

func TestOpenAICall(t *testing.T) {
	var testModel = NewOpenAIModel(authToken, "", "text-ada-001")
	output, err := testModel.Call(context.Background(), "we are us, we are us, we are ", model.WithTemperature(0))
	if err != nil {
		t.Error(err)
	} else {
		t.Log("output : ", output)
	}
}
