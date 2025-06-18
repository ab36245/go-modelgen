package gogen

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
)

func newType(d *defx.Type, level int) *Type {
	t := &Type{
		Kind:  d.Kind,
		Level: level,
	}
	switch d.Kind {
	case defx.ArrayType:
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("[]%s", t.Sub.Name)
	case defx.BoolType:
		t.Name = "bool"
	case defx.BytesType:
		t.Name = "[]byte"
	case defx.FloatType:
		t.Name = "float64"
	case defx.IntType:
		t.Name = "int"
	case defx.MapType:
		t.Key = newType(d.Key, level+1)
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("map[%s]%s", t.Key.Name, t.Sub.Name)
	case defx.ModelType:
		t.Name = d.Name
	case defx.RefType:
		t.Name = "model.Ref"
	case defx.StringType:
		t.Name = "string"
	case defx.TimeType:
		t.Name = "time.Time"
	}
	return t
}

type Type struct {
	Kind  defx.TypeKind
	Name  string
	Level int
	Key   *Type
	Sub   *Type
}

func (t *Type) doDecode(w writer.GenWriter, source, target string) string {
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
		doDecodeError(w)
		w.Put("%s := make(%s, %s.Length())", target, t.Name, d)
		w.Inc("for %s := range %s.Length() {", i, d)
		{
			v := t.Sub.doDecode(w, "", "v")
			w.Put("%s[%s] = %s", target, i, v)
		}
		w.Dec("}")

	case defx.BoolType:
		w.Put("%s, err := %s.GetBool(%s)", target, decoder, source)
		doDecodeError(w)

	case defx.BytesType:
		w.Put("%s, err := %s.GetBytes(%s)", target, decoder, source)
		doDecodeError(w)

	case defx.FloatType:
		w.Put("%s, err := %s.GetFloat(%s)", target, decoder, source)
		doDecodeError(w)

	case defx.IntType:
		w.Put("%s, err := %s.GetInt(%s)", target, decoder, source)
		doDecodeError(w)

	case defx.MapType:
		d := fmt.Sprintf("d%d", t.Level)
		w.Put("%s, err := %s.GetMap(%s)", d, decoder, source)
		doDecodeError(w)
		w.Put("%s := make(%s, %s.Length())", target, t.Name, d)
		w.Inc("for range %s.Length() {", d)
		{
			k := t.Key.doDecode(w, "", "k")
			v := t.Sub.doDecode(w, "", "v")
			w.Put("%s[%s] = %s", target, k, v)
		}
		w.Dec("}")

	case defx.ModelType:
		d := fmt.Sprintf("d%d", t.Level)
		w.Put("%s, err := %s.GetObject(%s)", d, decoder, source)
		doDecodeError(w)
		w.Put("%s, err := %sCodec.Decode(%s)", target, t.Name, d)
		doDecodeError(w)

	case defx.RefType:
		w.Put("%s, err := %s.GetRef(%s)", target, decoder, source)
		doDecodeError(w)

	case defx.StringType:
		w.Put("%s, err := %s.GetString(%s)", target, decoder, source)
		doDecodeError(w)

	case defx.TimeType:
		w.Put("%s, err := %s.GetTime(%s)", target, decoder, source)
		doDecodeError(w)
	}
	return target
}

func (t *Type) doEncode(w writer.GenWriter, source, target string) {
	if target != "" {
		target = fmt.Sprintf("%q, ", target)
	}
	encoder := "e"
	if t.Level > 0 {
		encoder += fmt.Sprintf("%d", t.Level-1)
	}
	switch t.Kind {
	case defx.ArrayType:
		e := fmt.Sprintf("e%d", t.Level)
		w.Put("%s, err = %s.PutArray(%slen(%s))", e, encoder, target, source)
		doEncodeError(w)
		v := fmt.Sprintf("v%d", t.Level)
		w.Inc("for _, %s := range %s {", v, source)
		{
			t.Sub.doEncode(w, v, "")
		}
		w.Dec("}")

	case defx.BoolType:
		w.Put("err = %s.PutBool(%s%s)", encoder, target, source)
		doEncodeError(w)

	case defx.BytesType:
		w.Put("err = %s.PutBytes(%s%s)", encoder, target, source)
		doEncodeError(w)

	case defx.FloatType:
		w.Put("err = %s.PutFloat(%s%s)", encoder, target, source)
		doEncodeError(w)

	case defx.IntType:
		w.Put("err = %s.PutInt(%s%s)", encoder, target, source)
		doEncodeError(w)

	case defx.MapType:
		e := fmt.Sprintf("e%d", t.Level)
		w.Put("%s, err = %s.PutMap(%slen(%s))", e, encoder, target, source)
		doEncodeError(w)
		k := fmt.Sprintf("k%d", t.Level)
		v := fmt.Sprintf("v%d", t.Level)
		w.Inc("for %s, %s := range %s {", k, v, source)
		{
			t.Key.doEncode(w, k, "")
			t.Sub.doEncode(w, v, "")
		}
		w.Dec("}")

	case defx.ModelType:
		e := fmt.Sprintf("e%d", t.Level)
		w.Put("%s, err = %s.PutObject(%s)", e, encoder, target)
		doEncodeError(w)
		w.Put("err = %sCodec.Encode(%s, %s)", t.Name, e, source)
		doEncodeError(w)

	case defx.RefType:
		w.Put("err = %s.PutRef(%s%s)", encoder, target, source)
		doEncodeError(w)

	case defx.StringType:
		w.Put("err = %s.PutString(%s%s)", encoder, target, source)
		doEncodeError(w)
	}
}

func doDecodeError(w writer.GenWriter) {
	w.Inc("if err != nil {")
	{
		w.Put("return m, nil")
	}
	w.Dec("}")
}

func doEncodeError(w writer.GenWriter) {
	w.Inc("if err != nil {")
	{
		w.Put("return nil")
	}
	w.Dec("}")
}
