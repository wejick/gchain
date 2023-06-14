package datastore

import "context"

// TODO :
// function to create index with specified schema
// ability to pass extra option to the vectorstore like minimum distance, certainty, limit, and many other

type DataStore interface {
	// Search using a query string
	Search(ctx context.Context, indexName string, query string) ([]interface{}, error)
	// AddText store text to vector index
	// it will also store embedding of the text using specified embedding model
	// If the index doesnt exist, it will create one with default schema
	AddText(ctx context.Context, indexName string, input string) (err error)
	// AddDocuments store a batch of documents to vector index
	// it will also store embedding of the document using specified embedding model
	// If the index doesnt exist, it will create one with default schema
	AddDocuments(ctx context.Context, indexName string, documents []string) (batchErr []error, err error)
	// DeleteIndex drop the index
	DeleteIndex(ctx context.Context, indexName string) (err error)
}

type VectorStore interface {
	DataStore
	//SearchVector by providing the vector from embedding function
	SearchVector(ctx context.Context, indexName string, vector []float32) ([]interface{}, error)
}

type Retriever interface {
	// Search using a query string
	Search(ctx context.Context, indexName string, query string) ([]interface{}, error)
}

// Option
// vectorstore you use may not support everything
type Option struct {
	Limit      int64   // max result to return
	Similarity float32 // minimum similarity score
}
