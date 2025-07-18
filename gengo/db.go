package gengo

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/writer"
)

func genDb(opts Opts, ms Models) error {
	w := writer.WithPrefix("\t")
	dbFile(w, ms)
	return genSave(opts, "db.go", w.Code())
}

func dbFile(w writer.GenWriter, ms Models) {
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("package models")
	w.Put("")
	dbImports(w, ms)
	for _, m := range ms.List {
		w.Put("")
		dbCodec(w, m)
	}
}

func dbImports(w writer.GenWriter, ms Models) {
	imports := &Imports{}
	imports.add("fmt")
	imports.add("go.mongodb.org/mongo-driver/v2/bson")
	imports.add("github.com/ab36245/go-db")
	if ms.Types.HasOption() {
		imports.add("github.com/ab36245/go-model")
	} else if ms.Types.HasRef() {
		imports.add("github.com/ab36245/go-model")
	}
	if ms.Types.HasTimeArray() || ms.Types.HasTimeMap() {
		imports.add("time")
	}
	w.Put(imports.String())
}

func dbCodec(w writer.GenWriter, m Model) {
	w.Inc("var %sDbCodec = db.Codec[%s]{", m.Name, m.Name)
	{
		dbDecode(w, m)
		dbEncode(w, m)
	}
	w.Dec("}")
}

func dbDecode(w writer.GenWriter, m Model) {
	w.Inc("Decode: func(d bson.M) (%s, error) {", m.Name)
	{
		w.Put("m := %s{}", m.Name)
		if len(m.Fields) > 0 {
			w.Put("var ok bool")
			for _, f := range m.Fields {
				dbDecodeField(w, f)
			}
		}
		w.Put("")
		w.Put("return m, nil")
	}
	w.Dec("},")
}

func dbDecodeField(w writer.GenWriter, f Field) {
	w.Put("")
	w.Put("// %s", f.Name)
	w.Inc("{")
	{
		source := fmt.Sprintf("d[%q]", f.Orig)
		target := dbDecodeType(w, f.Type, source, "v")
		w.Put("m.%s = %s", f.Name, target)
	}
	w.Dec("}")
}

func dbDecodeType(w writer.GenWriter, t *Type, source, target string) string {
	doGet := func(local, dbType string) {
		w.Put("var %s %s", local, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", local, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")
	}

	v := t.varName(target)
	d := t.varName("d") // raw bson data
	i := t.varName("i")
	e := t.varName("e")
	k := t.varName("k")

	switch t.Kind {
	case defs.ArrayType:
		doGet(d, "bson.A")
		w.Put("%s := make([]%s, len(%s))", v, t.Sub.Name, d)
		w.Inc("for %s, %s := range %s {", i, e, d)
		{
			e := dbDecodeType(w, t.Sub, e, "v")
			w.Put("%s[%s] = %s", v, i, e)
		}
		w.Dec("}")

	case defs.BoolType:
		doGet(v, "bool")

	case defs.BytesType:
		doGet(v, "[]byte")

	case defs.FloatType:
		doGet(v, "float64")

	case defs.IntType:
		doGet(d, "int32")
		w.Put("%s := int(%s)", v, d)

	case defs.MapType:
		doGet(d, "bson.M")
		w.Put("%s := make(map[%s]%s, len(%s))", v, t.Key.Name, t.Sub.Name, d)
		w.Inc("for %s, %s := range %s {", k, e, d)
		{
			k := dbDecodeType(w, t.Key, k, "k")
			e := dbDecodeType(w, t.Sub, e, "e")
			w.Put("%s[%s] = %s", v, k, e)
		}
		w.Dec("}")

	case defs.ModelType:
		doGet(d, "bson.M")
		w.Put("var %s %s", v, t.Name)
		w.Put("var err error")
		w.Inc("if %s, err = %sDbCodec.Decode(%s); err != nil {", v, t.Name, d)
		{
			w.Put("return m, err")
		}
		w.Dec("}")

	case defs.OptionType:
		w.Put("var %s model.Option[%s]", v, t.Sub.Name)
		w.Inc("if %s == nil {", source)
		{
			w.Put("%s = model.EmptyOption[%s]()", v, t.Sub.Name)
		}
		w.Dec("")
		w.Inc("} else {")
		{
			e := dbDecodeType(w, t.Sub, source, "e")
			w.Put("%s = model.NewOption(%s)", v, e)
		}
		w.Dec("}")

	case defs.RefType:
		doGet(d, "bson.ObjectID")
		w.Put("%s := model.Ref(%s.Hex())", v, d)

	case defs.StringType:
		doGet(v, "string")

	case defs.TimeType:
		doGet(d, "bson.DateTime")
		w.Put("%s := %s.Time()", v, d)
	}
	return v
}

