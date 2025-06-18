package dartgen

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defx"
)

func newField(d defx.Field) *fieldGen {
	return &fieldGen{
		name:  d.Name,
		type_: newType(d.Type),
	}
}

type fieldGen struct {
	name  string
	type_ typeGen
}

func (g *fieldGen) String() string {
	return fmt.Sprintf("final %s %s;", g.type_, g.name)
}

func (g *fieldGen) doDecode() {
	g.type_.doDecode(g.name, 0)
}

func (g *fieldGen) doEncode() {
	g.type_.doEncode(fmt.Sprintf("m.%s", g.name), 0)
}
