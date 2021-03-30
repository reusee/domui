package domui

import (
	"fmt"
	"reflect"
	"strings"
)

type Spec interface {
	IsSpec()
}

// utils

func If(cond bool, specs ...Spec) Spec {
	if cond {
		return Specs(specs)
	}
	return nil
}

func Alt(cond bool, spec1 Spec, spec2 Spec) Spec {
	if cond {
		return spec1
	}
	return spec2
}

func For(slice any, fn any) Specs {
	sliceValue := reflect.ValueOf(slice)
	fnValue := reflect.ValueOf(fn)
	var specs Specs
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i)
		ret := fnValue.Call([]reflect.Value{elem})
		s := ret[0].Interface()
		if s == nil {
			continue
		}
		specs = append(specs, s.(Spec))
	}
	return specs
}

func Range(slice any, fn any) Specs {
	sliceValue := reflect.ValueOf(slice)
	fnValue := reflect.ValueOf(fn)
	var specs Specs
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i)
		ret := fnValue.Call([]reflect.Value{reflect.ValueOf(i), elem})
		s := ret[0].Interface()
		if s == nil {
			continue
		}
		specs = append(specs, s.(Spec))
	}
	return specs
}

// elements

type IDSpec struct {
	Value string
}

func (_ IDSpec) IsSpec() {}

func ID(id string) IDSpec {
	return IDSpec{
		Value: id,
	}
}

type StyleSpec struct {
	Value string
}

func (_ StyleSpec) IsSpec() {}

func Style(value string) StyleSpec {
	return StyleSpec{
		Value: strings.TrimSpace(value),
	}
}

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

type ClassesSpec struct {
	Classes map[string]bool
}

func (_ ClassesSpec) IsSpec() {}

func Classes(names ...string) ClassesSpec {
	m := make(map[string]bool)
	for _, name := range names {
		m[name] = true
	}
	return ClassesSpec{
		Classes: m,
	}
}

var Class = Classes

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

var Attr = Attrs

type Specs []Spec

func (_ Specs) IsSpec() {}

type FocusSpec struct{}

func (_ FocusSpec) IsSpec() {}

var Focus = FocusSpec{}

type Lazy func() Spec

func (_ Lazy) IsSpec() {}
