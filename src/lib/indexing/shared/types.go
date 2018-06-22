package shared

type FieldName string

type Type int

const (
	// DO NOT CHANGE
	NilType       Type = 0
	ErrorType     Type = 10
	UnknownType   Type = 20
	NumberType    Type = 100
	IntType       Type = 110
	BoolType      Type = 120
	TimeType      Type = 130
	ByteSliceType Type = 140
	StringType    Type = 200
	StringsType   Type = 210
	GeoPointType  Type = 220
)

//type ValueIface interface {
//	Val() interface{}
//	Type() Type
//}

type DocumentIface interface {
	AsMap() map[string]interface{}
}

type Value struct {
	Type Type
	Name string
	Val  interface{}
}
type Document struct {
	Id     string
	Fields map[string][]*Value
}

func (d *Document) AsMap() map[string]interface{} {
	mapp := map[string]interface{}{}
	for fieldName, vals := range d.Fields {
		if len(vals) == 0 {
			continue
		} else if len(vals) == 1 {
			mapp[fieldName] = vals[0].Val
		} else {
			docVals := []interface{}{}
			for _, v := range vals {
				docVals = append(docVals, v.Val)
			}
		}
	}
	return mapp
}
