package shared

import (
	"encoding/json"
	"time"

	"github.com/blevesearch/bleve/document"
)

func IndexNameValid(name string) bool {
	if name == "" {
		return false
	}
	return true
}

func DocumentIDValid(name string) bool {
	if name == "" {
		return false
	}
	return true
}

/// The following convert documents into a type thats easy to marshal into json.
///
func NewJsonDocument(docID string) *JsonDocument {
	return &JsonDocument{
		ID:     docID,
		Fields: map[string]interface{}{},
	}
}

type JsonDocument struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

func (jd *JsonDocument) Marshal() ([]byte, error) {
	return json.Marshal(jd)
}

func ToJsonableDocument(doc document.Document) (*JsonDocument, error) {
	jsonDoc := NewJsonDocument(doc.ID)

	for _, field := range doc.Fields {
		var newval interface{}
		switch field := field.(type) {
		case *document.TextField:
			newval = string(field.Value())
		case *document.NumericField:
			n, err := field.Number()
			if err == nil {
				newval = n
			}
		case *document.DateTimeField:
			d, err := field.DateTime()
			if err == nil {
				newval = d.Format(time.RFC3339Nano)
			}
		}

		val, ok := jsonDoc.Fields[field.Name()]
		if ok {
			switch fval := val.(type) {
			case []interface{}:
				jsonDoc.Fields[field.Name()] = append(fval, newval)
			case interface{}:
				arr := make([]interface{}, 2)
				arr[0] = fval
				arr[1] = newval
				jsonDoc.Fields[field.Name()] = arr
			}
		} else {
			jsonDoc.Fields[field.Name()] = newval
		}
	}
	return jsonDoc, nil
}

/*
import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)


func showError(w http.ResponseWriter, r *http.Request,
	msg string, code int) {
	logger.Printf("Reporting error %v/%v", code, msg)
	http.Error(w, msg, code)
}

func mustEncode(w io.Writer, i interface{}) {
	if headered, ok := w.(http.ResponseWriter); ok {
		headered.Header().Set("Cache-Control", "no-cache")
		headered.Header().Set("Content-type", "application/json")
	}

	e := json.NewEncoder(w)
	if err := e.Encode(i); err != nil {
		panic(err)
	}
}

type varLookupFunc func(req *http.Request) string

var logger = log.New(ioutil.Discard, "indexing", log.LstdFlags)

// SetLog sets the logger used for logging
// by default log messages are sent to ioutil.Discard
func SetLog(l *log.Logger) {
	logger = l
}
*/
