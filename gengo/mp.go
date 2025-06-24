package gengo

import (
	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
)

func genMp(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("\t")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("package models")
	w.Put("")
	w.Inc("import (")
	{
		needTime := false
	loop:
		for _, m := range ms {
			for _, f := range m.Fields {
				if f.Type.Kind == defx.TimeType {
					needTime = true
					break loop
				}
			}
		}
		if needTime {
			w.Put("\"time\"")
		}

		w.Put("\"github.com/ab36245/go-model\"")
	}
	w.Dec(")")
	w.Put("")
	w.Put("// For convenience")
	w.Put("type Ref = model.Ref")
	for _, d := range ms {
		// w.Put("")
		// d.doStruct(w)
		// w.Put("")
		// d.doString(w)
		_ = d
	}
	return genSave(dir, "mp.go", opts, w.Code())
}