func dbEncode(w writer.GenWriter, m Model) {
	w.Inc("Encode: func(m %s) (bson.M, error) {", m.Name)
	{
		w.Put("e := make(bson.M, %d)", len(m.Fields))
		if len(m.Fields) > 0 {
			for _, f := range m.Fields {
				dbEncodeField(w, f)
			}
		}
		w.Put("")
		w.Put("return e, nil")
	}
	w.Dec("},")
}

func dbEncodeField(w writer.GenWriter, f Field) {
	w.Put("")
	w.Put("// %s", f.Name)
	w.Inc("{")
	{
		source := fmt.Sprintf("m.%s", f.Name)
		target := dbEncodeType(w, f.Type, source, "v")
		w.Put("e[%q] = %s", f.Orig, target)
	}
	w.Dec("}")
}

func dbEncodeType(w writer.GenWriter, t *Type, source, target string) string {
	v := t.varName(target)
	switch t.Kind {
	case defs.ArrayType:
		dbType := "bson.A"
		w.Put("%s := make(%s, len(%s))", v, dbType, source)
		i := t.varName("i")
		e := t.varName("e")
		w.Inc("for %s, %s := range %s {", i, e, source)
		{
			o := dbEncodeType(w, t.Sub, e, "v")
			w.Put("%s[%s] = %s", v, i, o)
		}
		w.Dec("}")
		return v

	case defs.BoolType:
		return source

	case defs.BytesType:
		return source

	case defs.FloatType:
		return source

	case defs.IntType:
		return fmt.Sprintf("int32(%s)", source)

	case defs.MapType:
		dbType := "bson.M"
		w.Put("%s := make(%s, len(%s))", v, dbType, source)
		k := t.varName("k")
		e := t.varName("e")
		w.Inc("for %s, %s := range %s {", k, e, source)
		{
			k := dbEncodeType(w, t.Key, k, "k")
			e := dbEncodeType(w, t.Sub, e, "e")
			w.Put("%s[%s] = %s", v, k, e)
		}
		w.Dec("}")
		return v

	case defs.ModelType:
		w.Put("%s, err := %sDbCodec.Encode(%s)", v, t.Name, source)
		w.Inc("if err != nil {")
		{
			w.Put("return nil, err")
		}
		w.Dec("}")
		return v

	case defs.OptionType:
		w.Put("var %s any", v)
		w.Inc("if %s.IsSet() {", source)
		{
			s := fmt.Sprintf("%s.Value()", source)
			e := dbEncodeType(w, t.Sub, s, "e")
			w.Put("%s = %s", v, e)
		}
		w.Dec("")
		w.Inc("} else {")
		{
			w.Put("%s = nil", v)
		}
		w.Dec("}")
		return v

	case defs.RefType:
		w.Put("%s, err := bson.ObjectIDFromHex(string(%s))", v, source)
		w.Inc("if err != nil {")
		{
			w.Put("return nil, err")
		}
		w.Dec("}")
		return v

	case defs.StringType:
		return source

	case defs.TimeType:
		return fmt.Sprintf("bson.NewDateTimeFromTime(%s)", source)
	}

	panic(fmt.Sprintf("unknown data type %d", t.Kind))
}
