package types

import "fmt"

type mapGen struct {
	baseGen
	key Gen
	sub Gen
}

func (g mapGen) Decode(from, to string) string {
	di := g.in("d")
	do := g.out("d")
	var m string
	if from == "" {
		m = "GetMap()"
	} else {
		m = fmt.Sprintf("GetMap(%q)", from)
	}
	g.Put("%s, err := %s.%s", do, di, m)
	g.Inc("if err != nil")
	{
		g.Put("return m, err")
	}
	g.Dec("}")

	vo := g.out(to)
	g.Put("%s := make(map[%s]%s, %s.Length())", vo, g.key, g.sub, do)
	g.Inc("for range %s.Length() {", do)
	{
		ke := g.key.Decode("", "k")
		ve := g.sub.Decode("", "v")
		g.Put("%s[%s] = %s", vo, ke, ve)
	}
	g.Dec("}")

	return vo
}

func (g mapGen) Encode(from, to string) {
	ei := g.in("e")
	eo := g.out("e")
	var m string
	if to == "" {
		m = fmt.Sprintf("PutMap(len(%s))", from)
	} else {
		m = fmt.Sprintf("PutMap(%q, len(%s))", to, from)
	}
	g.Put("%s, err := %s.%s", eo, ei, m)
	g.Inc("if err != nil")
	{
		g.Put("return err")
	}
	g.Dec("}")

	ko := g.out("k")
	vo := g.out("v")
	g.Inc("for %s, %s := range %s {", ko, vo, from)
	{
		g.key.Encode(ko, "")
		g.sub.Encode(vo, "")
	}
	g.Dec("}")
}

func (g mapGen) String() string {
	return fmt.Sprintf("map[%s]%s", g.key, g.sub)
}
