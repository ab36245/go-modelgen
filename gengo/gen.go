package gengo

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ab36245/go-modelgen/defx"
)

func Generate(path string, ds []defx.Model, opts Opts) error {
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
	if err := genDb(dir, ms, opts); err != nil {
		return err
	}
	if err := genMp(dir, ms, opts); err != nil {
		return err
	}
	return nil
}

func genSave(dir string, name string, opts Opts, code string) error {
	if opts.Reformat {
		var err error
		code, err = format(code)
		if err != nil {
			return fmt.Errorf("can't reformat code: %w", err)
		}
	}

	file := filepath.Join(dir, name)
	fmt.Printf("Creating %s\n", file)
	if err := os.WriteFile(file, []byte(code), fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	return nil
}
