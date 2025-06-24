package gengo

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/writer"
)

func genDb(dir string, ms []Model, opts Opts) error {
	w := writer.WithPrefix("\t")
	w.Put("// WARNING!")
	w.Put("// This code was generated automatically.")
	w.Put("package models")
	w.Put("")
	w.Inc("import (")
	{
		w.Put("\"fmt\"")
		w.Put("")
		w.Put("\"go.mongodb.org/mongo-driver/v2/bson\"")
		w.Put("")
		w.Put("\"github.com/ab36245/go-db\"")
		w.Put("\"github.com/ab36245/go-model\"")
	}
	w.Dec(")")
	for _, m := range ms {
		w.Put("")
		dbModel(w, m)
	}
	return genSave(dir, "db.go", opts, w.Code())
}

func dbModel(w writer.GenWriter, m Model) {
	w.Inc("var %sDbCodec = db.Codec[%s]{", m.Name, m.Name)
	{
		dbDecodeModel(w, m)
		dbEncodeModel(w, m)
	}
	w.Dec("}")
}

func dbDecodeModel(w writer.GenWriter, m Model) {
	w.Inc("Decode: func(d bson.M) (%s, error) {", m.Name)
	{
		w.Put("m := %s{}", m.Name)
		if len(m.Fields) > 0 {
			w.Put("var ok bool")
			for _, f := range m.Fields {
				dbDecodeField(w, f)
			}
		}
		w.Put("return m, nil")
	}
	w.Dec("},")
}

func dbDecodeField(w writer.GenWriter, f Field) {
	w.Inc("{")
	{
		source := fmt.Sprintf("d[%q]", f.Orig)
		target := dbDecodeType(w, f.Type, source, "v")
		w.Put("m.%s = %s", f.Name, target)
	}
	w.Dec("}")
}

func dbDecodeType(w writer.GenWriter, t *Type, source, target string) string {
	d := t.varName("d")
	v := t.varName(target)
	switch t.Kind {
	case defs.ArrayType:
		dbType := "bson.A"
		w.Put("var %s %s", d, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", d, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")
		w.Put("%s := make([]%s, len(%s))", v, t.Sub.Name, d)
		i := t.varName("i")
		e := t.varName("e")
		w.Inc("for %s, %s := range %s {", i, e, d)
		{
			o := dbDecodeType(w, t.Sub, e, "v")
			w.Put("%s[%s] = %s", v, i, o)
		}
		w.Dec("}")

	case defs.BoolType:
		dbType := "bool"
		w.Put("var %s %s", v, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", v, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")

	case defs.BytesType:
		dbType := "[]byte"
		w.Put("var %s %s", v, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", v, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")

	case defs.FloatType:
		dbType := "float64"
		w.Put("var %s %s", v, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", v, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")

	case defs.IntType:
		dbType := "int32"
		w.Put("var %s %s", d, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", d, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")
		w.Put("%s := int(%s)", v, d)

	case defs.MapType:
		dbType := "bson.M"
		w.Put("var %s %s", d, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", d, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")
		w.Put("%s := make(map[%s]%s, len(%s))", v, t.Key.Name, t.Sub.Name, d)
		k := t.varName("k")
		e := t.varName("e")
		w.Inc("for %s, %s := range %s {", k, e, d)
		{
			kn := dbDecodeType(w, t.Key, k, "k")
			en := dbDecodeType(w, t.Sub, e, "e")
			w.Put("%s[%s] = %s", v, kn, en)
		}
		w.Dec("}")

	case defs.ModelType:
		dbType := "bson.M"
		w.Put("var %s %s", d, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", d, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")
		w.Put("var %s %s", v, t.Name)
		w.Put("var err error")
		w.Inc("if %s, err = %sDbCodec.Decode(%s); err != nil {", v, t.Name, d)
		{
			w.Put("return m, err")
		}
		w.Dec("}")

	case defs.RefType:
		dbType := "bson.ObjectID"
		w.Put("var %s %s", d, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", d, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")
		w.Put("%s := model.Ref(%s.Hex())", v, d)

	case defs.StringType:
		dbType := "string"
		w.Put("var %s %s", v, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", v, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")

	case defs.TimeType:
		dbType := "bson.DateTime"
		w.Put("var %s %s", d, dbType)
		w.Inc("if %s, ok = %s.(%s); !ok {", d, source, dbType)
		{
			w.Put("return m, fmt.Errorf(\"invalid %s\")", dbType)
		}
		w.Dec("}")
		w.Put("%s := %s.Time()", v, d)
	}
	return v
}

func dbEncodeModel(w writer.GenWriter, m Model) {
	w.Inc("Encode: func(m %s) (bson.M, error) {", m.Name)
	{
		w.Put("e := make(bson.M, %d)", len(m.Fields))
		if len(m.Fields) > 0 {
			for _, f := range m.Fields {
				dbEncodeField(w, f)
			}
		}
		w.Put("return e, nil")
	}
	w.Dec("},")
}

func dbEncodeField(w writer.GenWriter, f Field) {
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
		w.Put("%s := int32(%s)", v, source)
		return v

	case defs.ModelType:
		w.Put("%s, err := %sDbCodec.Encode(%s)", v, t.Name, source)
		w.Inc("if err != nil {")
		{
			w.Put("return nil, err")
		}
		w.Dec("}")

	case defs.RefType:
		w.Put("%s, err := bson.ObjectIDFromHex(string(%s))", v, source)
		w.Inc("if err != nil {")
		{
			w.Put("return nil, err")
		}
		w.Dec("}")
		return v

	case defs.TimeType:
		w.Put("%s := bson.NewDateTimeFromTime(%s)", v, source)
		return v

	case defs.StringType:
		return source
	}
	return v
}
