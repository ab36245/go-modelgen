package defs

func GetTypes(ms []Model) Types {
	set := make(map[TypeKind]typeUsage)

	var check func(*Type, typeUsage)
	check = func(t *Type, usage typeUsage) {
		switch t.Kind {
		case ArrayType:
			check(t.Sub, tuArray)
		case MapType:
			check(t.Key, tuMap)
			check(t.Sub, tuMap)
		default:
			set[t.Kind] |= usage
		}
	}

	for _, m := range ms {
		for _, f := range m.Fields {
			check(f.Type, tuSimple)
		}
	}

	return Types{set}
}

type Types struct {
	set map[TypeKind]typeUsage
}

func (t Types) HasOption() bool {
	return t.set[OptionType] != 0
}

func (t Types) HasRef() bool {
	return t.set[RefType] != 0
}

func (t Types) HasTime() bool {
	return t.set[TimeType] != 0
}

func (t Types) HasTimeArray() bool {
	return t.set[TimeType]&tuArray == tuArray
}

func (t Types) HasTimeMap() bool {
	return t.set[TimeType]&tuMap == tuMap
}

type typeUsage int

const (
	tuSimple typeUsage = 1 << iota
	tuArray
	tuMap
)
