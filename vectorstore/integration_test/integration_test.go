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
	"github.com/wejick/gochain/model"
	_openai "github.com/wejick/gochain/model/openAI"
	weaviateVS "github.com/wejick/gochain/vectorstore/weaviate"
)

var llmModel model.LLMModel
var embeddingModel model.EmbeddingModel

var OAIauthToken = os.Getenv("OPENAI_API_KEY")

const (
	wvhost   = ""
	wvscheme = "https"
	wvApiKey = ""
)

func TestMain(m *testing.M) {
	fmt.Println("Running integration tests...")

	llmModel = _openai.NewOpenAIModel(OAIauthToken, "", "text-ada-001")
	embeddingModel = _openai.NewOpenAIEmbedModel(OAIauthToken, "", openai.AdaEmbeddingV2)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestWeaviate(t *testing.T) {
	wvClient, err := weaviateVS.NewWeaviateVectorStore(wvhost, wvscheme, wvApiKey, embeddingModel, nil)
	assert.NoError(t, err, err)

	className := "Story"
	story := "In the depths of the forest, a lone wolf found an abandoned puppy and raised it as its own. Years later, the once-lonely wolf and the playful dog became an inseparable duo, roaming the wilderness together."
	stories := []string{
		"As the sun set over the city skyline, a street musician's haunting melody caught the attention of a passerby, transporting them to a world of forgotten dreams and lost love in just a few melancholic notes.",
		"In a bustling café, a barista noticed a regular customer's worn-out shoes and secretly left a brand new pair beside their table, inspiring a ripple of anonymous acts of kindness that spread throughout the community.",
		"A bookworm stumbled upon a dusty, forgotten tome in the attic, and with each turn of the page, they were transported to extraordinary worlds, becoming the hero of their own epic adventures.",
		"As the rain poured relentlessly, a gardener watched in awe as her wilting flowers began to bloom, realizing that sometimes the greatest growth comes from enduring the storms of life.",
	}

	err = wvClient.AddText(context.Background(), className, story)
	assert.NoError(t, err, "AddText")

	batchErr, err := wvClient.AddDocuments(context.Background(), className, stories)
	assert.NoError(t, err, "addDocuments")
	for _, e := range batchErr {
		assert.NoError(t, e, "addDocuments batchErr")
	}

	response, err := wvClient.SearchKeyword(context.Background(), className, "city skyline")
	o := convertInterfaceToMap(response[0])
	assert.Contains(t, o["text"], "skyline")

	vectorQuery, err := embeddingModel.EmbedQuery("city skyline")
	response, err = wvClient.SearchVector(context.Background(), className, vectorQuery)
	p := convertInterfaceToMap(response[0])
	assert.Contains(t, p["text"], "skyline")

	err = wvClient.DeleteIndex(context.Background(), className)
	assert.NoError(t, err)

}

func convertInterfaceToMap(input interface{}) map[string]string {
	inputMap := input.(map[string]interface{})
	resultMap := make(map[string]string)
	for key, value := range inputMap {
		resultMap[key] = value.(string)
	}
	return resultMap
}
