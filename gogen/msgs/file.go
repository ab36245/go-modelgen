package msgs

import (
	"github.com/ab36245/go-modelgen/writer"

	"github.com/ab36245/go-modelgen/gogen/godefs"
)

func doFile(w writer.GenWriter, ds []godefs.Model) {
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("")
	w.Put("package models")

	w.Put("")
	w.Inc("import (")
	{
		w.Put("\"github.com/ab36245/go-msgs\"")
	}
	w.Dec(")")

	w.Put("")
	doDecoder(w, ds)

	w.Put("")
	doEncoder(w, ds)
}
