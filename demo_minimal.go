// +build ignore

package main

import (
	"github.com/reusee/domui"
	"syscall/js"
	"time"
)

// ergonomic aliases
var (
	Div = domui.Tag("div")
	T   = domui.Text
)

// A type to hold all definitions
type Def struct{}

// define the root UI element
func (_ Def) RootElement() domui.RootElement {
	// same to <div>Hello, world!</div>
	return Div(T("Hello, world!"))
}

func main() {
	domui.NewApp(
		// render on <div id="app">
		js.Global().Get("document").Call("getElementById", "app"),
		// pass Def's methods as definitions
		domui.Methods(new(Def))...,
	)
	// prevent from exiting
	time.Sleep(time.Hour * 24 * 365 * 100)
}
