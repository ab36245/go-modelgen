package godefs

import "github.com/ab36245/go-modelgen/defx"

type Model struct {
	Fields []Field
	Id     int
	Name   string
}

func newModel(d defx.Model) Model {
	return Model{
		Fields: doMap(d.Fields, newField),
		Id:     d.Id,
		Name:   d.Name,
	}
}
