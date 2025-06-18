package godefs

import (
	"github.com/ab36245/go-strcase"

	"github.com/ab36245/go-modelgen/defx"
)

type Field struct {
	Name   string
	Source string
	Type   *Type
}

func newField(d defx.Field) Field {
	return Field{
		Name:   strcase.ToPascal(d.Name),
		Source: d.Name,
		Type:   newType(d.Type, 0),
	}
}
