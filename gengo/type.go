package gengo

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/writer"
)

func newType(d *defs.Type, level int) *Type {
	t := &Type{
		Kind:  d.Kind,
		Level: level,
	}
	switch d.Kind {
	case defs.ArrayType:
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("[]%s", t.Sub.Name)
	case defs.BoolType:
		t.Name = "bool"
	case defs.BytesType:
		t.Name = "[]byte"
	case defs.FloatType:
		t.Name = "float64"
	case defs.IntType:
		t.Name = "int"
	case defs.MapType:
		t.Key = newType(d.Key, level+1)
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("map[%s]%s", t.Key.Name, t.Sub.Name)
	case defs.ModelType:
		t.Name = d.Name
	case defs.OptionType:
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("model.Option[%s]", t.Sub.Name)
	case defs.RefType:
		t.Name = "model.Ref"
	case defs.StringType:
		t.Name = "string"
	case defs.TimeType:
		t.Name = "time.Time"
	}
	return t
}

type Type struct {
	Kind  defs.TypeKind
	Name  string
	Level int
	Key   *Type
	Sub   *Type
}

func (t *Type) varName(base string) string {
	return fmt.Sprintf("%s%d", base, t.Level)
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
	case defs.ArrayType:
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

	case defs.BoolType:
		w.Put("%s, err := %s.GetBool(%s)", target, decoder, source)
		doDecodeError(w)

	case defs.BytesType:
		w.Put("%s, err := %s.GetBytes(%s)", target, decoder, source)
		doDecodeError(w)

	case defs.FloatType:
		w.Put("%s, err := %s.GetFloat(%s)", target, decoder, source)
		doDecodeError(w)

	case defs.IntType:
		w.Put("%s, err := %s.GetInt(%s)", target, decoder, source)
		doDecodeError(w)

	case defs.MapType:
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

	case defs.ModelType:
		d := fmt.Sprintf("d%d", t.Level)
		w.Put("%s, err := %s.GetObject(%s)", d, decoder, source)
		doDecodeError(w)
		w.Put("%s, err := %sCodec.Decode(%s)", target, t.Name, d)
		doDecodeError(w)

	case defs.RefType:
		w.Put("%s, err := %s.GetRef(%s)", target, decoder, source)
		doDecodeError(w)

	case defs.StringType:
		w.Put("%s, err := %s.GetString(%s)", target, decoder, source)
		doDecodeError(w)

	case defs.TimeType:
		w.Put("%s, err := %s.GetTime(%s)", target, decoder, source)
		doDecodeError(w)
	}
	return target
}

func (t *Type) doEncode(w writer.GenWriter, source, target string) {
	encoder := "e"
	if t.Level > 0 {
		encoder += fmt.Sprintf("%d", t.Level-1)
	}
	method := func(kind string, source string) string {
		var args string
		if target != "" && source != "" {
			args = fmt.Sprintf("%q, %s", target, source)
		} else if target != "" {
			args = fmt.Sprintf("%q", target)
		} else {
			args = fmt.Sprintf("%s", source)
		}
		return fmt.Sprintf("Put%s(%s)", kind, args)
	}

	switch t.Kind {
	case defs.ArrayType:
		e := fmt.Sprintf("e%d", t.Level)
		m := method("Array", fmt.Sprintf("len(%s)", source))
		w.Put("%s, err := %s.%s", e, encoder, m)
		doEncodeError(w)
		v := fmt.Sprintf("v%d", t.Level)
		w.Inc("for _, %s := range %s {", v, source)
		{
			t.Sub.doEncode(w, v, "")
		}
		w.Dec("}")

	case defs.BoolType:
		m := method("Bool", source)
		w.Put("err = %s.%s", encoder, m)
		doEncodeError(w)

	case defs.BytesType:
		m := method("Bytes", source)
		w.Put("err = %s.%s", encoder, m)
		doEncodeError(w)

	case defs.FloatType:
		m := method("Float", source)
		w.Put("err = %s.%s", encoder, m)
		doEncodeError(w)

	case defs.IntType:
		m := method("Int", source)
		w.Put("err = %s.%s", encoder, m)
		doEncodeError(w)

	case defs.MapType:
		e := fmt.Sprintf("e%d", t.Level)
		m := method("Map", fmt.Sprintf("len(%s)", source))
		w.Put("%s, err := %s.%s", e, encoder, m)
		doEncodeError(w)
		k := fmt.Sprintf("k%d", t.Level)
		v := fmt.Sprintf("v%d", t.Level)
		w.Inc("for %s, %s := range %s {", k, v, source)
		{
			t.Key.doEncode(w, k, "")
			t.Sub.doEncode(w, v, "")
		}
		w.Dec("}")

	case defs.ModelType:
		e := fmt.Sprintf("e%d", t.Level)
		m := method("Object", "")
		w.Put("%s, err := %s.%s", e, encoder, m)
		doEncodeError(w)
		w.Put("err = %sCodec.Encode(%s, %s)", t.Name, e, source)
		doEncodeError(w)

	case defs.RefType:
		m := method("Ref", source)
		w.Put("err = %s.%s", encoder, m)
		doEncodeError(w)

	case defs.StringType:
		m := method("String", source)
		w.Put("err = %s.%s", encoder, m)
		doEncodeError(w)

	case defs.TimeType:
		m := method("Time", source)
		w.Put("err = %s.%s", encoder, m)
		doEncodeError(w)
	}
}

func doDecodeError(w writer.GenWriter) {
	w.Inc("if err != nil {")
	{
		w.Put("return m, err")
	}
	w.Dec("}")
}

func doEncodeError(w writer.GenWriter) {
	w.Inc("if err != nil {")
	{
		w.Put("return err")
	}
	w.Dec("}")
}
