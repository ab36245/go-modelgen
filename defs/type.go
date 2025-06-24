package defs

import (
	"fmt"
)

type Type struct {
	Kind TypeKind
	Name string
	Key  *Type
	Sub  *Type
}

func (t *Type) String() string {
	switch t.Kind {
	case ArrayType:
		return fmt.Sprintf("[]%s", t.Sub)
	case BoolType:
		return "bool"
	case BytesType:
		return "bytes"
	case FloatType:
		return "float"
	case IntType:
		return "int"
	case MapType:
		return fmt.Sprintf("[%s]%s", t.Key, t.Sub)
	case ModelType:
		return fmt.Sprintf("model %s", t.Name)
	case StringType:
		return "string"
	case TimeType:
		return "time"
	default:
		return fmt.Sprintf("Unknown (%d)", t.Kind)
	}
}

type TypeKind int

const (
	ArrayType TypeKind = iota
	BoolType
	BytesType
	FloatType
	IntType
	MapType
	ModelType
	RefType
	StringType
	TimeType
)
