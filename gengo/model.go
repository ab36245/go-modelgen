package gengo

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

func (m Model) doDb(w writer.GenWriter) {
	w.Inc("var %sDbCodec = db.Codec[%s]{", m.Name, m.Name)
	{
		w.Inc("Decode: func(d db.M) (%s, error) {", m.Name)
		{
			w.Put("m := %s{}", m.Name)
			for _, f := range m.Fields {
				f.doDbDecode(w)
			}
			w.Put("return m, nil")
		}
		w.Dec("},")
		w.Inc("Encode: func(m %s) (db.M, error) {", m.Name)
		{
			w.Put("d := make(db.M, %d)", len(m.Fields))
			for _, f := range m.Fields {
				f.doDbEncode(w)
			}
			w.Put("return d, nil")
		}
		w.Dec("},")
	}
	w.Dec("}")
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
