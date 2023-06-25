package elasticsearchVS

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/wejick/gochain/datastore"
	"github.com/wejick/gochain/document"
	"github.com/wejick/gochain/model"
)

var _ datastore.VectorStore = &ElasticsearchVectorStore{}

const default_mapping = `{
	"mappings": {
	  "properties": {
		"vector": {
		  "type": "dense_vector",
		  "dims": 1536,
		  "index": true,
		  "similarity": "cosine"
		},
		"text": {
		  "type": "text"
		}
	  }
	}
  }
`

type elasticDocument map[string]interface{}

type ESOption struct {
	Username string
	Password string

	CloudID      string
	APIKey       string
	ServiceToken string
}

// ElasticsearchVectorStore provide access to elasticsearch
type ElasticsearchVectorStore struct {
	esClient       *elasticsearch.TypedClient
	embeddingModel model.EmbeddingModel
}

// NewElasticsearchVectorStore return new Elasticsearch instance
func NewElasticsearchVectorStore(host string, embeddingModel model.EmbeddingModel, esOption ...func(*ESOption)) (EVS *ElasticsearchVectorStore, err error) {
	opts := ESOption{}
	for _, opt := range esOption {
		opt(&opts)
	}

	cfg := elasticsearch.Config{
		Addresses:    []string{host},
		Username:     opts.Username,
		Password:     opts.Password,
		CloudID:      opts.CloudID,
		APIKey:       opts.APIKey,
		ServiceToken: opts.ServiceToken,
	}
	esClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return
	}
	EVS = &ElasticsearchVectorStore{
		esClient:       esClient,
		embeddingModel: embeddingModel,
	}

	return
}

// Search using a query string
func (ES *ElasticsearchVectorStore) Search(ctx context.Context, indexName string, query string, options ...func(*datastore.Option)) (docs []document.Document, err error) {
	vector, err := ES.embeddingModel.EmbedQuery(query)
	if err != nil {
		return
	}
	docs, err = ES.SearchVector(ctx, indexName, vector, options...)

	return
}

// AddText store text to vector index
// it will also store embedding of the text using specified embedding model
// If the index doesnt exist, it will create one with default schema
func (ES *ElasticsearchVectorStore) AddText(ctx context.Context, indexName string, input string) (err error) {
	_, err = ES.AddDocuments(ctx, indexName, []document.Document{{Text: input}})
	return
}

// AddDocuments store a batch of documents to vector index
// it will also store embedding of the document using specified embedding model
// If the index doesnt exist, it will create one with default schema
func (ES *ElasticsearchVectorStore) AddDocuments(ctx context.Context, indexName string, documents []document.Document) (batchErr []error, err error) {
	err = ES.createIndexIfNotExist(ctx, indexName)
	if err != nil {
		return
	}

	objVectors, err := ES.embeddingModel.EmbedDocuments(document.DocumentsToStrings(documents))
	if err != nil {
		return
	}

	// TODO make it bulk request
	esDocs := dataToESDoc(documents, objVectors)
	for idx := range esDocs {
		_, err := ES.esClient.Index(indexName).Request(esDocs[idx]).Do(ctx)
		if err != nil {
			log.Println(err)
			batchErr = append(batchErr, err)
			continue
		}
	}

	return
}

// DeleteIndex drop the index
func (ES *ElasticsearchVectorStore) DeleteIndex(ctx context.Context, indexName string) (err error) {
	_, err = ES.esClient.API.Indices.Delete(indexName).Do(ctx)
	return
}

// SearchVector by providing the vector from embedding function
func (ES *ElasticsearchVectorStore) SearchVector(ctx context.Context, indexName string, vector []float32, options ...func(*datastore.Option)) (docs []document.Document, err error) {
	opts := datastore.Option{}
	for _, opt := range options {
		opt(&opts)
	}

	R := search.NewRequest()

	R.Knn = []types.KnnQuery{{
		Field:         "vector",
		QueryVector:   vector,
		K:             5,
		NumCandidates: 5,
	}}

	resp, err := ES.esClient.API.Search().Index(indexName).Request(R).Do(ctx)
	if err != nil {
		return
	}

	docs, err = hitToDocs(resp.Hits.Hits, opts.AdditionalFields)

	return
}

func (ES *ElasticsearchVectorStore) createIndexIfNotExist(ctx context.Context, indexName string) (err error) {
	exist, err := ES.isIndexExist(ctx, indexName)
	if err != nil || exist {
		return
	}
	_, err = ES.esClient.Indices.Create(indexName).Raw(strings.NewReader(default_mapping)).Do(ctx)

	return
}

func (ES *ElasticsearchVectorStore) isIndexExist(ctx context.Context, indexName string) (exist bool, err error) {
	exist, err = ES.esClient.API.Indices.Exists(indexName).IsSuccess(ctx)
	return
}

// merge documents and vector into elasticDocument
func dataToESDoc(documents []document.Document, vector [][]float32) (output []elasticDocument) {
	output = make([]elasticDocument, len(documents))

	for idx, doc := range documents {
		output[idx] = elasticDocument{
			"text":   doc.Text,
			"vector": vector[idx],
		}
		// put metadata
		for k, v := range doc.Metadata {
			output[idx][k] = v
		}
	}

	return
}

// hitToDocs convert es hits to document
func hitToDocs(hits []types.Hit, additionalFields []string) (docs []document.Document, err error) {
	for _, hit := range hits {
		var source map[string]interface{}
		err = json.Unmarshal(hit.Source_, &source)
		if err != nil {
			log.Println(err)
			continue
		}
		doc := document.Document{Text: source["text"].(string), Metadata: make(map[string]interface{})}
		for _, field := range additionalFields {
			doc.Metadata[field] = source[field]
		}
		docs = append(docs, doc)
	}
	return
}
