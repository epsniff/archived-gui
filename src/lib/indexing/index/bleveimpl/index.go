package bleveimpl

import (
	"fmt"
	"sync"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/index"
	"github.com/epsniff/gui/src/lib/indexing/shared"
)

//TODO consider making it private?  In either case we should return behind an interface.
type BleveIndex struct {
	mu *sync.RWMutex

	index    index.Index
	analyzer map[shared.FieldName]*analysis.Analyzer
}

func (b *BleveIndex) DocumentCount() (cnt uint64, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// open a reader for this search
	indexReader, err := b.index.Reader()
	if err != nil {
		return 0, fmt.Errorf("error opening index reader %v", err)
	}
	defer func() {
		if cerr := indexReader.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	cnt, err = indexReader.DocCount()
	if err != nil {
		return cnt, err
	}
	return cnt, nil
}

func (b *BleveIndex) DocumentDelete(docID string) (bool, error) {
	// locate the document by id
	if docID == "" {
		return false, shared.ErrInvalidDocId
	}

	//TODO check what error Delete returns if the docID isn't found
	//     Consider if we should return false,nil on a non found error
	if err := b.index.Delete(docID); err != nil {
		return false, err
	}
	return true, nil
}

func (b *BleveIndex) DocumentGet(docID string) (doc shared.Document, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// locate the document by id
	if !shared.DocumentIDValid(docID) {
		return nil, shared.ErrInvalidDocId
	}

	//TODO check what error Delete returns if the docID isn't found
	//     Consider if we should return false,nil on a non found error
	indexReader, err := b.index.Reader()
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := indexReader.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	doc, err = indexReader.Document(docID)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

/*
func (b *BleveIndex) FieldsStats(indexName string) ([]string, error) {

	fields, err := b.index.Fields()
	if err != nil {
		return err
	}

	return fields
}

*/
