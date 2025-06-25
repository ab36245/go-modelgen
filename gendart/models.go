package gendart

import (
	"github.com/ab36245/go-modelgen/writer"
)

func genModels(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("  ")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	modelImports(w, ms)
	w.Put("")
	w.Put("// For convenience")
	w.Put("export 'package:flutter_model/flutter_model.dart' show ModelRef;")
	for _, m := range ms {
		w.Put("")
		modelClass(w, m)
	}
	return genSave(dir, "models.dart", opts, w.Code())
}

func modelImports(w writer.GenWriter, ms []Model) {
	names := map[string]bool{
		"package:flutter_model/flutter_model.dart": true,
	}
	types := genTypes(ms)
	_ = types
	// if types[defs.TimeType] {
	// 	names["time"] = true
	// }
	for name := range names {
		w.Put("import '%s';", name)
	}
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
