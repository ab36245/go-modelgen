package types

import (
	"fmt"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-modelgen/writer"
)

func New(w writer.GenWriter, t *defx.Type, level int) Gen {
	switch t.Kind {
	case defx.ArrayType:
		return arrayGen{
			baseGen: newBaseGen(w, level),
			sub:     New(w, t.Sub, level+1),
		}

	case defx.IntType:
		return intGen{
			baseGen: newBaseGen(w, level),
		}

	case defx.MapType:
		return mapGen{
			baseGen: newBaseGen(w, level),
			key:     New(w, t.Key, level+1),
			sub:     New(w, t.Sub, level+1),
		}

	case defx.ModelType:
		return objectGen{
			baseGen: newBaseGen(w, level),
			name:    t.Name,
		}

	case defx.RefType:
		return refGen{
			baseGen: newBaseGen(w, level),
		}

	case defx.StringType:
		return stringGen{
			baseGen: newBaseGen(w, level),
		}

	default:
		text := fmt.Sprintf("bad type %d", t.Kind)
		panic(text)
	}
}

type Gen interface {
	Decode(string, string) string
	Encode(string, string)
	String() string
}
