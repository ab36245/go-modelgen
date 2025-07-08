package gengo

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/writer"
)

const msgpackExtType = 0

func genMsgpack(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("\t")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("package models")
	w.Put("")
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
	return genSave(dir, "msgpack.go", opts, w.Code())
}

func msgpackImports(w writer.GenWriter, ms []Model) {
	names := map[string]bool{
		"fmt":                           true,
		"github.com/ab36245/go-msgpack": true,
	}
	types := genTypes(ms)
	if types[defs.OptionType] || types[defs.RefType] {
		names["github.com/ab36245/go-model"] = true
	}
	if types[defs.TimeType] {
		names["time"] = true
	}
	if len(names) > 0 {
		w.Inc("import (")
		{
			for name := range names {
				w.Put("%q", name)
			}
		}
		w.Dec(")")
	}
}

func msgpackDecodeFunc(w writer.GenWriter, ms []Model) {
	w.Inc("func DecodeMsgpack(mpd *msgpack.Decoder) (any, error) {")
	{
		w.Put("typ, id, err := mpd.GetExtUint()")
		w.Inc("if err != nil {")
		{
			w.Put("return nil, err")
		}
		w.Dec("}")
		w.Inc("if typ != %d {", msgpackExtType)
		{
			w.Put("return nil, fmt.Errorf(\"unexpected extension type %d\", typ)")
		}
		w.Dec("}")

		w.Put("switch id {")
		for _, m := range ms {
			w.Inc("case %d:", m.Id)
			{
				w.Put("return decode%sMsgpack(mpd)", m.Name)
			}
			w.Dec("")
		}
		w.Inc("default:")
		{
			w.Put("return nil, fmt.Errorf(\"unknown model id %d\", id)")
		}
		w.Dec("")
		w.Put("}")
	}
	w.Dec("}")
}

func msgpackEncodeFunc(w writer.GenWriter, ms []Model) {
	w.Inc("func EncodeMsgpack(mpe *msgpack.Encoder, v any) error {")
	{
		w.Put("var err error")
		w.Put("switch v := v.(type) {")
		for _, m := range ms {
			w.Inc("case %s:", m.Name)
			{
				w.Put("err = encode%sMsgpack(mpe, v)", m.Name)
			}
			w.Dec("")
		}
		w.Inc("default:")
		{
			w.Put("err = fmt.Errorf(\"unknown model %T\", v)")
		}
		w.Dec("")
		w.Put("}")
		w.Put("return err")
	}
	w.Dec("}")
}

func msgpackDecodeModel(w writer.GenWriter, m Model) {
	param := "mpd"
	if len(m.Fields) == 0 {
		param = "_"
	}
	w.Inc("func decode%sMsgpack(%s *msgpack.Decoder) (%s, error) {", m.Name, param, m.Name)
	{
		w.Put("m := %s{}", m.Name)
		if len(m.Fields) > 0 {
			for _, f := range m.Fields {
				msgpackDecodeField(w, f)
			}
			w.Put("")
		}
		w.Put("return m, nil")
	}
	w.Dec("}")
}

func msgpackDecodeField(w writer.GenWriter, f Field) {
	w.Put("")
	w.Put("// %s", f.Name)
	w.Inc("{")
	{
		target := msgpackDecodeType(w, f.Type, "v")
		w.Put("m.%s = %s", f.Name, target)
	}
	w.Dec("}")
}

