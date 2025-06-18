package msgs

import (
	"github.com/ab36245/go-modelgen/writer"

	"github.com/ab36245/go-modelgen/gogen/godefs"
)

func doDecoder(w writer.GenWriter, ds []godefs.Model) {
	w.Inc("func DecodeMsg(b []byte) (any, error) {")
	{
		w.Put("d, err := msgs.Decoder(b)")
		w.Inc("if err != nil {")
		{
			w.Put("return nil, err")
		}
		w.Dec("}")
		w.Put("switch d.Id() {")
		for _, d := range ds {
			w.Inc("case %d:", d.Id)
			{
				w.Put("return %sCodec.Decode(d)", d.Name)
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
}
