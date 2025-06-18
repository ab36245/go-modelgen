package godefs

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defx"
)

type Type struct {
	Kind  defx.TypeKind
	Name  string
	Level int
	Key   *Type
	Sub   *Type
}

func newType(d *defx.Type, level int) *Type {
	t := &Type{
		Kind:  d.Kind,
		Level: level,
	}
	switch d.Kind {
	case defx.ArrayType:
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("[]%s", t.Sub.Name)
	case defx.BoolType:
		t.Name = "bool"
	case defx.BytesType:
		t.Name = "[]byte"
	case defx.FloatType:
		t.Name = "float64"
	case defx.IntType:
		t.Name = "int"
	case defx.MapType:
		t.Key = newType(d.Key, level+1)
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("map[%s]%s", t.Key.Name, t.Sub.Name)
	case defx.ModelType:
		t.Name = d.Name
	case defx.RefType:
		t.Name = "model.Ref"
	case defx.StringType:
		t.Name = "string"
	case defx.TimeType:
		t.Name = "time.Time"
	}
	return t
}
