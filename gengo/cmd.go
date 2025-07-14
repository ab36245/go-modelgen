package gengo

import (
	"github.com/ab36245/go-cli"

	"github.com/ab36245/go-modelgen/load"
)

var (
	opts    = Opts{}
	sources []string
)

var Command = cli.Command{
	Name: "go",
	Options: cli.Options{
		&cli.Option{
			Name:        "db",
			Short:       "d",
			Description: "Generate codecs for db",
			Binding:     cli.BoolFlag().Bind(&opts.Db),
		},
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
		&cli.Option{
			Name:        "reformat",
			Short:       "r",
			Description: "Reformat the generated code with gofmt",
			Binding:     cli.BoolFlag().Bind(&opts.Reformat),
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

func run(cmd *cli.Command, _ []string) {
	// info, err := getInfo(opts)
	// if err != nil {
	// 	cmd.Fatal(1, "%v", err)
	// }
	// fmt.Printf("info %#v\n", info)

	ds, err := load.Models(sources)
	if err != nil {
		cmd.Fatal(1, "%v", err)
	}
	ms := newModels(ds)

	if err := genModels(opts, ms); err != nil {
		cmd.Fatal(1, "%v", err)
	}
	if err := genDb(opts, ms); err != nil {
		cmd.Fatal(1, "%v", err)
	}
	if err := genMsgpack(opts, ms); err != nil {
		cmd.Fatal(1, "%v", err)
	}
}
