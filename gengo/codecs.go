package gengo

import (
	"github.com/ab36245/go-modelgen/writer"
)

func genCodecs(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("\t")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("package models")
	w.Put("")
	w.Inc("import (")
	{
		w.Put("\"github.com/ab36245/go-model\"")
	}
	w.Dec(")")
	w.Put("")
	for _, d := range ms {
		w.Put("")
		d.doCodec(w)
	}
	return genSave(dir, "codecs.go", opts, w.Code())
}
