package types

import (
	"fmt"
)

type objectGen struct {
	baseGen
	name string
}

func (g objectGen) Decode(from, to string) string {
	di := g.in("d")
	do := g.out("d")
	var m string
	if from == "" {
		m = "GetObject()"
	} else {
		m = fmt.Sprintf("GetObject(%q)", from)
	}
	g.Put("%s, err := %s.%s", do, di, m)
	g.Inc("if err != nil")
	{
		g.Put("return m, err")
	}
	g.Dec("}")

	vo := g.out(to)
	g.Put("%s, err := %sCodec.Decode(%s)", vo, g.name, do)
	g.Inc("if err != nil")
	{
		g.Put("return m, err")
	}
	g.Dec("}")

	return vo
}

func (g objectGen) Encode(from, to string) {
	ei := g.in("e")
	eo := g.out("e")
	var m string
	if to == "" {
		m = "PutObject()"
	} else {
		m = fmt.Sprintf("PutObject(%q)", to)
	}
	g.Put("%s, err := %s.%s", eo, ei, m)
	g.Inc("if err != nil")
	{
		g.Put("return err")
	}
	g.Dec("}")

	g.Put("err = %sCodec.Encode(%s, %s)", g.name, eo, from)
	g.Inc("if err != nil")
	{
		g.Put("return err")
	}
	g.Dec("}")
}

func (g objectGen) String() string {
	return fmt.Sprintf("%s", g.name)
}
