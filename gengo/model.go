package gengo

import (
	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/writer"
)

func newModel(d defs.Model) Model {
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
		m.doCodecDecode(w)
		m.doCodecEncode(w)
	}
	w.Dec("}")
}

func (m Model) doCodecDecode(w writer.GenWriter) {
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

func (m Model) doCodecEncode(w writer.GenWriter) {
	w.Inc("Encode: func(e model.ObjectEncoder, m %s) error {", m.Name)
	{
		w.Put("var err error")
		w.Put("// Horrible hack to avoid unused variable error")
		w.Put("_ = err")
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
			f.doDeclaration(w)
		}
	}
	w.Dec("}")
}
