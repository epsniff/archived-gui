package bleveimpl

import (
	"context"

	"github.com/araddon/qlbridge/rel"
	"github.com/blevesearch/bleve"
)

func (b *BleveIndex) Search(ctx context.Context, fql *rel.FilterStatement) (*bleve.SearchResult, error) {
	res & SearchResult{Hits: 0}

	return searchResponse, nil
}
