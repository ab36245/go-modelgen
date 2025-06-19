package gendart

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
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
	w.Put("// TODO")
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
