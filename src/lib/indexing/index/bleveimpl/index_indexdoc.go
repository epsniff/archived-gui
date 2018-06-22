package bleveimpl

import (
	"fmt"

	"github.com/blevesearch/bleve/index"
	"github.com/epsniff/gui/src/lib/indexing/shared"
)

func (b *BleveIndex) UpdateDocument(id string, doc shared.DocumentIface) error {
	return b.BatchUpdateDocument(map[string]shared.DocumentIface{id: doc})
}

//
// map[string]*document.Document
//    nil value means to delete the document
//    otherwise its a new update
//
func (b *BleveIndex) BatchUpdateDocument(docs map[string]shared.DocumentIface) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ba := index.NewBatch()
	for id, doc := range docs {
		if !shared.DocumentIDValid(id) {
			return shared.ErrInvalidDocId
		}

		if doc == nil {
			ba.Delete(id)
		} else {
			bdoc, err := b.docIfaceToBleveDoc(id, doc)
			if err != nil {
				return err
			}

			ba.Update(bdoc)
		}
	}
	return b.index.Batch(ba)
}

func (b *BleveIndex) Fields(indexName string) ([]string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	indexReader, err := b.index.Reader()
	if err != nil {
		return nil, fmt.Errorf("error opening index reader %v", err)
	}
	defer func() {
		if cerr := indexReader.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	fields, err := indexReader.Fields()
	if err != nil {
		return nil, fmt.Errorf("bleve.Fields error get fields list: err:%v", err)
	}
	return fields, nil
}
