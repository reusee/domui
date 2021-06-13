// +build ignore

package main

import (
	"github.com/reusee/domui"
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
	Spec   = domui.Spec
	Update = domui.Update
)

type Greetings string

func (_ Def) Greetings() Greetings {
	return "Hello, world!"
}

type GreetingsElement Spec

func (_ Def) GreetingsElement(
	greetings Greetings,
) GreetingsElement {
	return Div(T("%s", greetings))
}

func (_ Def) RootElement(
	greetingsElem GreetingsElement,
	// use the Update function
	update Update,
) domui.RootElement {
	return Div(
		greetingsElem,

		// when clicked, do update
		OnClick(func() {
			// provide a new definition for Greetings
			update(func() Greetings {
				return "Hello, DomUI!"
			})
		}),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		domui.Methods(new(Def))...,
	)
	time.Sleep(time.Hour * 24 * 365 * 100)
}