func msgpackDecodeType(w writer.GenWriter, t *Type, target string) string {
	doGet := func(local, method string) {
		w.Put("%s, err := mpd.Get%s()", local, method)
		w.Inc("if err != nil {")
		{
			w.Put("return m, err")
		}
		w.Dec("}")
	}

	v := t.varName(target)
	d := t.varName("d") // raw msgpack data
	// e := t.varName("e")
	i := t.varName("i")
	// k := t.varName("k")
	n := t.varName("n")

	switch t.Kind {
	case defs.ArrayType:
		doGet(n, "ArrayLength")
		w.Put("%s := make([]%s, %s)", v, t.Sub.Name, n)
		w.Inc("for %s := range %s {", i, n)
		{
			e := msgpackDecodeType(w, t.Sub, "v")
			w.Put("%s[%s] = %s", v, i, e)
		}
		w.Dec("}")

	case defs.BoolType:
		doGet(v, "Bool")

	case defs.BytesType:
		doGet(v, "Bytes")

	case defs.FloatType:
		doGet(v, "Float")

	case defs.IntType:
		doGet(d, "Int")
		w.Put("%s := int(%s)", v, d)

	case defs.MapType:
		doGet(n, "MapLength")
		w.Put("%s := make(map[%s]%s, %s)", v, t.Key.Name, t.Sub.Name, n)
		w.Inc("for range %s {", i, n)
		{
			k := msgpackDecodeType(w, t.Key, "k")
			e := msgpackDecodeType(w, t.Sub, "e")
			w.Put("%s[%s] = %s", v, k, e)
		}
		w.Dec("}")

	case defs.ModelType:
		w.Put("var %s %s", v, t.Name)
		w.Put("var err error")
		w.Inc("if %s, err = decode%sMsgpack(mpd); err != nil {", v, t.Name)
		{
			w.Put("return m, err")
		}
		w.Dec("}")

	case defs.RefType:
		doGet(d, "String")
		w.Put("%s := model.Ref(%s)", v, d)

	case defs.OptionType:
		w.Put("isnil, err := mpd.IfNil()")
		w.Inc("if err != nil {")
		{
			w.Put("return m, err")
		}
		w.Dec("}")
		w.Put("var %s model.Option[%s]", v, t.Sub.Name)
		w.Inc("if isnil {")
		{
			w.Put("%s = model.EmptyOption[%s]()", v, t.Sub.Name)
		}
		w.Dec("")
		w.Inc("} else {")
		{
			e := msgpackDecodeType(w, t.Sub, "e")
			w.Put("%s = model.NewOption(%s)", v, e)
		}
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
	param := "m"
	if len(m.Fields) == 0 {
		param = "_"
	}
	w.Inc("func encode%sMsgpack(mpe *msgpack.Encoder, %s %s) error {", m.Name, param, m.Name)
	{
		w.Put("mpe.PutExtUint(%d, %d)", msgpackExtType, m.Id)
		if len(m.Fields) > 0 {
			// w.Put("var err error")
			for _, f := range m.Fields {
				msgpackEncodeField(w, f)
			}
			w.Put("")
		}
		w.Put("return nil")
	}
	w.Dec("}")
}

func msgpackEncodeField(w writer.GenWriter, f Field) {
	w.Put("")
	w.Put("// %s", f.Name)
	w.Inc("{")
	{
		source := fmt.Sprintf("m.%s", f.Name)
		msgpackEncodeType(w, f.Type, source)
	}
	w.Dec("}")
}

func msgpackEncodeType(w writer.GenWriter, t *Type, source string) {
	doPut := func(method, local string) {
		switch method {
		case "Bytes", "String":
			w.Inc("if err := mpe.Put%s(%s) ; err != nil {", method, local)
			{
				w.Put("return err")
			}
			w.Dec("}")
		default:
			w.Put("mpe.Put%s(%s)", method, local)
		}
	}

	switch t.Kind {
	case defs.ArrayType:
		doPut("ArrayLength", fmt.Sprintf("uint32(len(%s))", source))
		e := t.varName("e")
		w.Inc("for _, %s := range %s {", e, source)
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
		doPut("Int", fmt.Sprintf("int64(%s)", source))

	case defs.MapType:
		doPut("MapLength", fmt.Sprintf("uint32(len(%s))", source))
		k := t.varName("k")
		e := t.varName("e")
		w.Inc("for %s, %s := range %s {", k, e, source)
		{
			msgpackEncodeType(w, t.Key, k)
			msgpackEncodeType(w, t.Sub, e)
		}
		w.Dec("}")

	case defs.ModelType:
		w.Inc("if err := encode%sMsgpack(mpe, %s); err != nil {", t.Name, source)
		{
			w.Put("return err")
		}
		w.Dec("}")

	case defs.OptionType:
		w.Inc("if %s.IsSet() {", source)
		{
			s := fmt.Sprintf("%s.Value()", source)
			msgpackEncodeType(w, t.Sub, s)
		}
		w.Dec("")
		w.Inc("} else {")
		{
			doPut("Nil", "")
		}
		w.Dec("}")

	case defs.RefType:
		doPut("String", fmt.Sprintf("string(%s)", source))

	case defs.StringType:
		doPut("String", source)

	case defs.TimeType:
		doPut("Time", source)

	default:
		panic(fmt.Sprintf("unknown type to encode %d", t.Kind))
	}
}
