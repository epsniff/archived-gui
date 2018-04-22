package shared

import (
	"github.com/araddon/qlbridge/value"
)

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
)

type Value interface {
	Val() interface{}
	Type() Type
}

type Document interface {
	Row() map[FieldName]value.Value
}
