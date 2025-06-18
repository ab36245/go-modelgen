package godefs

import "github.com/ab36245/go-modelgen/defx"

func Models(ds []defx.Model) []Model {
	return doMap(ds, newModel)
}
