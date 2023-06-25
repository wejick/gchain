package weaviateVS

import (
	"context"
	"errors"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/wejick/gochain/datastore"
	"github.com/wejick/gochain/document"
	"github.com/wejick/gochain/model"
)

var _ datastore.VectorStore = &WeaviateVectorStore{}

// WeaviateVectorStore provide access to weaviate vector db
type WeaviateVectorStore struct {
	client         *weaviate.Client
	embeddingModel model.EmbeddingModel

	existClass map[string]bool
}

// NewWeaviateVectorStore return new Weaviate Vector Store instance
// headers is optional, if you want to add additional headers to the request
func NewWeaviateVectorStore(host string, scheme string, apiKey string, embeddingModel model.EmbeddingModel, headers map[string]string) (WVS *WeaviateVectorStore, err error) {
	WVS = &WeaviateVectorStore{
		existClass:     map[string]bool{},
		embeddingModel: embeddingModel,
	}
	cfg := weaviate.Config{
		Host:       host,
		Scheme:     scheme,
		Headers:    headers,
		AuthConfig: auth.ApiKey{Value: apiKey},
	}
	WVS.client, err = weaviate.NewClient(cfg)

	return
}

// SearchVector query weaviate using vector
// for weaviate support to return additional field / metadata is not yet implemented,
func (W *WeaviateVectorStore) SearchVector(ctx context.Context, className string, vector []float32, options ...func(*datastore.Option)) (output []document.Document, err error) {
	opts := datastore.Option{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.Similarity == 0 {
		opts.Similarity = 0.8
	}

	query := W.client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithCertainty(opts.Similarity)
	fields := []graphql.Field{
		{Name: "text"},
	}
	// add additional fields
	for _, fieldName := range opts.AdditionalFields {
		fields = append(fields, graphql.Field{
			Name: fieldName,
		})
	}
	resp, err := W.client.GraphQL().Get().WithClassName(className).WithNearVector(query).WithFields(fields...).WithLimit(5).Do(ctx)
	if err != nil {
		return
	}

	output, err = objectsToDocument(className, resp.Data["Get"], opts.AdditionalFields)

	return
}

// Search query weaviate db, the query parameter will be translated into embedding
// the underlying query is the same with SearchVector
func (W *WeaviateVectorStore) Search(ctx context.Context, className string, query string, options ...func(*datastore.Option)) (output []document.Document, err error) {
	vectorQuery, err := W.embeddingModel.EmbedQuery(query)
	if err != nil {
		return
	}

	output, err = W.SearchVector(ctx, className, vectorQuery)

	return
}

// AddText add single string document
func (W *WeaviateVectorStore) AddText(ctx context.Context, className string, input string) (err error) {
	_, err = W.AddDocuments(ctx, className, []document.Document{{Text: input}})
	return
}

// AddDocuments add multiple string documents
func (W *WeaviateVectorStore) AddDocuments(ctx context.Context, className string, documents []document.Document) (batchErr []error, err error) {
	err = W.createClassIfNotExist(ctx, className)
	if err != nil {
		return
	}

	objVectors, err := W.embeddingModel.EmbedDocuments(document.DocumentsToStrings(documents))
	if err != nil {
		return
	}
	objs := documentsToObject(className, documents, objVectors)
	batchResp, err := W.client.Batch().ObjectsBatcher().WithObjects(objs...).Do(ctx)
	if err != nil {
		return
	}
	for _, res := range batchResp {
		if res.Result.Errors != nil {
			batchErr = append(batchErr, errors.New(res.Result.Errors.Error[0].Message))
		}
	}

	return
}

// objectsToDocument convert objects of weaviate query result to gochain document
func objectsToDocument(className string, getObjects models.JSONObject, additionalField []string) (docs []document.Document, err error) {
	/* Response from weaviate
		{
	    "data": {
	        "Get": {
	            "className	": [
	                {
	                    "answer": "DNA",
	                    "category": "SCIENCE",
	                    "question": "In 1953 Watson & Crick built a model of the molecular structure of this, the gene-carrying substance"
	                },
	                {
	                    "answer": "Liver",
	                    "category": "SCIENCE",
	                    "question": "This organ removes excess glucose from the blood & stores it as glycogen"
	                }
	            ]
	        }
	    }
		}
	*/

	Get, ok := getObjects.(map[string]interface{})
	if !ok {
		return
	}

	result, ok := Get[className].([]interface{})
	if !ok {
		return
	}

	for _, data := range result {
		if dataMap, ok := data.(map[string]interface{}); ok {
			doc := document.Document{
				Text:     dataMap["text"].(string),
				Metadata: make(map[string]interface{}),
			}
			for _, field := range additionalField {
				doc.Metadata[field] = dataMap[field]
			}
			docs = append(docs, doc)
		}
	}

	return
}

func documentsToObject(className string, documents []document.Document, vectors [][]float32) (objs []*models.Object) {
	for idx, doc := range documents {
		properties := map[string]any{
			"text": doc.Text,
		}
		// Put metadata to properties
		for key, val := range doc.Metadata {
			properties[key] = val
		}
		objs = append(objs, &models.Object{
			Class:      className,
			Properties: properties,
			Vector:     vectors[idx],
		})
	}
	return
}

func (W *WeaviateVectorStore) createClassIfNotExist(ctx context.Context, className string) (err error) {
	classExist, err := W.isClassExist(ctx, className)
	if !classExist {
		// create classHere
		err = W.createClass(ctx, className)
		if err != nil {
			return
		}
		W.existClass[className] = true
	}

	return
}

// createClass with default schema
func (W *WeaviateVectorStore) createClass(ctx context.Context, className string) (err error) {
	classSchema := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{
				Name:     "text",
				DataType: []string{"text"},
			},
		},
	}
	err = W.client.Schema().ClassCreator().WithClass(classSchema).Do(ctx)

	return
}

// isClassExist check existance of a class
func (W *WeaviateVectorStore) isClassExist(ctx context.Context, className string) (exist bool, err error) {
	if val, ok := W.existClass[className]; ok {
		return val, nil
	}
	exist, err = W.client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil {
		return
	}
	W.existClass[className] = exist

	return
}

// DeleteIndex will delete a class
func (W *WeaviateVectorStore) DeleteIndex(ctx context.Context, className string) (err error) {
	err = W.client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
	return
}
