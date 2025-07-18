package gendart

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func genSave(opts Opts, name string, code string) error {
	if err := os.MkdirAll(opts.Path, os.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", opts.Path, err)
	}

	file := filepath.Join(opts.Path, name)
	fmt.Printf("Creating %s\n", file)
	if err := os.WriteFile(file, []byte(code), fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", file, err)
	}
	return nil
}
