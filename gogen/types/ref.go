package types

import "fmt"

type refGen struct {
	baseGen
}

func (g refGen) Decode(from, to string) string {
	di := g.in("d")
	vo := g.out(to)
	var m string
	if from == "" {
		m = "GetRef()"
	} else {
		m = fmt.Sprintf("GetRef(%q)", from)
	}
	g.Put("%s, err := %s.%s", vo, di, m)
	g.Inc("if err != nil")
	{
		g.Put("return m, err")
	}
	g.Dec("}")
	return vo
}

func (g refGen) Encode(from, to string) {
	ei := g.in("e")
	var m string
	if to == "" {
		m = fmt.Sprintf("PutRef(%s)", from)
	} else {
		m = fmt.Sprintf("PutRef(%q, %s)", to, from)
	}
	g.Put("err = %s.%s", ei, m)
	g.Inc("if err != nil")
	{
		g.Put("return err")
	}
	g.Dec("}")
}

func (g refGen) String() string {
	return "model.Ref"
}
