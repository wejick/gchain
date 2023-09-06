package elasticsearchVS

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/wejick/gchain/datastore"
	"github.com/wejick/gchain/document"
	"github.com/wejick/gchain/model"
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
	esClient       *elasticsearch.Client
	embeddingModel model.EmbeddingModel
}

type ElasticHit struct {
	Index   string          `json:"_index,omitempty"`
	ID      string          `json:"_id,omitempty"`
	Score   float64         `json:"_score,omitempty"`
	Source_ json.RawMessage `json:"_source,omitempty"`
}

type ElasticResponse struct {
	Took     int  `json:"took,omitempty"`
	TimedOut bool `json:"timed_out,omitempty"`
	Hits     struct {
		MaxScore float64      `json:"max_score,omitempty"`
		Hits     []ElasticHit `json:"hits,omitempty"`
	} `json:"hits,omitempty"`
}

// NewElasticsearchVectorStore return new Elasticsearch instance
func NewElasticsearchVectorStore(host string, embeddingModel model.EmbeddingModel, esOption ...func(*ESOption)) (EVS *ElasticsearchVectorStore, err error) {
	opts := ESOption{}
	for _, opt := range esOption {
		opt(&opts)
	}

	cfg := elasticsearch.Config{
		Username:     opts.Username,
		Password:     opts.Password,
		Addresses:    []string{host},
		CloudID:      opts.CloudID,
		APIKey:       opts.APIKey,
		ServiceToken: opts.ServiceToken,
	}
	esClient, err := elasticsearch.NewClient(cfg)
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

// A function for marshaling structs to JSON string
func jsonStruct(doc interface{}) (string, error) {
	// Marshal the struct to JSON and check for errors
	b, err := json.Marshal(doc)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return "", err
	}
	return string(b), nil
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
	for _, esDoc := range esDocs {
		jsonDoc, err := jsonStruct(esDoc)
		if err != nil {
			log.Printf("%+v\n", err)
			batchErr = append(batchErr, err)
			continue
		}
		response, err := ES.esClient.Index(
			indexName,
			strings.NewReader(jsonDoc),
			ES.esClient.Index.WithContext(ctx),
		)
		if err != nil {
			log.Printf("%+v\n", err)
			batchErr = append(batchErr, err)
			continue
		}
		if response.IsError() {
			log.Printf("%+v\n", err)
			batchErr = append(batchErr, err)
			continue
		}
		defer response.Body.Close()

	}

	return
}

// DeleteIndex drop the index
func (ES *ElasticsearchVectorStore) DeleteIndex(ctx context.Context, indexName string) (err error) {
	_, err = ES.esClient.Indices.Delete(
		[]string{indexName},
		ES.esClient.Indices.Delete.WithContext(ctx),
	)
	return
}

type KNNSearchBody struct {
	KNN    KNNField `json:"knn"`
	Fields []string `json:"fields"`
}

type KNNField struct {
	Field         string    `json:"field"`
	QueryVector   []float32 `json:"query_vector"`
	K             int64     `json:"k"`
	NumCandidates int64     `json:"num_candidates"`
}

// SearchVector by providing the vector from embedding function
func (ES *ElasticsearchVectorStore) SearchVector(ctx context.Context, indexName string, vector []float32, options ...func(*datastore.Option)) (docs []document.Document, err error) {
	opts := datastore.Option{}
	for _, opt := range options {
		opt(&opts)
	}

	knnSearchBody := KNNSearchBody{
		KNN: KNNField{
			Field:         "vector",
			QueryVector:   vector,
			K:             5,
			NumCandidates: 5,
		},
		Fields: []string{"dense_vector"},
	}

	jsonVector, err := jsonStruct(knnSearchBody)
	if err != nil {
		return
	}

	resp, err := ES.esClient.API.KnnSearch(
		[]string{indexName},
		ES.esClient.API.KnnSearch.WithContext(ctx),
		ES.esClient.API.KnnSearch.WithBody(strings.NewReader(jsonVector)),
	)
	defer resp.Body.Close()

	if resp.IsError() {
		err = fmt.Errorf("es8 got error with status : %+v and response: %+v", resp.Status(), resp)
		return []document.Document{}, err
	}
	respBody := ElasticResponse{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)

	docs, err = hitToDocs(respBody, opts.AdditionalFields)

	return
}

func (ES *ElasticsearchVectorStore) createIndexIfNotExist(ctx context.Context, indexName string) (err error) {
	exist, err := ES.isIndexExist(ctx, indexName)
	if err != nil || exist {
		return
	}
	resp, err := ES.esClient.Indices.Create(
		indexName,
		ES.esClient.Indices.Create.WithBody(strings.NewReader(default_mapping)),
		ES.esClient.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		log.Println("createIndexIfNotExist Indices.Create err: ", err)
		return
	}
	if resp.IsError() {
		err = fmt.Errorf("error createIndex with response: %+v", resp)
	}
	return
}

func (ES *ElasticsearchVectorStore) isIndexExist(ctx context.Context, indexName string) (exist bool, err error) {
	resp, err := ES.esClient.API.Indices.Exists(
		[]string{indexName},
		ES.esClient.API.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 400 || resp.StatusCode == 404 {
		return false, nil
	}
	exist = true
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
func hitToDocs(esRespBody ElasticResponse, additionalFields []string) (docs []document.Document, err error) {
	hits := esRespBody.Hits.Hits
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
