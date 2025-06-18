package types

import "fmt"

type stringGen struct {
	baseGen
}

func (g stringGen) Decode(from, to string) string {
	di := g.in("d")
	vo := g.out(to)
	var m string
	if from == "" {
		m = "GetString()"
	} else {
		m = fmt.Sprintf("GetString(%q)", from)
	}
	g.Put("%s, err := %s.%s", vo, di, m)
	g.Inc("if err != nil")
	{
		g.Put("return m, err")
	}
	g.Dec("}")
	return vo
}

func (g stringGen) Encode(from, to string) {
	ei := g.in("e")
	var m string
	if to == "" {
		m = fmt.Sprintf("PutString(%s)", from)
	} else {
		m = fmt.Sprintf("PutString(%q, %s)", to, from)
	}
	g.Put("err = %s.%s", ei, m)
	g.Inc("if err != nil")
	{
		g.Put("return err")
	}
	g.Dec("}")
}

func (g stringGen) String() string {
	return "string"
}
