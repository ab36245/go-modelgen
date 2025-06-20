package gendart

import (
	"github.com/ab36245/go-modelgen/writer"
)

func genCodecs(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("  ")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("import 'package:flutter_model/flutter_model.dart';")
	w.Put("import 'models.dart';")
	for _, d := range ms {
		w.Put("")
		d.doCodec(w)
	}
	return genSave(dir, "codecs.dart", opts, w.Code())
}
