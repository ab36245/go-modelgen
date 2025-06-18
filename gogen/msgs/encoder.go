package msgs

import (
	"github.com/ab36245/go-modelgen/writer"

	"github.com/ab36245/go-modelgen/gogen/godefs"
)

func doEncoder(w writer.GenWriter, ds []godefs.Model) {
	w.Inc("func EncodeMsg(v any) ([]byte, error) {")
	{
		w.Put("e, err := msgs.Encoder()")
		w.Inc("if err != nil {")
		{
			w.Put("return nil, err")
		}
		w.Dec("}")
		w.Put("switch v := v.(type) {")
		for _, d := range ds {
			w.Inc("case %s:", d.Name)
			{
				w.Inc("if err := e.Id(%d); err != nil {", d.Id)
				{
					w.Put("return nil, err")
				}
				w.Dec("}")
				w.Inc("if err := %sCodec.Encode(e, v); err != nil {", d.Name)
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
}
