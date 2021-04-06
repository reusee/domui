// +build ignore

package main

import (
	"syscall/js"
	"time"

	"github.com/reusee/domui"
)

type Def struct{}

func (_ Def) RootElement() domui.RootElement {
	return domui.Tdiv(
		domui.S("hello, world!"),
	)
}

func (_ Def) RenderElement() domui.RenderElement {
	return domui.RenderElement(
		js.Global().Get("document").Call("getElementById", "app"),
	)
}

func main() {
	domui.NewApp(new(Def))
	time.Sleep(time.Hour * 24 * 365 * 200)
}
