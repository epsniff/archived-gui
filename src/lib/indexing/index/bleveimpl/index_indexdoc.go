package bleveimpl

import (
	"fmt"

	"github.com/araddon/qlbridge/expr"
	"github.com/blevesearch/bleve/document"
	"github.com/epsniff/gui/src/lib/indexing/shared"
)

func (b *BleveIndex) appendField(fieldname string, val shared.Value, doc *document.Document ){
	if fm.Type == "text" {
		analyzer := fm.analyzerForField(path, context)
		field := document.NewTextFieldCustom(fieldName, indexes, []byte(propertyValueString), options, analyzer)
		context.doc.AddField(field)

		if !fm.IncludeInAll {
			context.excludedFromAll = append(context.excludedFromAll, fieldName)
		}
	} else if fm.Type == "datetime" {
		dateTimeFormat := context.im.DefaultDateTimeParser
		if fm.DateFormat != "" {
			dateTimeFormat = fm.DateFormat
		}
		dateTimeParser := context.im.DateTimeParserNamed(dateTimeFormat)
		if dateTimeParser != nil {
			parsedDateTime, err := dateTimeParser.ParseDateTime(propertyValueString)
			if err == nil {
				fm.processTime(parsedDateTime, pathString, path, indexes, context)
			}
		}
	}
}

func (b *BleveIndex) toBleveDoc(id string, doc shared.Document) *document.Document {
	data := doc.Row()

	bdoc := document.NewDocument(id)
	for field, val := range data {

	}

	return doc
}

func (b *BleveIndex) DocumentIndex(docID string, doc expr.EvalContext) error {

	// locate the document by id
	if !shared.DocumentIDValid(docID) {
		return shared.ErrInvalidDocId
	}

	//err = i.m.MapDocument(bdoc, data)
	//if err != nil {
	//	return
	//}

	err := b.index.Update(docID, bdoc)
	if err != nil {
		return err
	}
	return nil
}

func (b *BleveIndex) Fields(indexName string) ([]string, error) {

	fields, err := b.index.Fields()
	if err != nil {
		return nil, fmt.Errorf("bleve.Fields error get fields list: err:%v", err)
	}

	return fields, nil
}
