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

// define the root UI element
func defRootElement() domui.RootElement {
	return Div(T("Hello, world!"))
}

func main() {
	domui.NewApp(
		// render on <div id="app">
		js.Global().Get("document").Call("getElementById", "app"),
		// provide definitions
		defRootElement,
	)
	// prevent from exiting
	time.Sleep(time.Hour * 24 * 365 * 100)
}
