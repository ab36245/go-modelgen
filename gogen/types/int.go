package types

import "fmt"

type intGen struct {
	baseGen
}

func (g intGen) Decode(from, to string) string {
	di := g.in("d")
	vo := g.out(to)
	var m string
	if from == "" {
		m = "GetInt()"
	} else {
		m = fmt.Sprintf("GetInt(%q)", from)
	}
	g.Put("%s, err := %s.%s", vo, di, m)
	g.Inc("if err != nil")
	{
		g.Put("return m, err")
	}
	g.Dec("}")
	return vo
}

func (g intGen) Encode(from, to string) {
	ei := g.in("e")
	var m string
	if to == "" {
		m = fmt.Sprintf("PutInt(%s)", from)
	} else {
		m = fmt.Sprintf("PutInt(%q, %s)", to, from)
	}
	g.Put("err = %s.%s", ei, m)
	g.Inc("if err != nil")
	{
		g.Put("return err")
	}
	g.Dec("}")
}

func (g intGen) String() string {
	return "int"
}
