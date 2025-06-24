package gendart

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ab36245/go-modelgen/defs"
)

func Generate(path string, ds []defs.Model, opts Opts) error {
	dir := filepath.Join(path, "models")
	if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	ms := doMap(ds, newModel)
	if err := genModels(dir, ms, opts); err != nil {
		return err
	}
	if err := genCodecs(dir, ms, opts); err != nil {
		return err
	}
	if err := genMsgs(dir, ms, opts); err != nil {
		return err
	}
	return nil
}

func genSave(dir string, name string, opts Opts, code string) error {
	file := filepath.Join(dir, name)
	fmt.Printf("Creating %s\n", file)
	if err := os.WriteFile(file, []byte(code), fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	return nil
}
