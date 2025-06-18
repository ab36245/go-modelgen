package models

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"

	"github.com/ab36245/go-modelgen/gogen/godefs"
)

func doTypeDecode(w writer.GenWriter, t *godefs.Type, source string, target string) string {
	if source != "" {
		source = fmt.Sprintf("%q", source)
	}
	target += fmt.Sprintf("%d", t.Level)
	decoder := "d"
	if t.Level > 0 {
		decoder += fmt.Sprintf("%d", t.Level-1)
	}
	switch t.Kind {
	case defx.ArrayType:
		d := fmt.Sprintf("d%d", t.Level)
		i := fmt.Sprintf("i%d", t.Level)
		w.Put("%s, err := %s.GetArray(%s)", d, decoder, source)
		doTypeDecodeError(w)
		w.Put("%s := make(%s, %s.Length())", target, t.Name, d)
		w.Inc("for %s := range %s.Length() {", i, d)
		{
			v := doTypeDecode(w, t.Sub, "", "v")
			w.Put("%s[%s] = %s", target, i, v)
		}
		w.Dec("}")

	case defx.BoolType:
		w.Put("%s, err := %s.GetBool(%s)", target, decoder, source)
		doTypeDecodeError(w)

	case defx.BytesType:
		w.Put("%s, err := %s.GetBytes(%s)", target, decoder, source)
		doTypeDecodeError(w)

	case defx.FloatType:
		w.Put("%s, err := %s.GetFloat(%s)", target, decoder, source)
		doTypeDecodeError(w)

	case defx.IntType:
		w.Put("%s, err := %s.GetInt(%s)", target, decoder, source)
		doTypeDecodeError(w)

	case defx.RefType:
		w.Put("%s, err := %s.GetRef(%s)", target, decoder, source)
		doTypeDecodeError(w)

	case defx.MapType:
		d := fmt.Sprintf("d%d", t.Level)
		w.Put("%s, err := %s.GetMap(%s)", d, decoder, source)
		doTypeDecodeError(w)
		w.Put("%s := make(%s, %s.Length())", target, t.Name, d)
		w.Inc("for range %s.Length() {", d)
		{
			k := doTypeDecode(w, t.Key, "", "k")
			v := doTypeDecode(w, t.Sub, "", "v")
			w.Put("%s[%s] = %s", target, k, v)
		}
		w.Dec("}")

	case defx.ModelType:
		d := fmt.Sprintf("d%d", t.Level)
		w.Put("%s, err := %s.GetObject(%s)", d, decoder, source)
		doTypeDecodeError(w)
		w.Put("%s, err := %sCodec.Decode(%s)", target, t.Name, d)
		doTypeDecodeError(w)

	case defx.StringType:
		w.Put("%s, err := %s.GetString(%s)", target, decoder, source)
		doTypeDecodeError(w)

	case defx.TimeType:
		w.Put("%s, err := %s.GetTime(%s)", target, decoder, source)
		doTypeDecodeError(w)
	}
	return target
}

func doTypeDecodeError(w writer.GenWriter) {
	w.Inc("if err != nil {")
	{
		w.Put("return m, nil")
	}
	w.Dec("}")
}
