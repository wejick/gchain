package datastore

import (
	"context"

	"github.com/wejick/gochain/document"
)

// TODO :
// function to create index with specified schema
// ability to pass extra option to the vectorstore like minimum distance, certainty, limit, and many other

// DataStore is the interface for storing and retrieving data
type DataStore interface {
	// Search using a query string
	Search(ctx context.Context, indexName string, query string) ([]document.Document, error)
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
	SearchVector(ctx context.Context, indexName string, vector []float32) ([]document.Document, error)
}

// Retriever is the interface for retrieving data
// can be used to interact with read only data source like get API
// DataStore and VectorStore are interface compatible with retriever
type Retriever interface {
	// Search using a query string
	Search(ctx context.Context, indexName string, query string) ([]document.Document, error)
}

// Option
// vectorstore you use may not support everything
type Option struct {
	Limit      int64   // max result to return
	Similarity float32 // minimum similarity score
}
