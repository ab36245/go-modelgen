package gendart

import (
	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
	"github.com/ab36245/go-strcase"
)

func newModel(d defx.Model) Model {
	return Model{
		Fields: doMap(d.Fields, newField),
		Id:     d.Id,
		Name:   d.Name,
	}
}

type Model struct {
	Fields []Field
	Id     int
	Name   string
}

func (m Model) doClass(w writer.GenWriter) {
	w.Inc("class %s {", m.Name)
	{
		if len(m.Fields) == 0 {
			w.Put("const %s();", m.Name)
			w.Put("")
			w.Put("@override")
			w.Inc("String toString() =>")
			{
				w.Put("'%s';", m.Name)
			}
			w.Dec("")
		} else {
			for _, f := range m.Fields {
				f.doClassDeclaration(w)
			}
			w.Put("")
			w.Inc("const %s({", m.Name)
			{
				for _, f := range m.Fields {
					f.doClassConstructor(w)
				}
			}
			w.Dec("});")
			w.Put("")
			w.Put("@override")
			w.Inc("String toString() =>")
			{
				w.Inc("ObjectWriter('%s')", m.Name)
				{
					for _, f := range m.Fields {
						w.Put(".field('%s', %s)", f.Name, f.Name)
					}
					w.Put(".toString();")
				}
				w.Dec("")
			}
			w.Dec("")
		}
	}
	w.Dec("}")
}

func (m Model) doCodec(w writer.GenWriter) {
	name := strcase.ToCamel(m.Name)
	w.Inc("final %sCodec = ModelCodec<%s>(", name, m.Name)
	{
		m.doCodecDecode(w)
		m.doCodecEncode(w)
	}
	w.Dec(");")
}

func (m Model) doCodecDecode(w writer.GenWriter) {
	w.Inc("decode: (d) {")
	{
		for _, f := range m.Fields {
			f.doDecodeAssignment(w)
		}
		w.Inc("return %s(", m.Name)
		{
			for _, f := range m.Fields {
				f.doDecodeConstructor(w)
			}
		}
		w.Dec(");")
	}
	w.Dec("},")
}

func (m Model) doCodecEncode(w writer.GenWriter) {
	w.Inc("encode: (e, m) {")
	{
		for _, f := range m.Fields {
			f.doEncode(w)
		}
	}
	w.Dec("},")
}
