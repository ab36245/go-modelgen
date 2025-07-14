module github.com/ab36245/go-modelgen

go 1.24.4

replace github.com/ab36245/go-source => ../go-source

replace github.com/ab36245/go-writer => ../go-writer

replace github.com/ab36245/go-cli => ../go-cli

require (
	github.com/ab36245/go-cli v0.0.0-20250514074543-660d55bcd3e5
	github.com/ab36245/go-runner v0.0.1
	github.com/ab36245/go-writer v0.0.0-20250619012835-04848829953b
)

require github.com/ab36245/go-strcase v0.0.0-20250613073624-43f95b65dcc2

require golang.org/x/mod v0.25.0

require (
	github.com/ab36245/go-source v0.0.0-20250610102038-4f637c704786
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
)

replace github.com/ab36245/go-runner => ../go-runner

replace github.com/ab36245/go-strcase => ../go-strcase

replace github.com/ab36245/go-model => ../go-model

replace github.com/ab36245/go-msgpack => ../go-msgpack
