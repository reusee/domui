// +build ignore

package main

import (
	"github.com/reusee/domui"
	"syscall/js"
	"time"
)

var (
	Div = domui.Tag("div")
	T   = domui.Text
)

type (
	Spec = domui.Spec
)

// A string-typed state
type Greetings string

// define Greetings
func defGreetings() Greetings {
	return "Hello, world!"
}

// An UI element
type GreetingsElement Spec

// define GreetingsElement
func defGreetingsElement(
	// use Greetings
	greetings Greetings,
) GreetingsElement {
	return Div(T("%s", greetings))
}

// The root UI element
func defRootElement(
	// use GreetingsElement
	greetingsElem GreetingsElement,
) domui.RootElement {
	return Div(
		greetingsElem,
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		// provide definitions
		defRootElement,
		defGreetings,
		defGreetingsElement,
	)
	time.Sleep(time.Hour * 24 * 365 * 100)
}
