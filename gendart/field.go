package gendart

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/writer"
)

func newField(d defs.Field) Field {
	return Field{
		Name: d.Name,
		Orig: d.Name,
		Type: newType(d.Type, 0),
	}
}

type Field struct {
	Name string
	Orig string
	Type *Type
}

func (f Field) doClassConstructor(w writer.GenWriter) {
	w.Put("required this.%s,", f.Name)
}

func (f Field) doClassDeclaration(w writer.GenWriter) {
	w.Put("final %s %s;", f.Type.Name, f.Name)
}

func (f Field) doDecodeAssignment(w writer.GenWriter) {
	source := f.Orig
	target := f.Name
	target = f.Type.doDecode(w, source, target)
}

func (f Field) doDecodeConstructor(w writer.GenWriter) {
	w.Put("%s: %s,", f.Name, f.Name)
}

func (f Field) doEncode(w writer.GenWriter) {
	source := fmt.Sprintf("m.%s", f.Name)
	target := f.Orig
	f.Type.doEncode(w, source, target)
}
