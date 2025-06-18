package gogen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/gogen/godefs"

	"github.com/ab36245/go-modelgen/gogen/models"
	"github.com/ab36245/go-modelgen/gogen/msgs"
)

func Generate(ds []defx.Model, path string, opts Opts) error {
	gs := godefs.Models(ds)

	dir := filepath.Join(path, "models")
	fmt.Printf("Creating %s\n", dir)
	if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	if err := models.Generate(dir, gs); err != nil {
		return err
	}
	if err := msgs.Generate(dir, gs); err != nil {
		return err
	}
	return nil
}
