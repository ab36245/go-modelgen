package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ab36245/go-modelgen/defs"
	"github.com/ab36245/go-modelgen/gengo"
	"github.com/ab36245/go-modelgen/parser"
)

func TestSimpleGen(t *testing.T) {
	name := "simple"
	defs, err := load(name)
	if err != nil {
		t.Fatal(err)
	}
	err = genGo(name, defs)
	if err != nil {
		t.Fatal(err)
	}
}

func load(name string) ([]defs.Model, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't get working directory: %w", err)
	}
	file := fmt.Sprintf("%s.def", name)
	path := filepath.Join(dir, "defs", file)
	parser, err := parser.NewParser(path)
	if err != nil {
		return nil, fmt.Errorf("can't create parser: %w", err)
	}
	defs, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("parser failed: %w", err)
	}
	return defs, nil
}

func genGo(name string, defs []defs.Model) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("can't get working directory: %w", err)
	}
	path := filepath.Join(dir, "output", name, "go")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't create output directory: %w", err)
	}
	opts := gengo.Opts{}
	err = gengo.Generate(path, defs, opts)
	if err != nil {
		return fmt.Errorf("can't generate output: %w", err)
	}
	return nil
}
