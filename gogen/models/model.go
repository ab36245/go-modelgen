package models

import (
	"github.com/ab36245/go-modelgen/writer"

	"github.com/ab36245/go-modelgen/gogen/godefs"
)

func doModel(w writer.GenWriter, d godefs.Model) error {
	doModelStruct(w, d)
	w.Put("")
	doModelString(w, d)
	w.Put("")
	doModelCodec(w, d)
	return nil
}

func doModelStruct(w writer.GenWriter, d godefs.Model) {
	w.Inc("type %s struct {", d.Name)
	{
		for _, f := range d.Fields {
			doFieldStruct(w, f)
		}
	}
	w.Dec("}")
}

func doModelString(w writer.GenWriter, d godefs.Model) {
	w.Inc("func (m %s) String() string {", d.Name)
	{
		w.Put("return writer.Value(m)")
	}
	w.Dec("}")
}

func doModelCodec(w writer.GenWriter, d godefs.Model) {
	w.Inc("var %sCodec = model.Codec[%s]{", d.Name, d.Name)
	{
		doModelDecode(w, d)
		doModelEncode(w, d)
	}
	w.Dec("}")
}

func doModelDecode(w writer.GenWriter, d godefs.Model) {
	w.Inc("Decode: func(d model.ObjectDecoder) (%s, error) {", d.Name)
	{
		w.Put("m := %s{}", d.Name)
		for _, f := range d.Fields {
			doFieldDecode(w, f)
		}
		w.Put("return m, nil")
	}
	w.Dec("},")
}

func doModelEncode(w writer.GenWriter, d godefs.Model) {
	w.Inc("Encode: func(e model.ObjectEncoder, m %s) error {", d.Name)
	{
		w.Put("var err error")
		for _, f := range d.Fields {
			doFieldEncode(w, f)
		}
		w.Put("return nil")
	}
	w.Dec("},")
}
