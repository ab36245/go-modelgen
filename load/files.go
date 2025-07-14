package load

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/parser"
)

func filesToModels(models *[]defs.Model, files []string) error {
	for _, file := range files {
		if err := fileToModels(models, file); err != nil {
			return err
		}
	}
	return nil
}

func fileToModels(models *[]defs.Model, file string) error {
	fmt.Printf("Loading %s\n", file)
	parser, err := parser.NewParser(file)
	if err != nil {
		return err
	}
	list, err := parser.Parse()
	if err != nil {
		msg1 := fmt.Sprintf("error parsing %q", file)
		msg2 := err.Error()
		return fmt.Errorf("%s\n%s", msg1, msg2)
	}
	*models = append(*models, list...)
	return nil
}
