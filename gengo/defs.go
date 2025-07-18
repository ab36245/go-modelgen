package gengo

import (
	"fmt"

	"github.com/ab36245/go-strcase"

	"github.com/ab36245/go-modelgen/defs"
)

func newModels(ds []defs.Model) Models {
	return Models{
		List:  newMap(ds, newModel),
		Types: defs.GetTypes(ds),
	}
}

type Models struct {
	List  []Model
	Types defs.Types
}

func newModel(d defs.Model) Model {
	return Model{
		Fields: newMap(d.Fields, newField),
		Id:     d.Id,
		Name:   d.Name,
	}
}

type Model struct {
	Fields []Field
	Id     int
	Name   string
}

func newField(d defs.Field) Field {
	return Field{
		Name: strcase.ToPascal(d.Name),
		Orig: d.Name,
		Type: newType(d.Type, 0),
	}
}

type Field struct {
	Name string
	Orig string
	Type *Type
}

func newType(d *defs.Type, level int) *Type {
	t := &Type{
		Kind:  d.Kind,
		Level: level,
	}
	switch d.Kind {
	case defs.ArrayType:
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("[]%s", t.Sub.Name)
	case defs.BoolType:
		t.Name = "bool"
	case defs.BytesType:
		t.Name = "[]byte"
	case defs.FloatType:
		t.Name = "float64"
	case defs.IntType:
		t.Name = "int"
	case defs.MapType:
		t.Key = newType(d.Key, level+1)
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("map[%s]%s", t.Key.Name, t.Sub.Name)
	case defs.ModelType:
		t.Name = d.Name
	case defs.OptionType:
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("model.Option[%s]", t.Sub.Name)
	case defs.RefType:
		t.Name = "model.Ref"
	case defs.StringType:
		t.Name = "string"
	case defs.TimeType:
		t.Name = "time.Time"
	}
	return t
}

type Type struct {
	Kind  defs.TypeKind
	Name  string
	Level int
	Key   *Type
	Sub   *Type
}

func (t *Type) varName(base string) string {
	return fmt.Sprintf("%s%d", base, t.Level)
}

func newMap[A, B any](as []A, f func(A) B) []B {
	var bs []B
	for _, a := range as {
		bs = append(bs, f(a))
	}
	return bs
}
