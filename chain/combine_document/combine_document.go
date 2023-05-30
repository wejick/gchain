package combine_document

import (
	"context"

	"github.com/wejick/gochain/model"
)

type CombinedDocument interface {
	// Combine document and do something
	// the implementation of combine can run LLM againts the doc
	Combine(ctx context.Context, docs []string, options ...func(*model.Option)) (output string, err error)
}
