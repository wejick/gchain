package datastore

import (
	"context"

	"github.com/wejick/gchain/document"
)

// DataStore is the interface for storing and retrieving data
type DataStore interface {
	// Search using a query string
	Search(ctx context.Context, indexName string, query string, options ...func(*Option)) ([]document.Document, error)
	// AddText store text to vector index
	// it will also store embedding of the text using specified embedding model
	// If the index doesnt exist, it will create one with default schema
	AddText(ctx context.Context, indexName string, input string) (err error)
	// AddDocuments store a batch of documents to vector index
	// it will also store embedding of the document using specified embedding model
	// If the index doesnt exist, it will create one with default schema
	AddDocuments(ctx context.Context, indexName string, documents []document.Document) (batchErr []error, err error)
	// DeleteIndex drop the index
	DeleteIndex(ctx context.Context, indexName string) (err error)
}

// VectorStore is the interface for storing and retrieving vector data
type VectorStore interface {
	DataStore
	//SearchVector by providing the vector from embedding function
	SearchVector(ctx context.Context, indexName string, vector []float32, options ...func(*Option)) ([]document.Document, error)
}

// Retriever is the interface for retrieving data
// can be used to interact with read only data source like get API
// DataStore and VectorStore are interface compatible with retriever
type Retriever interface {
	// Search using a query string
	Search(ctx context.Context, indexName string, query string, options ...func(*Option)) ([]document.Document, error)
}

// Option give way to pass additional parameter to the datastore
type Option struct {
	Limit            int64    // max result to return
	Similarity       float32  // minimum similarity score
	AdditionalFields []string // list of fields to return, Text field is always returned part of the Document
}

// WithLimit set the limit of the result
func WithLimit(limit int64) func(*Option) {
	return func(o *Option) {
		o.Limit = limit
	}
}

// WithAdditionalFields set the additional fields to query
// Some datastore has different way to handle non existent field, some return empty some return error
func WithAdditionalFields(fields []string) func(*Option) {
	return func(o *Option) {
		o.AdditionalFields = fields
	}
}

// WithSimilarity set the minimum similarity score
// Different datastore has different way to handle similarity score
func WithSimilarity(similarity float32) func(*Option) {
	return func(o *Option) {
		o.Similarity = similarity
	}
}
