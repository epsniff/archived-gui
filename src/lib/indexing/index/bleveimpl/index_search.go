package bleveimpl

import (
	"github.com/araddon/qlbridge/rel"
	"github.com/blevesearch/bleve"
)

func (b *BleveIndex) Search(indexName string, fql *rel.FilterStatement) (*bleve.SearchResult, error) {
	res & SearchResult{Hits: 0}

	return searchResponse, nil
}
