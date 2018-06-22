package bleveimpl

import (
	"fmt"

	"github.com/blevesearch/bleve/document"
	"github.com/epsniff/gui/src/lib/indexing/shared"
	"github.com/epsniff/gui/src/lib/logging"
)

/*
docIfaceToBleveDoc converts our document type to one that bleve can index natively.

NOTES: instead of using `b.mappings.MapDocument(bdoc, data)` it would be possible to manage the mappings ourselfs using something like the snippet below.  Which would allow for custom analizers or to do things like round datetime...document

	switch val.Type() {
	case shared.StringType:
		analyzer := fm.analyzerForField(path, context)
		field := document.NewTextFieldCustom(fieldName, indexes, []byte(propertyValueString), options, analyzer)
		context.doc.AddField(field)

		if !fm.IncludeInAll {
			context.excludedFromAll = append(context.excludedFromAll, fieldName)
		}
	case shared.TimeType:
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
*/
func (b *BleveIndex) docIfaceToBleveDoc(id string, doc shared.DocumentIface) (*document.Document, error) {
	data := doc.AsMap()

	bdoc := document.NewDocument(id)
	if err := b.mappings.MapDocument(bdoc, data); err != nil {
		return nil, err
	}
	bdoc.ID = id
	return bdoc, nil
}

func (b *BleveIndex) bleveDocToDoc(doc *document.Document) (*shared.Document, error) {
	gdoc := &shared.Document{Id: doc.ID, Fields: map[string][]*shared.Value{}}

	for _, field := range doc.Fields {
		val := &shared.Value{}
		val.Name = field.Name()

		switch field := field.(type) {
		case *document.TextField:
			val.Type = shared.StringType
			val.Val = string(field.Value())
		case *document.NumericField:
			val.Type = shared.NumberType
			n, err := field.Number()
			if err != nil {
				return nil, fmt.Errorf("bleve error: field `%v`'s Number returned: err:%v", val.Name, err)
			}
			val.Val = n
		case *document.DateTimeField:
			val.Type = shared.TimeType
			d, err := field.DateTime()
			if err != nil {
				return nil, fmt.Errorf("bleve error: field `%v`'s DateTime returned: err:%v", val.Name, err)
			}
			val.Val = d
		case *document.BooleanField:
			val.Type = shared.BoolType
			b, err := field.Boolean()
			if err != nil {
				return nil, fmt.Errorf("bleve error: field `%v`'s Boolean returned: err:%v", val.Name, err)
			}
			val.Val = b
		case *document.GeoPointField:
			val.Type = shared.GeoPointType
			lon, err := field.Lon()
			if err == nil {
				lat, err := field.Lat()
				if err == nil {
					val.Val = []float64{lon, lat}
				}
			}
		//TODO add case *document.CompositeField:
		default:
			logging.Logger.Warnf("unsupported type? field:%v: type:%T", field.Name(), field)
		}

		vals, existed := gdoc.Fields[field.Name()]
		if existed {
			gdoc.Fields[field.Name()] = append(vals, val) //existing value, so append
		} else {
			gdoc.Fields[field.Name()] = []*shared.Value{val} //new value
		}
	}
	return gdoc, nil
}
