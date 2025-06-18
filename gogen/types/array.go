package types

import "fmt"

type arrayGen struct {
	baseGen
	sub Gen
}

func (g arrayGen) Decode(from, to string) string {
	di := g.in("d")
	do := g.out("d")
	var m string
	if from == "" {
		m = "GetArray()"
	} else {
		m = fmt.Sprintf("GetArray(%q)", from)
	}
	g.Put("%s, err := %s.%s", do, di, m)
	g.Inc("if err != nil")
	{
		g.Put("return m, err")
	}
	g.Dec("}")

	vo := g.out(to)
	g.Put("%s := make([]%s, %s.Length())", vo, g.sub, do)
	io := g.out("i")
	g.Inc("for %s := range %s.Length() {", io, do)
	{
		ve := g.sub.Decode("", "v")
		g.Put("%s[%s] = %s", io, vo, ve)
	}
	g.Dec("}")

	return vo
}

func (g arrayGen) Encode(from, to string) {
	ei := g.in("e")
	eo := g.out("e")
	var m string
	if from == "" {
		m = fmt.Sprintf("PutArray(len(%s))", from)
	} else {
		m = fmt.Sprintf("PutArray(%q, len(%s))", to, from)
	}
	g.Put("%s, err := %s.%s", eo, ei, m)
	g.Inc("if err != nil")
	{
		g.Put("return err")
	}
	g.Dec("}")

	vo := g.out("v")
	g.Inc("for _, %s := range %s {", vo, from)
	{
		g.sub.Encode(vo, "")
	}
	g.Dec("}")
}

func (g arrayGen) String() string {
	return fmt.Sprintf("[]%s", g.sub)
}
