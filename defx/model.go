package defx

import (
	"github.com/ab36245/go-writer"
)

type Model struct {
	Fields []Field
	Id     int
	Name   string
}

func (m Model) String() string {
	return writer.Reflect(m)
}
