package load

import (
	"github.com/ab36245/go-modelgen/defs"
)

const defExtension = ".def"

func Models(paths []string) ([]defs.Model, error) {
	if paths == nil {
		paths = []string{"."}
	}
	files := make([]string, 0)
	if err := pathsToFiles(&files, paths); err != nil {
		return nil, err
	}
	models := make([]defs.Model, 0)
	if err := filesToModels(&models, files); err != nil {
		return nil, err
	}
	return models, nil
}
