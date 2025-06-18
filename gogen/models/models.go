package models

import (
	"github.com/ab36245/go-modelgen/writer"

	godefs "github.com/ab36245/go-modelgen/gogen/godefs"
)

func doModels(w writer.GenWriter, ds []godefs.Model) {
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("")
	w.Put("package models")

	w.Put("")
	w.Inc("import (")
	{
		w.Put("\"github.com/ab36245/go-model\"")
		w.Put("\"github.com/ab36245/go-writer\"")
	}
	w.Dec(")")

	for _, d := range ds {
		w.Put("")
		doModel(w, d)
	}
}
