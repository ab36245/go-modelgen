package dartgen

import (
	"github.com/ab36245/go-modelgen/defx"
)

func newModel(d defx.Model) *modelGen {
	fields := []*fieldGen{}
	for _, f := range d.Fields {
		fields = append(fields, newField(f))
	}
	return &modelGen{
		fields: fields,
		id:     d.Id,
		name:   d.Name,
	}
}

type modelGen struct {
	fields []*fieldGen
	id     int
	name   string
}

func (g *modelGen) doClass() {
	put("// id %v", g.id)
	inc("class %s {", g.name)
	{
		for _, f := range g.fields {
			put("%s", f)
		}

		put("")
		if len(g.fields) == 0 {
			put("const %s();", g.name)
		} else {
			inc("const %s({", g.name)
			{
				for _, f := range g.fields {
					put("required this.%s,", f.name)
				}
			}
			dec("});")
		}

		put("")
		put("@override")
		inc("String toString() =>")
		{
			inc("ObjectWriter('%s')", g.name)
			{
				for _, f := range g.fields {
					put(".field('%s', %s)", f.name, f.name)
				}
				put(".toString();")
			}
			dec("")
		}
		dec("")
	}
	dec("}")
}

func (g *modelGen) doCodec() {
	put("// id %s", g.id)

	inc("%s _decode%s(MsgPackDecoder mp) {", g.name, g.name)
	if len(g.fields) == 0 {
		put("return const %s();", g.name)
	} else {
		inc("try {")
		{
			for _, f := range g.fields {
				f.doDecode()
			}
			inc("return %s(", g.name)
			{
				for _, f := range g.fields {
					put("%s: %s,", f.name, f.name)
				}
			}
			dec(");")
		}
		dec("")
		inc("} on MsgPackException catch (e) {")
		{
			put("throw MsgsException(e.toString());")
		}
		dec("}")
	}
	dec("}")

	put("")

	inc("void _encode%s(MsgPackEncoder mp, %s m) {", g.name, g.name)
	{
		for _, f := range g.fields {
			f.doEncode()
		}
	}
	dec("}")
}
