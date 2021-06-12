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
	Spec   = domui.Spec
	Update = domui.Update
)

type Greetings string

func defGreetings() Greetings {
	return "Hello, world!"
}

type GreetingsElement Spec

func defGreetingsElement(
	greetings Greetings,
) GreetingsElement {
	return Div(T("%s", greetings))
}

func defRootElement(
	greetingsElem GreetingsElement,
	// use the Update function
	update Update,
) domui.RootElement {
	return Div(
		greetingsElem,

		// when clicked, update Greetings
		OnClick(func() {
			greetings := Greetings("Hello, DomUI!")
			update(&greetings)
		}),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		defRootElement,
		defGreetings,
		defGreetingsElement,
	)
	time.Sleep(time.Hour * 24 * 365 * 100)
}
