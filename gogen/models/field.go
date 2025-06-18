package models

import (
	"github.com/ab36245/go-modelgen/writer"

	"github.com/ab36245/go-modelgen/gogen/godefs"
)

func doFieldStruct(w writer.GenWriter, f godefs.Field) {
	w.Put("%s %s", f.Name, f.Type.Name)
}

func doFieldDecode(w writer.GenWriter, f godefs.Field) {
	w.Put("// %s", f.Name)
	w.Inc("{")
	{
		target := doTypeDecode(w, f.Type, f.Source, "v")
		w.Put("m.%s = %s", f.Name, target)
	}
	w.Dec("}")
}

func doFieldEncode(w writer.GenWriter, f godefs.Field) {
	w.Put("// %s", f.Name)
	w.Inc("{")
	{
		w.Put("// TODO")
	}
	w.Dec("}")
}
