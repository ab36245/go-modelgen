package gendart

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
	"github.com/ab36245/go-strcase"
)

func newType(d *defx.Type, level int) *Type {
	t := &Type{
		Kind:  d.Kind,
		Level: level,
	}
	switch d.Kind {
	case defx.ArrayType:
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("List<%s>", t.Sub.Name)
	case defx.BoolType:
		t.Name = "bool"
	case defx.BytesType:
		t.Name = "Uint8List"
	case defx.FloatType:
		t.Name = "double"
	case defx.IntType:
		t.Name = "int"
	case defx.MapType:
		t.Key = newType(d.Key, level+1)
		t.Sub = newType(d.Sub, level+1)
		t.Name = fmt.Sprintf("Map<%s, %s>", t.Key.Name, t.Sub.Name)
	case defx.ModelType:
		t.Name = d.Name
	case defx.RefType:
		t.Name = "ModelRef"
	case defx.StringType:
		t.Name = "String"
	case defx.TimeType:
		t.Name = "DateTime"
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
	if t.Level > 0 {
		target += fmt.Sprintf("%d", t.Level-1)
	}
	decoder := "d"
	if t.Level > 0 {
		decoder += fmt.Sprintf("%d", t.Level-1)
	}
	switch t.Kind {
	case defx.ArrayType:
		w.Put("final %s = <%s>[];", target, t.Sub.Name)
		w.Inc("{")
		{
			d := fmt.Sprintf("d%d", t.Level)
			i := fmt.Sprintf("i%d", t.Level)
			w.Put("final %s = %s.getArray(%s);", d, decoder, source)
			w.Inc("for (int %s = 0; %s < %s.length; %s++) {", i, i, d, i)
			{
				v := t.Sub.doDecode(w, "", "v")
				w.Put("%s.add(%s);", target, v)
			}
			w.Dec("}")
		}
		w.Dec("}")

	case defx.BoolType:
		w.Put("final %s = %s.getBool(%s);", target, decoder, source)

	case defx.BytesType:
		w.Put("final %s = %s.getBytes(%s);", target, decoder, source)

	case defx.FloatType:
		w.Put("final %s = %s.getFloat(%s);", target, decoder, source)

	case defx.IntType:
		w.Put("final %s = %s.getInt(%s);", target, decoder, source)

	case defx.MapType:
		w.Put("final %s = <%s,%s>{};", target, t.Key.Name, t.Sub.Name)
		w.Inc("{")
		{
			d := fmt.Sprintf("d%d", t.Level)
			i := fmt.Sprintf("i%d", t.Level)
			w.Put("final %s = %s.getMap(%s);", d, decoder, source)
			w.Inc("for (int %s = 0; %s < %s.length; %s++) {", i, i, d, i)
			{
				k := t.Key.doDecode(w, "", "k")
				v := t.Sub.doDecode(w, "", "v")
				w.Put("%s[%s] = %s;", target, k, v)
			}
			w.Dec("}")
		}
		w.Dec("}")

	case defx.ModelType:
		w.Put("%s %s;", t.Name, target)
		w.Inc("{")
		{
			d := fmt.Sprintf("d%d", t.Level)
			w.Put("final %s = %s.getObject(%s);", d, decoder, source)
			c := strcase.ToCamel(t.Name)
			w.Put("%s = %sCodec.decode(%s);", target, c, d)
		}
		w.Dec("}")

	case defx.RefType:
		w.Put("final %s = %s.getRef(%s);", target, decoder, source)

	case defx.StringType:
		w.Put("final %s = %s.getString(%s);", target, decoder, source)

	case defx.TimeType:
		w.Put("final %s = %s.getTime(%s);", target, decoder, source)
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
		return fmt.Sprintf("put%s(%s)", kind, args)
	}

	switch t.Kind {
	case defx.ArrayType:
		w.Inc("{")
		{
			e := fmt.Sprintf("e%d", t.Level)
			m := method("Array", fmt.Sprintf("%s.length", source))
			w.Put("final %s = %s.%s;", e, encoder, m)
			v := fmt.Sprintf("v%d", t.Level)
			w.Inc("for (final %s in %s) {", v, source)
			{
				t.Sub.doEncode(w, v, "")
			}
			w.Dec("}")
		}
		w.Dec("}")

	case defx.BoolType:
		m := method("Int", source)
		w.Put("%s.%s;", encoder, m)

	case defx.BytesType:
		m := method("Bytes", source)
		w.Put("%s.%s;", encoder, m)

	case defx.FloatType:
		m := method("Float", source)
		w.Put("%s.%s;", encoder, m)

	case defx.IntType:
		m := method("Int", source)
		w.Put("%s.%s;", encoder, m)

	case defx.MapType:
		w.Inc("{")
		{
			e := fmt.Sprintf("e%d", t.Level)
			m := method("Map", fmt.Sprintf("%s.length", source))
			w.Put("final %s = %s.%s;", e, encoder, m)
			p := fmt.Sprintf("p%d", t.Level)
			w.Inc("for (final %s in %s.entries) {", p, source)
			{
				k := fmt.Sprintf("%s.key", p)
				t.Key.doEncode(w, k, "")
				v := fmt.Sprintf("%s.value", p)
				t.Sub.doEncode(w, v, "")
			}
			w.Dec("}")
		}
		w.Dec("}")

	case defx.ModelType:
		w.Inc("{")
		{
			e := fmt.Sprintf("e%d", t.Level)
			m := method("Object", "")
			w.Put("final %s = %s.%s;", e, encoder, m)
			c := strcase.ToCamel(t.Name)
			w.Put("%sCodec.encode(%s, %s);", c, e, source)
		}
		w.Dec("}")

	case defx.RefType:
		m := method("Ref", source)
		w.Put("%s.%s;", encoder, m)

	case defx.StringType:
		m := method("String", source)
		w.Put("%s.%s;", encoder, m)

	case defx.TimeType:
		m := method("String", source)
		w.Put("%s.%s;", encoder, m)
	}
}
