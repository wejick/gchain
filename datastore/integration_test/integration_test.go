//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/wejick/gchain/callback"
	"github.com/wejick/gchain/datastore"
	elasticsearchVS "github.com/wejick/gchain/datastore/elasticsearch_vector"
	weaviateVS "github.com/wejick/gchain/datastore/weaviate_vector"
	wikipedia "github.com/wejick/gchain/datastore/wikipedia_retriever"
	"github.com/wejick/gchain/document"
	"github.com/wejick/gchain/model"
	_openai "github.com/wejick/gchain/model/openAI"
)

var llmModel model.LLMModel
var embeddingModel model.EmbeddingModel

var OAIauthToken = os.Getenv("OPENAI_API_KEY")

var className string
var story string
var stories []document.Document

const (
	wvhost   = "localhost:8080"
	wvscheme = "http"
	wvApiKey = ""
)

func TestMain(m *testing.M) {
	fmt.Println("Running integration tests...")

	llmModel = _openai.NewOpenAIModel(OAIauthToken, "text-ada-001", callback.NewManager())
	embeddingModel = _openai.NewOpenAIEmbedModel(OAIauthToken, openai.AdaEmbeddingV2)
	metadata := map[string]interface{}{
		"url":  "https://wejick.wordpress.com",
		"time": 1847,
	}
	className = "Story"
	story = "In the depths of the forest, a lone wolf found an abandoned puppy and raised it as its own. Years later, the once-lonely wolf and the playful dog became an inseparable duo, roaming the wilderness together."
	stories = []document.Document{
		{Text: "As the sun set over the city skyline, a street musician's haunting melody caught the attention of a passerby, transporting them to a world of forgotten dreams and lost love in just a few melancholic notes.", Metadata: metadata},
		{Text: "In a bustling café, a barista noticed a regular customer's worn-out shoes and secretly left a brand new pair beside their table, inspiring a ripple of anonymous acts of kindness that spread throughout the community.", Metadata: metadata},
		{Text: "A bookworm stumbled upon a dusty, forgotten tome in the attic, and with each turn of the page, they were transported to extraordinary worlds, becoming the hero of their own epic adventures.", Metadata: metadata},
		{Text: "As the rain poured relentlessly, a gardener watched in awe as her wilting flowers began to bloom, realizing that sometimes the greatest growth comes from enduring the storms of life.", Metadata: metadata},
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestWeaviate(t *testing.T) {
	wvClient, err := weaviateVS.NewWeaviateVectorStore(wvhost, wvscheme, wvApiKey, embeddingModel, nil)
	assert.NoError(t, err, err)

	err = wvClient.AddText(context.Background(), className, story)
	assert.NoError(t, err, "AddText")

	batchErr, err := wvClient.AddDocuments(context.Background(), className, stories)
	assert.NoError(t, err, "addDocuments")
	for _, e := range batchErr {
		assert.NoError(t, e, "addDocuments batchErr")
	}

	response, err := wvClient.Search(context.Background(), className, "musician's melody")
	assert.NoError(t, err)
	if len(response) > 0 {
		assert.Contains(t, response[0].Text, "skyline")
	} else {
		t.Error("response is empty")
	}

	vectorQuery, err := embeddingModel.EmbedQuery("musician's melody")
	response, err = wvClient.SearchVector(context.Background(), className, vectorQuery, datastore.WithAdditionalFields([]string{"url", "time"}))
	assert.NoError(t, err)
	if len(response) > 0 {
		assert.Contains(t, response[0].Text, "skyline")
		assert.Equal(t, response[0].Metadata["url"], "https://wejick.wordpress.com")
		assert.Equal(t, response[0].Metadata["time"], float64(1847))
	} else {
		t.Error("response is empty")
	}

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

func TestWikipedia(t *testing.T) {
	w := &wikipedia.Wikipedia{}
	result, err := w.Search(context.Background(), "", "indonesia")
	assert.NoError(t, err)
	for _, res := range result {
		assert.Contains(t, res.Text, "Indonesia")
	}
}

func TestElastic(t *testing.T) {
	if os.Getenv("INTEGRATION_SKIP_ES") == "true" {
		t.Skip("Skipping TestElastic")
	}
	esClient, err := elasticsearchVS.NewElasticsearchVectorStore("http://localhost:9200", embeddingModel)
	assert.NoError(t, err)

	batchErr, err := esClient.AddDocuments(context.Background(), strings.ToLower(className), stories)
	assert.NoError(t, err, "addDocuments")
	for _, e := range batchErr {
		assert.NoError(t, e, "addDocuments batchErr")
	}

	time.Sleep(1 * time.Second) // to give time es to ingest the documents

	response, err := esClient.Search(context.Background(), strings.ToLower(className), "city skyline")
	assert.NoError(t, err)
	if len(response) > 0 {
		assert.Contains(t, response[0].Text, "skyline")
	} else {
		t.Error("response is empty")
	}

	vectorQuery, err := embeddingModel.EmbedQuery("city skyline")
	response, err = esClient.SearchVector(context.Background(), strings.ToLower(className), vectorQuery, datastore.WithAdditionalFields([]string{"url", "time", "nothing"}))
	assert.NoError(t, err)
	if len(response) > 0 {
		assert.Contains(t, response[0].Text, "skyline")
		assert.Contains(t, response[0].Text, "skyline")
		assert.Equal(t, response[0].Metadata["url"], "https://wejick.wordpress.com")
		assert.Equal(t, response[0].Metadata["time"], float64(1847))
		assert.Equal(t, response[0].Metadata["nothing"], nil)
	} else {
		t.Error("response is empty")
	}

	err = esClient.DeleteIndex(context.Background(), strings.ToLower(className))
	assert.NoError(t, err)
}
