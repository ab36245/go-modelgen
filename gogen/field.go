package gogen

import (
	"fmt"

	"github.com/ab36245/go-strcase"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
)

func newField(d defx.Field) Field {
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

func (f Field) doStruct(w writer.GenWriter) {
	w.Put("%s %s", f.Name, f.Type.Name)
}

func (f Field) doDecode(w writer.GenWriter) {
	w.Put("// %s", f.Name)
	w.Inc("{")
	{
		source := f.Orig
		target := "v"
		target = f.Type.doDecode(w, source, target)
		w.Put("m.%s = %s", f.Name, target)
	}
	w.Dec("}")
}

func (f Field) doEncode(w writer.GenWriter) {
	w.Put("// %s", f.Name)
	w.Inc("{")
	{
		source := fmt.Sprintf("m.%s", f.Name)
		target := f.Orig
		f.Type.doEncode(w, source, target)
	}
	w.Dec("}")
}
