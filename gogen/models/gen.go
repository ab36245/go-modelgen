package models

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ab36245/go-runner"

	"github.com/ab36245/go-modelgen/writer"

	"github.com/ab36245/go-modelgen/gogen/godefs"
)

func Generate(dir string, ds []godefs.Model) error {
	w := writer.WithPrefix("\t")
	doModels(w, ds)
	code := w.Code()

	// code, err := format(code)
	// if err != nil {
	// 	return fmt.Errorf("can't reformat code: %w", err)
	// }

	file := filepath.Join(dir, "models.go")
	fmt.Printf("Creating %s\n", file)
	if err := os.WriteFile(file, []byte(code), fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	return nil
}

func format(code string) (string, error) {
	var output []byte
	cmd := runner.New("gofmt")
	cmd.Stdin(strings.NewReader(code))
	cmd.Stdout(runner.CaptureOutput(&output))
	err := cmd.Run()
	if err != nil {
		return "", nil
	}
	return string(output), nil
}
