// +build ignore

package main

import (
	"github.com/reusee/domui"
	"syscall/js"
	"time"
)

var (
	Div        = domui.Tag("div")
	Link       = domui.Tag("a")
	Ahref      = domui.Attr("href")
	Sfont_size = domui.Style("font-size")
	T          = domui.Text
	ID         = domui.ID
	Class      = domui.Class
)

func defRootElement() domui.RootElement {
	return Div(
		Link(
			T("Hello, world!"),
			// id
			ID("link"),
			// class
			Class("link1", "link2"),
			// href attribute
			Ahref("http://github.com"),
			// font-size style
			Sfont_size("1.6rem"),
		),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		defRootElement,
	)
	time.Sleep(time.Hour * 24 * 365 * 100)
}
