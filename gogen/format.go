package gogen

import (
	"strings"

	"github.com/ab36245/go-runner"
)

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
