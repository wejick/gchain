package vectorstore

import "context"

// TODO :
// function to create index with specified schema
// ability to pass extra option to the vectorstore like minimum distance, certainty, limit, and many other

type VectorStore interface {
	//SearchVector by providing the vector from embedding function
	SearchVector(ctx context.Context, indexName string, vector []float32) ([]interface{}, error)
	// SearchKeyword using a query string
	SearchKeyword(ctx context.Context, indexName string, query string) ([]interface{}, error)
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
