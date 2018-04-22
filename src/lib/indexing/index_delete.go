package indexing

import (
	"fmt"
	"net/http"
	"os"
)

type DeleteIndexHandler struct {
	basePath        string
	IndexNameLookup varLookupFunc
}

func NewDeleteIndexHandler(basePath string) *DeleteIndexHandler {
	return &DeleteIndexHandler{
		basePath: basePath,
	}
}

func (h *DeleteIndexHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// find the name of the index to delete
	var indexName string
	if h.IndexNameLookup != nil {
		indexName = h.IndexNameLookup(req)
	}
	if indexName == "" {
		showError(w, req, "index name is required", 400)
		return
	}

	indexToDelete := UnregisterIndexByName(indexName)
	if indexToDelete == nil {
		showError(w, req, fmt.Sprintf("no such index '%s'", indexName), 404)
		return
	}

	// close the index
	err := indexToDelete.Close()
	if err != nil {
		showError(w, req, fmt.Sprintf("error closing index: %v", err), 500)
		return
	}

	// now delete it
	err = os.RemoveAll(h.indexPath(indexName))
	if err != nil {
		showError(w, req, fmt.Sprintf("error deleting index: %v", err), 500)
		return
	}

	rv := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	mustEncode(w, rv)
}

func (h *DeleteIndexHandler) indexPath(name string) string {
	return h.basePath + string(os.PathSeparator) + name
}
