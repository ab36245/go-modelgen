package gendart

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/writer"
)

const msgpackExtType = 0

func genMsgpack(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("  ")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	msgpackImports(w, ms)
	w.Put("")
	msgpackDecodeFunc(w, ms)
	w.Put("")
	msgpackEncodeFunc(w, ms)
	for _, m := range ms {
		w.Put("")
		msgpackDecodeModel(w, m)
		w.Put("")
		msgpackEncodeModel(w, m)
	}
	w.Put("")
	return genSave(dir, "msgpack.dart", opts, w.Code())
}

func msgpackImports(w writer.GenWriter, ms []Model) {
	names := map[string]bool{
		"package:dart_msgpack/dart_msgpack.dart": true,
		"models.dart":                            true,
	}
	types := genTypes(ms)
	_ = types
	// if types[defs.OptionType] || types[defs.RefType] {
	// 	names["github.com/ab36245/go-model"] = true
	// }
	for name := range names {
		w.Put("import '%s';", name)
	}
}

func msgpackDecodeFunc(w writer.GenWriter, ms []Model) {
	w.Inc("dynamic decodeMsgpack(MsgPackDecoder mpd) {")
	{
		w.Put("final (typ, id) = mpd.getExtUint();")
		w.Inc("if (typ != %d) {", msgpackExtType)
		{
			w.Put("throw Exception('unexpected extension type $typ');")
		}
		w.Dec("}")
		w.Inc("return switch (id) {")
		{
			for _, m := range ms {
				w.Put("%d => _decode%sMsgpack(mpd),", m.Id, m.Name)
			}
			w.Put("_ => throw Exception('unknown model id $id'),")
		}
		w.Dec("};")
	}
	w.Dec("}")
}

func msgpackEncodeFunc(w writer.GenWriter, ms []Model) {
	w.Inc("void encodeMsgpack(MsgPackEncoder mpe, dynamic v) {")
	{
		w.Inc("switch (v) {")
		for _, m := range ms {
			w.Inc("case %s():", m.Name)
			{
				w.Put("_encode%sMsgpack(mpe, v);", m.Name)
			}
			w.Dec("")
		}
		w.Inc("default:")
		{
			w.Put("throw Exception('unknown model ${v.runtimeType}');")
		}
		w.Dec("")
		w.Dec("}")
	}
	w.Dec("}")
}

func msgpackDecodeModel(w writer.GenWriter, m Model) {
	w.Inc("%s _decode%sMsgpack(MsgPackDecoder mpd) {", m.Name, m.Name)
	{
		for i, f := range m.Fields {
			if i > 0 {
				w.Put("")
			}
			w.Put("// %s", f.Name)
			msgpackDecodeField(w, f)
		}
		w.Put("")
		w.Inc("return %s(", m.Name)
		{
			for _, f := range m.Fields {
				w.Put("%s: %s,", f.Name, f.Name)
			}
		}
		w.Dec(");")
	}
	w.Dec("}")
}

func msgpackDecodeField(w writer.GenWriter, f Field) {
	msgpackDecodeType(w, f.Type, f.Name)
}

