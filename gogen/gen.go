package gogen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
)

func Generate(ds []defx.Model, path string, opts Opts) error {
	dir := filepath.Join(path, "models")
	fmt.Printf("Creating %s\n", dir)
	if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	ms := doMap(ds, newModel)
	if err := genModels(dir, ms); err != nil {
		return err
	}
	if err := genMsgs(dir, ms); err != nil {
		return err
	}
	return nil
}

func genModels(dir string, ms []Model) error {
	w := writer.WithPrefix("\t")
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
	for _, d := range ms {
		w.Put("")
		d.doStruct(w)
		w.Put("")
		d.doString(w)
		w.Put("")
		d.doCodec(w)
	}
	return genSave(dir, "models.go", w.Code())
}

func genMsgs(dir string, ms []Model) error {
	w := writer.WithPrefix("\t")
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
			w.Put("return nil, fmt.Errorf(\"unknown model id %d\", d.Id()")
		}
		w.Dec("")
		w.Put("}")
	}
	w.Dec("}")
	w.Put("")
	w.Inc("func EncodeMsg(v any) ([]byte, error) {")
	{
		w.Put("e, err := msgs.Encoder()")
		w.Inc("if err != nil {")
		{
			w.Put("return nil, err")
		}
		w.Dec("}")
		w.Put("switch v := v.(type) {")
		for _, m := range ms {
			w.Inc("case %s:", m.Name)
			{
				w.Inc("if err := e.Id(%d); err != nil {", m.Id)
				{
					w.Put("return nil, err")
				}
				w.Dec("}")
				w.Inc("if err := %sCodec.Encode(e, v); err != nil {", m.Name)
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
	return genSave(dir, "msgs.go", w.Code())
}

func genSave(dir string, name string, code string) error {
	// code, err := format(code)
	// if err != nil {
	// 	return fmt.Errorf("can't reformat code: %w", err)
	// }

	file := filepath.Join(dir, name)
	fmt.Printf("Creating %s\n", file)
	if err := os.WriteFile(file, []byte(code), fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	return nil
}
