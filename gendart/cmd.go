package gendart

import (
	"github.com/ab36245/go-cli"
	"github.com/ab36245/go-modelgen/load"
)

var (
	opts    Opts
	outputs []string
	sources []string
)

var Command = cli.Command{
	Name: "dart",
	Options: cli.Options{
		&cli.Option{
			Name:        "msgpack",
			Short:       "m",
			Description: "Generate codecs for msgpack",
			Binding:     cli.BoolFlag().Bind(&opts.Msgpack),
		},
		&cli.Option{
			Name:        "path",
			Short:       "p",
			Description: "Output path",
			Binding:     cli.String().Bind(&opts.Path),
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

func run(cmd *cli.Command, args []string) {
	ds, err := load.Models(sources)
	if err != nil {
		cmd.Fatal(1, "%v", err)
	}
	ms := newModels(ds)

	if err := genModels(opts, ms); err != nil {
		cmd.Fatal(1, "%v", err)
	}
	if err := genMsgpack(opts, ms); err != nil {
		cmd.Fatal(1, "%v", err)
	}
}
