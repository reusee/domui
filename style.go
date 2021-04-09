package domui

import (
	"fmt"
)

type StyleString string

func (_ StyleString) IsSpec() {}

type StylesSpec struct {
	Styles map[string]string
}

func (_ StylesSpec) IsSpec() {}

func Styles(args ...any) StylesSpec {
	m := make(map[string]string)
	for i := 0; i < len(args); i += 2 {
		k := args[i].(string)
		v := fmt.Sprintf("%v", args[i+1])
		m[k] = v
	}
	return StylesSpec{
		Styles: m,
	}
}

var CSS = Styles

type StyleSpec struct {
	Name  string
	Value string
}

func (_ StyleSpec) IsSpec() {}

func Style(name string) func(value any) StyleSpec {
	return func(value any) StyleSpec {
		return StyleSpec{
			Name:  name,
			Value: fmt.Sprintf("%v", value),
		}
	}
}
