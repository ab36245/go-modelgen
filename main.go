package main

import (
	"github.com/ab36245/go-cli"

	"github.com/ab36245/go-modelgen/gendart"
	"github.com/ab36245/go-modelgen/gengo"
)

func main() {
	cmd := cli.Command{
		Subcommands: []cli.Command{
			gendart.Command,
			gengo.Command,
		},
	}
	cmd.Run()
}
