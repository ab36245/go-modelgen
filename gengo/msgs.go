package gengo

import (
	"github.com/ab36245/go-modelgen/writer"
)

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
