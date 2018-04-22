package indexing

import (
	"fmt"
	"sync"

	"github.com/blevesearch/bleve"
)

var indexNameMapping map[string]bleve.Index
var indexNameMappingLock sync.RWMutex

func RegisterIndexName(name string, idx bleve.Index) {
	indexNameMappingLock.Lock()
	defer indexNameMappingLock.Unlock()

	if indexNameMapping == nil {
		indexNameMapping = make(map[string]bleve.Index)
	}
	indexNameMapping[name] = idx
}

func UnregisterIndexByName(name string) bleve.Index {
	indexNameMappingLock.Lock()
	defer indexNameMappingLock.Unlock()

	if indexNameMapping == nil {
		return nil
	}
	rv := indexNameMapping[name]
	if rv != nil {
		delete(indexNameMapping, name)
	}
	return rv
}

func IndexNames() []string {
	indexNameMappingLock.RLock()
	defer indexNameMappingLock.RUnlock()

	rv := make([]string, len(indexNameMapping))
	count := 0
	for k := range indexNameMapping {
		rv[count] = k
		count++
	}
	return rv
}

func UpdateAlias(alias string, add, remove []string) error {
	indexNameMappingLock.Lock()
	defer indexNameMappingLock.Unlock()

	index, exists := indexNameMapping[alias]
	if !exists {
		// new alias
		if len(remove) > 0 {
			return fmt.Errorf("cannot remove indexes from a new alias")
		}
		indexes := make([]bleve.Index, len(add))
		for i, addIndexName := range add {
			addIndex, indexExists := indexNameMapping[addIndexName]
			if !indexExists {
				return fmt.Errorf("index named '%s' does not exist", addIndexName)
			}
			indexes[i] = addIndex
		}
		indexAlias := bleve.NewIndexAlias(indexes...)
		indexNameMapping[alias] = indexAlias
	} else {
		// something with this name already exists
		indexAlias, isAlias := index.(bleve.IndexAlias)
		if !isAlias {
			return fmt.Errorf("'%s' is not an alias", alias)
		}
		// build list of add indexes
		addIndexes := make([]bleve.Index, len(add))
		for i, addIndexName := range add {
			addIndex, indexExists := indexNameMapping[addIndexName]
			if !indexExists {
				return fmt.Errorf("index named '%s' does not exist", addIndexName)
			}
			addIndexes[i] = addIndex
		}
		// build list of remove indexes
		removeIndexes := make([]bleve.Index, len(remove))
		for i, removeIndexName := range remove {
			removeIndex, indexExists := indexNameMapping[removeIndexName]
			if !indexExists {
				return fmt.Errorf("index named '%s' does not exist", removeIndexName)
			}
			removeIndexes[i] = removeIndex
		}
		indexAlias.Swap(addIndexes, removeIndexes)
	}
	return nil
}
