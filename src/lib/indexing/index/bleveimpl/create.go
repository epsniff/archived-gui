package bleveimpl

import (
	"fmt"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
)

type Cfg struct {
	Path     string
	Name     string
	ReadOnly bool

	Mappings mapping.IndexMapping //most likely of type *mapping.IndexMappingImpl
}

func New(cfg *Cfg) (*BleveIndex, error) {

	if cfg.Mappings == nil {
		cfg.Mappings = bleve.NewIndexMapping()
	}

	config := map[string]interface{}{
		"read_only": cfg.ReadOnly,
		"path":      cfg.Path,
	}
	newIndex, err := scorch.NewScorch(cfg.Name, config, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating scorch index: %v", err)
	}

	return &BleveIndex{
		mu:       &sync.RWMutex{},
		index:    newIndex,
		mappings: cfg.Mappings,
	}, nil

	//newIndex, err := bleve.New(cfg.Path, cfg.Mappings)
	//if err != nil {
	//	return nil, fmt.Errorf("error creating index: %v", err)
	//}
	//newIndex.SetName(cfg.Name)
	//return &BleveIndex{&sync.RWMutex{}, newIndex}, nil
}
