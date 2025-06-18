package gogen

import "fmt"

const (
	decoderName = "decoder"
	encoderName = "encoder"
	modelName   = "m"
)

func sourceName(name string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("%q", name)
}

func keyName(level int) string {
	return targetName("k", level)
}

func valueName(level int) string {
	return targetName("v", level)
}

func targetName(base string, level int) string {
	if level == 0 {
		return base
	}
	return fmt.Sprintf("%s%d", base, level)
}
