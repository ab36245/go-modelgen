package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ab36245/go-cli"

	"github.com/ab36245/go-modelgen/dartgen"
	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/gengo"
	"github.com/ab36245/go-modelgen/parser"
)

const defExtension = ".def"

var dartPaths []string
var dartOpts dartgen.Opts

var goPaths []string
var goOpts gengo.Opts

var sources []string

func main() {
	cmd := cli.Command{
		Options: cli.Options{
			&cli.Option{
				Name:        "dart-path",
				Description: "Output path for dart code",
				Short:       "d",
				Binding:     cli.StringSlice().Bind(&dartPaths),
			},
			&cli.Option{
				Name:        "go-format",
				Description: "Reformat the generated code with gofmt",
				Binding:     cli.BoolFlag().Bind(&goOpts.Reformat),
			},
			&cli.Option{
				Name:        "go-path",
				Description: "Output path for go code",
				Short:       "g",
				Binding:     cli.StringSlice().Bind(&goPaths),
			},
		},
		Params: cli.Params{
			&cli.Param{
				Name:        "sources",
				Description: "List of files or directories containing defs",
				Binding:     cli.StringSlice().Bind(&sources),
			},
		},
		OnRun: run,
	}
	cmd.Run()
}

func run(cmd *cli.Command, args []string) {
	if sources == nil {
		sources = []string{"."}
	}
	files := []string{}
	for _, source := range sources {
		info, err := os.Stat(source)
		if errors.Is(err, os.ErrNotExist) {
			cmd.Fatal(1, "%v", err)
		}
		if err != nil {
			cmd.Fatal(1, "%v", err)
		}
		if !info.IsDir() {
			files = append(files, source)
			continue
		}
		filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(d.Name(), defExtension) {
				return nil
			}
			files = append(files, path)
			return nil
		})
	}
	if len(files) == 0 {
		fmt.Printf("No %s files found\n", defExtension)
		os.Exit(0)
	}

	allDefs := []defx.Model{}
	for _, file := range files {
		fmt.Printf("Loading %s\n", file)
		parser, err := parser.NewParser(file)
		if err != nil {
			cmd.Fatal(1, "can't open %q: %v", file, err)
		}
		fileDefs, err := parser.Parse()
		if err != nil {
			cmd.Error("Error parsing %q", file)
			cmd.Error(err.Error())
			os.Exit(1)
		}
		allDefs = append(allDefs, fileDefs...)
	}
	if len(allDefs) == 0 {
		fmt.Printf("No definitions in the given files\n")
		os.Exit(0)
	}

	for _, dartPath := range dartPaths {
		dartgen.Generate(allDefs, dartPath, dartOpts)
	}

	for _, goPath := range goPaths {
		gengo.Generate(goPath, allDefs, goOpts)
	}
}
