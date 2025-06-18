package gogen

import (
	"fmt"

	"github.com/ab36245/go-strcase"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"

	"github.com/ab36245/go-modelgen/gogen/types"
)

func newField(w writer.GenWriter, d defx.Field) *fieldGen {
	typ := types.New(w, d.Type, 0)
	return &fieldGen{
		GenWriter: w,
		goName:    strcase.ToPascal(d.Name),
		name:      d.Name,
		typ:       typ,
	}
}

type fieldGen struct {
	writer.GenWriter
	goName string
	name   string
	typ    types.Gen
}

func (g *fieldGen) String() string {
	return fmt.Sprintf("%s %s", g.goName, g.typ)
}

func (g *fieldGen) doDecode() {
	g.Put("// %s", g.name)
	g.Inc("{")
	{
		v := g.typ.Decode(g.name, "v")
		g.Put("m.%s = %s", g.goName, v)
	}
	g.Dec("}")
}

func (g *fieldGen) doEncode() {
	g.Put("// %s", g.name)
	g.Inc("{")
	{
		n := fmt.Sprintf("m.%s", g.goName)
		g.typ.Encode(n, g.name)
	}
	g.Dec("}")
}