func msgpackDecodeType(w writer.GenWriter, t *Type, target string) string {
	doGet := func(local, method string) {
		w.Put("final %s = mpd.get%s();", local, method)
	}

	v := target
	if t.Level > 0 {
		v += fmt.Sprintf("%d", t.Level-1)
	}

	d := t.varName("d") // raw msgpack data
	// e := t.varName("e")
	i := t.varName("i")
	// k := t.varName("k")
	n := t.varName("n")

	switch t.Kind {
	case defs.ArrayType:
		w.Put("final %s = <%s>[];", v, t.Sub.Name)
		w.Inc("{")
		{
			doGet(n, "ArrayLength")
			w.Inc("for (int %s = 0; %s < %s; %s++) {", i, i, n, i)
			{
				e := msgpackDecodeType(w, t.Sub, "e")
				w.Put("%s.add(%s);", v, e)
			}
			w.Dec("}")
		}
		w.Dec("}")

	case defs.BoolType:
		doGet(v, "Bool")

	case defs.BytesType:
		doGet(v, "Bytes")

	case defs.FloatType:
		doGet(v, "Float")

	case defs.IntType:
		doGet(v, "Int")

	case defs.MapType:
		w.Put("final %s = <%s, %s>{};", target, t.Key.Name, t.Sub.Name)
		w.Inc("{")
		{
			doGet(n, "MapLength")
			w.Inc("for (int %s = 0; %s < %s; %s++) {", i, i, n, i)
			{
				k := msgpackDecodeType(w, t.Key, "k")
				e := msgpackDecodeType(w, t.Sub, "e")
				w.Put("%s[%s] = %s;", v, k, e)
			}
			w.Dec("}")
		}
		w.Dec("}")

	case defs.ModelType:
		w.Put("final %s = _decode%sMsgpack(mpd);", v, t.Name)

	case defs.RefType:
		doGet(d, "String")
		w.Put("final %s = ModelRef(%s);", v, d)

	case defs.OptionType:
		w.Put("%s %s;", t.Name, v)
		w.Inc("if (mpd.ifNil()) {")
		{
			w.Put("%s = null;", v)
		}
		w.Dec("")
		w.Inc("} else {")
		e := msgpackDecodeType(w, t.Sub, "v")
		w.Put("%s = %s;", v, e)
		w.Dec("}")

	case defs.StringType:
		doGet(v, "String")

	case defs.TimeType:
		doGet(v, "Time")

	default:
		panic(fmt.Sprintf("unknown type to decode %d", t.Kind))
	}
	return v
}

func msgpackEncodeModel(w writer.GenWriter, m Model) {
	w.Inc("void _encode%sMsgpack(MsgPackEncoder mpe, %s m) {", m.Name, m.Name)
	{
		w.Put("mpe.putExtUint(%d, %d);", msgpackExtType, m.Id)
		for _, f := range m.Fields {
			w.Put("")
			w.Put("// %s", f.Name)
			msgpackEncodeField(w, f)
		}
	}
	w.Dec("}")
}

func msgpackEncodeField(w writer.GenWriter, f Field) {
	source := fmt.Sprintf("m.%s", f.Name)
	msgpackEncodeType(w, f.Type, source)
}

func msgpackEncodeType(w writer.GenWriter, t *Type, source string) {
	doPut := func(method, local string) {
		w.Put("mpe.put%s(%s);", method, local)
	}

	switch t.Kind {
	case defs.ArrayType:
		doPut("ArrayLength", fmt.Sprintf("%s.length", source))
		e := t.varName("e")
		w.Inc("for (final %s in %s) {", e, source)
		{
			msgpackEncodeType(w, t.Sub, e)
		}
		w.Dec("}")

	case defs.BoolType:
		doPut("Bool", source)

	case defs.BytesType:
		doPut("Bytes", source)

	case defs.FloatType:
		doPut("Float", source)

	case defs.IntType:
		doPut("Int", source)

	case defs.MapType:
		doPut("MapLength", fmt.Sprintf("%s.length", source))
		e := t.varName("e")
		w.Inc("for (final %s in %s.entries) {", e, source)
		{
			k := fmt.Sprintf("%s.key", e)
			msgpackEncodeType(w, t.Key, k)
			o := fmt.Sprintf("%s.value", e)
			msgpackEncodeType(w, t.Sub, o)
		}
		w.Dec("}")

	case defs.ModelType:
		w.Put("_encode%sMsgpack(mpe, %s);", t.Name, source)

	case defs.OptionType:
		w.Inc("if (%s != null) {", source)
		{
			msgpackEncodeType(w, t.Sub, fmt.Sprintf("%s!", source))
		}
		w.Dec("")
		w.Inc("} else {")
		{
			doPut("Nil", "")
		}
		w.Dec("}")

	case defs.RefType:
		doPut("String", fmt.Sprintf("%s.id", source))

	case defs.StringType:
		doPut("String", source)

	case defs.TimeType:
		doPut("Time", source)

	default:
		panic(fmt.Sprintf("unknown type to encode %d", t.Kind))
	}
}
