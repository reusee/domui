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
	any    = interface{}
	Spec   = domui.Spec
	Update = domui.Update
	Specs  = domui.Specs
)

// Greetings with name parameter
type Greetings func(name any) Spec

func defGreetings(
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

func defTheWorld() TheWorld {
	return "world"
}

type TheDomUI string

func defTheDomUI() TheDomUI {
	return "DomUI"
}

func defRootElement(
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
		defRootElement,
		defGreetings,
		defTheWorld,
		defTheDomUI,
	)
	time.Sleep(time.Hour * 24 * 365 * 100)
}
