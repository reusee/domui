// +build ignore

package main

import (
	"fmt"
	"github.com/reusee/domui"
	"reflect"
	"strings"
	"syscall/js"
	"time"
)

var (
	Div     = domui.Tag("div")
	T       = domui.Text
	OnClick = domui.On("click")
)

type (
	Def    struct{}
	any    = interface{}
	Spec   = domui.Spec
	Update = domui.Update
	Specs  = domui.Specs
)

// Greetings with name parameter
type Greetings func(name any) Spec

func (_ Def) Greetings(
	update Update,
) Greetings {
	return func(name any) Spec {
		return Specs{
			T("Hello, %s!", name),
			OnClick(func() {
				// when clicked, update the name argument to upper case
				// use reflect to support all string-typed arguments
				nameValue := reflect.New(reflect.TypeOf(name))
				nameValue.Elem().SetString(
					strings.ToUpper(fmt.Sprintf("%s", name)))
				update(nameValue.Interface())
			}),
		}
	}
}

type TheWorld string

func (_ Def) TheWorld() TheWorld {
	return "world"
}

type TheDomUI string

func (_ Def) TheDomUI() TheDomUI {
	return "DomUI"
}

func (_ Def) RootElement(
	greetings Greetings,
	world TheWorld,
	domUI TheDomUI,
) domui.RootElement {
	return Div(
		// use Greetings
		Div(
			greetings(world),
		),
		Div(
			greetings(domUI),
		),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		domui.Methods(new(Def))...,
	)
	time.Sleep(time.Hour * 24 * 365 * 100)
}
