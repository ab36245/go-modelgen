package defs

import (
	"fmt"
)

type Field struct {
	Name string
	Type *Type
}

func (f Field) String() string {
	return fmt.Sprintf("%s: %s", f.Name, f.Type)
}
