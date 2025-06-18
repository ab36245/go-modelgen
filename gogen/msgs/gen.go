package msgs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ab36245/go-modelgen/writer"

	"github.com/ab36245/go-modelgen/gogen/godefs"
)

func Generate(dir string, ds []godefs.Model) error {
	w := writer.WithPrefix("\t")
	doFile(w, ds)
	code := w.Code()

	// code, err := format(code)
	// if err != nil {
	// 	return fmt.Errorf("can't reformat code: %w", err)
	// }

	file := filepath.Join(dir, "msgs.go")
	fmt.Printf("Creating %s\n", file)
	if err := os.WriteFile(file, []byte(code), fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	return nil
}
