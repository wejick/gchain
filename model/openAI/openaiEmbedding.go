package _openai

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"
	goopenai "github.com/sashabaranov/go-openai"
)

type OpenAIEmbedModel struct {
	c     *goopenai.Client
	model goopenai.EmbeddingModel
}

// NewOpenAIEmbedModel return new openAI Model instance
func NewOpenAIEmbedModel(authToken string, orgID string, modelName openai.EmbeddingModel) (model *OpenAIEmbedModel) {
	var client *goopenai.Client
	if orgID != "" {
		client = goopenai.NewClient(authToken)
	} else {
		config := goopenai.DefaultConfig(authToken)
		config.OrgID = orgID
		client = goopenai.NewClientWithConfig(config)
	}

	model = &OpenAIEmbedModel{
		c:     client,
		model: modelName,
	}

	return
}

// EmbedQuery produce embedding for a string of query
func (m *OpenAIEmbedModel) EmbedQuery(input string) (embedding []float32, err error) {
	embeddings, err := m.EmbedDocuments([]string{input})
	if err != nil || len(embeddings) == 0 {
		return
	}

	embedding = embeddings[0]

	return
}

// EmbedQuery produce embedding for a list of documents
func (m *OpenAIEmbedModel) EmbedDocuments(documents []string) (embeddings [][]float32, err error) {
	resp, err := m.c.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: documents,
			Model: m.model,
		})

	if err != nil || len(resp.Data) == 0 {
		return nil, errors.New("CreateEmbeddings failed" + err.Error())
	}

	for _, data := range resp.Data {
		embeddings = append(embeddings, data.Embedding)
	}

	return
}
