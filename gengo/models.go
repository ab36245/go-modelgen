package gengo

import (
	"github.com/ab36245/go-modelgen/writer"
)

func genModels(opts Opts, ms Models) error {
	w := writer.WithPrefix("\t")
	modelFile(w, ms)
	return genSave(opts, "models.go", w.Code())
}

func modelFile(w writer.GenWriter, ms Models) {
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("package models")
	w.Put("")
	modelImports(w, ms)
	w.Put("")
	w.Put("// For convenience")
	w.Put("type Ref = model.Ref")
	for _, m := range ms.List {
		w.Put("")
		modelStruct(w, m)
		w.Put("")
		modelMethods(w, m)
	}
}

func modelImports(w writer.GenWriter, ms Models) {
	imports := &Imports{}
	imports.add("github.com/ab36245/go-model")
	if ms.Types.HasTime() {
		imports.add("time")
	}
	w.Put(imports.String())
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
