package gendart

import (
	"fmt"
	"slices"
	"strings"
)

type Imports struct {
	standard   []string
	thirdParty []string
	ours       []string
	local      []string
}

func (i *Imports) String() string {
	s := ""
	add := func(names []string) {
		if len(names) == 0 {
			return
		}
		if s != "" {
			s += "\n"
		}
		for _, name := range names {
			s += fmt.Sprintf("import '%s';\n", name)
		}
	}
	add(i.standard)
	add(i.thirdParty)
	add(i.ours)
	add(i.local)
	return s
}

func (i *Imports) add(name string) {
	var slice *[]string
	if strings.HasPrefix(name, "package:dart_") {
		slice = &i.ours
	} else if strings.HasPrefix(name, "package:") {
		slice = &i.thirdParty
	} else if strings.Contains(name, ":") {
		slice = &i.standard
	} else {
		slice = &i.local
	}
	index := len(*slice)
	for i, n := range *slice {
		if n == name {
			return
		}
		if n > name {
			index = i
			continue
		}
	}
	*slice = slices.Insert(*slice, index, name)
}
