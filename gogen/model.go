package gogen

import (
	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
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

func (m Model) doCodec(w writer.GenWriter) {
	w.Inc("var %sCodec = model.Codec[%s]{", m.Name, m.Name)
	{
		m.doDecode(w)
		m.doEncode(w)
	}
	w.Dec("}")
}

func (m Model) doDecode(w writer.GenWriter) {
	w.Inc("Decode: func(d model.ObjectDecoder) (%s, error) {", m.Name)
	{
		w.Put("m := %s{}", m.Name)
		for _, f := range m.Fields {
			f.doDecode(w)
		}
		w.Put("return m, nil")
	}
	w.Dec("},")
}

func (m Model) doEncode(w writer.GenWriter) {
	w.Inc("Encode: func(e model.ObjectEncoder, m %s) error {", m.Name)
	{
		w.Put("var err error")
		for _, f := range m.Fields {
			f.doEncode(w)
		}
		w.Put("return nil")
	}
	w.Dec("},")
}

func (m Model) doString(w writer.GenWriter) {
	w.Inc("func (m %s) String() string {", m.Name)
	{
		w.Put("return model.String(m)")
	}
	w.Dec("}")
}

func (m Model) doStruct(w writer.GenWriter) {
	w.Inc("type %s struct {", m.Name)
	{
		for _, f := range m.Fields {
			f.doStruct(w)
		}
	}
	w.Dec("}")
}
