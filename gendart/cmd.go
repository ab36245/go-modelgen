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
			Name:        "mp",
			Description: "Generate msgpack message codecs",
			Short:       "m",
			Binding:     cli.BoolFlag().Bind(&opts.Msgpack),
		},
		&cli.Option{
			Name:        "output",
			Description: "Output path for code",
			Short:       "o",
			Binding:     cli.StringSlice().Bind(&outputs),
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
	models, err := load.Models(sources)
	if err != nil {
		cmd.Fatal(1, "%v", err)
	}
	for _, output := range outputs {
		gen(output, models, opts)
	}
}
