package gengo

import (
	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/writer"
)

func genModels(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("\t")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("package models")
	w.Put("")
	modelImports(w, ms)
	w.Put("")
	w.Put("// For convenience")
	w.Put("type Ref = model.Ref")
	for _, m := range ms {
		w.Put("")
		modelStruct(w, m)
		w.Put("")
		modelMethods(w, m)
	}
	return genSave(dir, "models.go", opts, w.Code())
}

func modelImports(w writer.GenWriter, ms []Model) {
	names := map[string]bool{
		"github.com/ab36245/go-model": true,
	}
	types := genTypes(ms)
	if types[defs.TimeType] {
		names["time"] = true
	}
	if len(names) > 0 {
		w.Inc("import (")
		{
			for name := range names {
				w.Put("%q", name)
			}
		}
		w.Dec(")")
	}
}

func modelStruct(w writer.GenWriter, m Model) {
	w.Inc("type %s struct {", m.Name)
	{
		for _, f := range m.Fields {
			w.Put("%s %s", f.Name, f.Type.Name)
		}
	}
	w.Dec("}")
}

func modelMethods(w writer.GenWriter, m Model) {
	w.Inc("func (m %s) String() string {", m.Name)
	{
		w.Put("return model.String(m)")
	}
	w.Dec("}")
}
