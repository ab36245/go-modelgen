package gogen

import "github.com/ab36245/go-modelgen/defx"

func newModels(ds []defx.Model) Models {
	return Models(doMap(ds, newModel))
}

type Models []Model
