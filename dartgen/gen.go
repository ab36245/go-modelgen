package dartgen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ab36245/go-modelgen/defx"
)

func Generate(ds []defx.Model, path string, opts Opts) error {
	var dir string
	var file string
	if strings.HasSuffix(path, ".dart") {
		dir = filepath.Dir(path)
		file = path
	} else {
		dir = filepath.Join(path, "models")
		file = filepath.Join(dir, "models.dart")
	}
	fmt.Printf("Creating %s\n", file)
	if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}

	w.Clear()
	msgs := newModels(ds, opts)
	code := msgs.code()

	if err := os.WriteFile(file, []byte(code), fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	return nil
}
