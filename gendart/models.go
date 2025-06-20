package gendart

import (
	"github.com/ab36245/go-modelgen/writer"
)

func genModels(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("  ")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("import 'package:flutter_model/flutter_model.dart';")
	w.Put("")
	w.Put("// For convenience")
	w.Put("export 'package:flutter_model/flutter_model.dart' show ModelRef;")
	for _, d := range ms {
		w.Put("")
		d.doClass(w)
	}
	return genSave(dir, "models.dart", opts, w.Code())
}
