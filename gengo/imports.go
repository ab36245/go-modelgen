package gengo

import (
	"fmt"
	"slices"
	"strings"
)

type Imports struct {
	standard   []string
	thirdParty []string
	ours       []string
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
			s += fmt.Sprintf("\t%q\n", name)
		}
	}
	add(i.standard)
	add(i.thirdParty)
	add(i.ours)
	if s != "" {
		s = fmt.Sprintf("import(\n%s)", s)
	}
	return s
}

func (i *Imports) add(name string) {
	var slice *[]string
	if strings.HasPrefix(name, "github.com/ab36245/") {
		slice = &i.ours
	} else if strings.HasPrefix(name, "github.com/aivoicesystems/") {
		slice = &i.ours
	} else if strings.Count(name, "/") < 3 {
		slice = &i.standard
	} else {
		slice = &i.thirdParty
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
