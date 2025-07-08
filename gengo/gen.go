package gengo

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ab36245/go-runner"

	"github.com/ab36245/go-modelgen/defs"
)

func Generate(path string, ds []defs.Model, opts Opts) error {
	dir := filepath.Join(path, "models")
	if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
		return fmt.Errorf("can't create %s: %w", dir, err)
	}
	ms := newModels(ds)
	if err := genModels(dir, ms, opts); err != nil {
		return err
	}
	if opts.Db {
		if err := genDb(dir, ms, opts); err != nil {
			return err
		}
	}
	if opts.Msgpack {
		if err := genMsgpack(dir, ms, opts); err != nil {
			return err
		}
	}
	return nil
}

func genSave(dir string, name string, opts Opts, code string) error {
	if opts.Reformat {
		var err error
		code, err = genFormat(code)
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

func genTypes(ms []Model) map[defs.TypeKind]bool {
	set := make(map[defs.TypeKind]bool)

	var check func(*Type)
	check = func(t *Type) {
		switch t.Kind {
		case defs.ArrayType:
			check(t.Sub)
		case defs.MapType:
			check(t.Key)
			check(t.Sub)
		default:
			set[t.Kind] = true
		}
	}

	for _, m := range ms {
		for _, f := range m.Fields {
			check(f.Type)
		}
	}
	return set
}
