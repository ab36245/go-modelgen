package dartgen

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defx"
)

func newType(d *defx.Type) typeGen {
	switch d.Kind {
	case defx.ArrayType:
		return &arrayType{
			sub: newType(d.Sub),
		}
	case defx.BoolType:
		return &boolType{}
	case defx.BytesType:
		return &bytesType{}
	case defx.FloatType:
		return &floatType{}
	case defx.IntType:
		return &intType{}
	case defx.MapType:
		return &mapType{
			key: newType(d.Key),
			sub: newType(d.Sub),
		}
	case defx.ModelType:
		return &modelType{
			name: d.Name,
		}
	case defx.StringType:
		return &stringType{}
	case defx.TimeType:
		return &timeType{}
	default:
		text := fmt.Sprintf("bad type %d", d.Kind)
		panic(text)
	}
}

type typeGen interface {
	String() string
	doDecode(string, int)
	doEncode(string, int)
}

type arrayType struct {
	sub typeGen
}

func (t *arrayType) String() string {
	return fmt.Sprintf("List<%s>", t.sub)
}

func (t *arrayType) doDecode(name string, level int) {
	put("final %s = <%s>[];", name, t.sub)
	inc("{")
	{
		i := fmt.Sprintf("$i%d", level)
		n := fmt.Sprintf("$n%d", level)
		v := fmt.Sprintf("$v%d", level)

		put("final %s = mp.getArrayLength();", n)
		inc("for (var %s = 0; %s < %s; %s++) {", i, i, n, i)
		{
			t.sub.doDecode(v, level+1)
			put("%s.add(%s);", name, v)
		}
		dec("}")
	}
	dec("}")
}

func (t *arrayType) doEncode(name string, level int) {
	inc("{")
	{
		v := fmt.Sprintf("$v%d", level)

		put("mp.putArrayLength(%s.length);", name)
		inc("for (final %s in %s) {", v, name)
		{
			t.sub.doEncode(v, level+1)
		}
		dec("}")
	}
	dec("}")
}

type boolType struct{}

func (t *boolType) String() string {
	return "bool"
}

func (t *boolType) doDecode(name string, level int) {
	put("final %s = mp.getBool();", name)
}

func (t *boolType) doEncode(name string, level int) {
	put("mp.putBool(%s);", name)
}

type bytesType struct{}

func (t *bytesType) String() string {
	return "Uint8List"
}

func (t *bytesType) doDecode(name string, level int) {
	put("final %s = mp.getBytes();", name)
}

func (t *bytesType) doEncode(name string, level int) {
	put("mp.putBytes(%s);", name)
}

type floatType struct{}

func (t *floatType) String() string {
	return "double"
}

func (t *floatType) doDecode(name string, level int) {
	put("final %s = mp.getFloat();", name)
}

func (t *floatType) doEncode(name string, level int) {
	put("mp.putFloat(%s);", name)
}

type intType struct{}

func (t *intType) String() string {
	return "int"
}

func (t *intType) doDecode(name string, level int) {
	put("final %s = mp.getInt();", name)
}

func (t *intType) doEncode(name string, level int) {
	put("mp.putInt(%s);", name)
}

type mapType struct {
	key typeGen
	sub typeGen
}

func (t *mapType) String() string {
	return fmt.Sprintf("Map<%s, %s>", t.key, t.sub)
}

func (t *mapType) doDecode(name string, level int) {
	put("final %s = <%s, %s>{};", name, t.key, t.sub)
	inc("{")
	{
		i := fmt.Sprintf("$i%d", level)
		k := fmt.Sprintf("$k%d", level)
		n := fmt.Sprintf("$n%d", level)
		v := fmt.Sprintf("$v%d", level)

		put("final %s = mp.getMapLength();", n)
		inc("for (var %s = 0; %s < %s; %s++) {", i, i, n, i)
		{
			t.key.doDecode(k, level+1)
			t.sub.doDecode(v, level+1)
			put("%s[%s] = %s;", name, k, v)
		}
		dec("}")
	}
	dec("}")
}

func (t *mapType) doEncode(name string, level int) {
	inc("{")
	{
		k := fmt.Sprintf("$k%d", level)
		v := fmt.Sprintf("$v%d", level)

		put("mp.putMapLength(%s.length);", name)
		inc("%s.forEach((%s, %s) {", name, k, v)
		{
			t.key.doEncode(k, level+1)
			t.sub.doEncode(v, level+1)
		}
		dec("});")
	}
	dec("}")
}

type modelType struct {
	name string
}

func (t *modelType) String() string {
	return t.name
}

func (t *modelType) doDecode(name string, level int) {
	put("final %s = _decode%s(mp);", name, t.name)
}

func (t *modelType) doEncode(name string, level int) {
	put("_encode%s(mp, %s);", t.name, name)
}

type stringType struct{}

func (t *stringType) String() string {
	return "String"
}

func (t *stringType) doDecode(name string, level int) {
	put("final %s = mp.getString();", name)
}

func (t *stringType) doEncode(name string, level int) {
	put("mp.putString(%s);", name)
}

type timeType struct{}

func (t *timeType) String() string {
	return "DateTime"
}

func (t *timeType) doDecode(name string, level int) {
	put("final %s = mp.getTime();", name)
}

func (t *timeType) doEncode(name string, level int) {
	put("mp.putTime(%s);\n", name)
}
