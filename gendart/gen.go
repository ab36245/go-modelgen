package gendart

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
	"github.com/ab36245/go-strcase"
)

func Generate(path string, ds []defx.Model, opts Opts) error {
	dir := filepath.Join(path, "models")
	if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	ms := doMap(ds, newModel)
	if err := genModels(dir, ms, opts); err != nil {
		return err
	}
	if err := genMsgs(dir, ms, opts); err != nil {
		return err
	}
	return nil
}

func genModels(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("  ")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("import 'package:flutter_model/flutter_model.dart';")
	for _, d := range ms {
		w.Put("")
		d.doClass(w)
		w.Put("")
		d.doCodec(w)
	}
	return genSave(dir, "models.dart", opts, w.Code())
}

func genMsgs(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("  ")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("import 'package:flutter_msgs/flutter_msgs.dart';")
	w.Put("import 'models.dart';")
	w.Put("")
	w.Inc("dynamic decodeMsg(Uint8List b) {")
	{
		w.Put("final d = MsgDecoder(b);")
		w.Inc("return switch (d.id) {")
		{
			for _, m := range ms {
				c := strcase.ToCamel(m.Name)
				w.Put("%d => %sCodec.decode(d),", m.Id, c)
			}
			w.Put("_ => throw Exception('unknown model id ${d.id}'),")
		}
		w.Dec("};")
	}
	w.Dec("}")
	w.Put("")
	w.Inc("Uint8List encodeMsg(dynamic v) {")
	{
		w.Put("final e = MsgEncoder();")
		w.Inc("switch (v) {")
		{
			for _, m := range ms {
				w.Inc("case %s():", m.Name)
				{
					w.Put("e.id = %d;", m.Id)
					c := strcase.ToCamel(m.Name)
					w.Put("%sCodec.encode(e, v);", c)
				}
				w.Dec("")
			}
			w.Inc("default:")
			{
				w.Put("throw Exception('unknown model ${v.runtimeType}');")
			}
			w.Dec("")
		}
		w.Dec("}")
		w.Put("return e.bytes;")
	}
	w.Dec("}")
	return genSave(dir, "msgs.dart", opts, w.Code())
}

func genSave(dir string, name string, opts Opts, code string) error {
	file := filepath.Join(dir, name)
	fmt.Printf("Creating %s\n", file)
	if err := os.WriteFile(file, []byte(code), fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	return nil
}
