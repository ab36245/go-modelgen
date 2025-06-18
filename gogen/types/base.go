package types

import (
	"fmt"

	"github.com/ab36245/go-modelgen/writer"
)

func newBaseGen(w writer.GenWriter, level int) baseGen {
	return baseGen{
		GenWriter: w,
		level:     level,
	}
}

type baseGen struct {
	writer.GenWriter
	level int
}

func (g baseGen) in(base string) string {
	if g.level == 0 {
		return base
	}
	return fmt.Sprintf("%s%d", base, g.level-1)
}

func (g baseGen) out(base string) string {
	return fmt.Sprintf("%s%d", base, g.level)
}

func (g baseGen) doDecodeErr() {
	g.Inc("if err != nil {")
	{
		g.Put("return m, err")
	}
	g.Dec("}")
}

func (g baseGen) doEncodeErr() {
	g.Inc("if err != nil {")
	{
		g.Put("return err")
	}
	g.Dec("}")
}
