package domui

type AttrsSpec struct {
	Attrs map[string]any
}

func (_ AttrsSpec) IsSpec() {}

func Attrs(args ...any) AttrsSpec {
	m := make(map[string]any)
	for i := 0; i < len(args); i += 2 {
		k := args[i].(string)
		v := args[i+1]
		m[k] = v
	}
	return AttrsSpec{
		Attrs: m,
	}
}

type AttrSpec struct {
	Name  string
	Value any
}

func (_ AttrSpec) IsSpec() {}

func Attr(name string) func(value any) AttrSpec {
	return func(value any) AttrSpec {
		return AttrSpec{
			Name:  name,
			Value: value,
		}
	}
}
