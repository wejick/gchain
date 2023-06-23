package wikipedia

import (
	"context"
	"errors"

	gowiki "github.com/trietmn/go-wiki"
	"github.com/wejick/gochain/datastore"
	"github.com/wejick/gochain/document"
)

var _ datastore.Retriever = &Wikipedia{}

type Wikipedia struct {
}

// Search wikipedia article and return the first article's content
func (W *Wikipedia) Search(ctx context.Context, indexName string, query string, options ...func(*datastore.Option)) (output []document.Document, err error) {
	titles, _, err := gowiki.Search(query, 1, false)
	if err != nil {
		return
	}

	if len(titles) == 0 {
		return nil, errors.New("Wikipedia couldn't find any article related to" + query)
	}

	page, err := gowiki.GetPage(titles[0], -1, false, true)
	if err != nil {
		return
	}

	content, err := page.GetContent()
	if err != nil {
		return
	}

	output = append(output, document.Document{Text: content})

	return
}
