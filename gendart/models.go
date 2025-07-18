package gendart

import (
	"github.com/ab36245/go-modelgen/writer"
)

func genModels(opts Opts, ms Models) error {
	w := writer.WithPrefix("  ")
	modelFile(w, ms)
	return genSave(opts, "models.dart", w.Code())
}

func modelFile(w writer.GenWriter, ms Models) {
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	modelImports(w, ms)
	w.Put("")
	w.Put("// For convenience")
	// w.Put("export 'package:flutter_model/flutter_model.dart' show ModelRef;")
	w.Put("export 'package:dart_model/dart_model.dart' show ModelRef;")
	for _, m := range ms.List {
		w.Put("")
		modelClass(w, m)
	}
}

func modelImports(w writer.GenWriter, ms Models) {
	imports := &Imports{}
	imports.add("package:dart_model/dart_model.dart")
	w.Put(imports.String())
}

func modelClass(w writer.GenWriter, m Model) {
	w.Inc("class %s {", m.Name)
	{
		if len(m.Fields) == 0 {
			w.Put("const %s();", m.Name)
			w.Put("")
			w.Put("@override")
			w.Inc("String toString() =>")
			{
				w.Put("'%s';", m.Name)
			}
			w.Dec("")
		} else {
			for _, f := range m.Fields {
				w.Put("final %s %s;", f.Type.Name, f.Name)
			}
			w.Put("")
			w.Inc("const %s({", m.Name)
			{
				for _, f := range m.Fields {
					w.Put("required this.%s,", f.Name)
				}
			}
			w.Dec("});")
			w.Put("")
			w.Put("@override")
			w.Inc("String toString() =>")
			{
				w.Inc("ObjectWriter('%s')", m.Name)
				{
					for _, f := range m.Fields {
						w.Put(".field('%s', %s)", f.Name, f.Name)
					}
					w.Put(".toString();")
				}
				w.Dec("")
			}
			w.Dec("")
		}
	}
	w.Dec("}")
}
