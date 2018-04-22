package indexing

import (
	"fmt"
	"os"
	"sync"

	"github.com/blevesearch/bleve"
)

type IndexProvider struct {
	mu       *sync.RWMutex
	basepath string
}

func (ip *IndexProvider) IndexCreate(indexName string) error {
	var index bleve.Index
	if index, err := IndexByName(indexName); err != nil {
		return err
	}
	if index != nil {
		return ErrIndexExists
	}
}

func (ip *IndexProvider) indexPath(name string) string {
	return ip.basePath + string(os.PathSeparator) + name
}

func IndexByName(name string) (bleve.Index, error) {
	if !IndexNameValid(name) {
		return nil, ErrInvalidIndex
	}

	indexNameMappingLock.RLock()
	defer indexNameMappingLock.RUnlock()

	idx := indexNameMapping[name]
	if idx == nil {
		return nil, ErrUnknownIndex
	}
	return idx, nil
}
