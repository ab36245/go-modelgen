package gogen

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
	w := writer.WithPrefix("\t")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("package models")
	w.Put("")
	w.Inc("import (")
	{
		w.Put("\"github.com/ab36245/go-model\"")
	}
	w.Dec(")")
	for _, d := range ms {
		w.Put("")
		d.doStruct(w)
		w.Put("")
		d.doString(w)
		w.Put("")
		d.doCodec(w)
	}
	return genSave(dir, "models.go", opts, w.Code())
}

func genMsgs(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("\t")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("package models")
	w.Put("")
	w.Inc("import (")
	{
		w.Put("\"fmt\"")
		w.Put("")
		w.Put("\"github.com/ab36245/go-msgs\"")
	}
	w.Dec(")")
	w.Put("")
	w.Inc("func DecodeMsg(b []byte) (any, error) {")
	{
		w.Put("d, err := msgs.Decoder(b)")
		w.Inc("if err != nil {")
		{
			w.Put("return nil, err")
		}
		w.Dec("}")
		w.Put("switch d.Id() {")
		for _, m := range ms {
			w.Inc("case %d:", m.Id)
			{
				w.Put("return %sCodec.Decode(d)", m.Name)
			}
			w.Dec("")
		}
		w.Inc("default:")
		{
			w.Put("return nil, fmt.Errorf(\"unknown model id %d\", d.Id())")
		}
		w.Dec("")
		w.Put("}")
	}
	w.Dec("}")
	w.Put("")
	w.Inc("func EncodeMsg(v any) ([]byte, error) {")
	{
		w.Put("var err error")
		w.Put("e := msgs.Encoder()")
		w.Put("switch v := v.(type) {")
		for _, m := range ms {
			w.Inc("case %s:", m.Name)
			{
				w.Put("err = e.Id(%d)", m.Id)
				w.Inc("if err != nil {")
				{
					w.Put("return nil, err")
				}
				w.Dec("}")
				w.Put("err = %sCodec.Encode(e, v)", m.Name)
				w.Inc("if err != nil {")
				{
					w.Put("return nil, err")
				}
				w.Dec("}")
			}
			w.Dec("")
		}
		w.Inc("default:")
		{
			w.Put("return nil, fmt.Errorf(\"unknown model %T\", v)")
		}
		w.Dec("")
		w.Put("}")
		w.Put("return e.Bytes(), nil")
	}
	w.Dec("}")
	return genSave(dir, "msgs.go", opts, w.Code())
}

func genSave(dir string, name string, opts Opts, code string) error {
	if opts.Reformat {
		var err error
		code, err = format(code)
		if err != nil {
			return fmt.Errorf("can't reformat code: %w", err)
		}
	}

	file := filepath.Join(dir, name)
	fmt.Printf("Creating %s\n", file)
	if err := os.WriteFile(file, []byte(code), fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	return nil
}
