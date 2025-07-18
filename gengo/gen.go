package gengo

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/ab36245/go-runner"
)

func genFormat(code string) (string, error) {
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

func genModule(path string) (string, string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	}
	root := filepath.VolumeName(abs) + string(filepath.Separator)
	top := abs
	var file string
	for {
		file = filepath.Join(top, "go.mod")
		if _, err := os.Stat(file); err == nil {
			break
		}
		file = ""
		if top == root {
			break
		}
		top = filepath.Dir(top)
	}
	if file == "" {
		return "", "", fmt.Errorf("can't find go.mod file above %s", abs)
	}
	bytes, err := os.ReadFile(file)
	if err != nil {
		return "", "", fmt.Errorf("can't read %s file: %w", file, err)
	}
	info, err := modfile.Parse(file, bytes, nil)
	if err != nil {
		return "", "", fmt.Errorf("can't parse %s file: %w", file, err)
	}
	name := info.Module.Mod.Path
	more := strings.TrimPrefix(abs, top)
	more = strings.TrimPrefix(more, string(filepath.Separator))
	return name, more, nil
}

func genSave(opts Opts, name string, code string) error {
	if opts.Reformat {
		var err error
		code, err = genFormat(code)
		if err != nil {
			return fmt.Errorf("can't reformat code: %w", err)
		}
	}
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
