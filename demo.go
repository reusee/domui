// +build ignore

package main

import (
	"syscall/js"
	"time"

	"github.com/reusee/domui"
)

type Def struct{}

func (_ Def) RootElement(
	counter CounterElement,
) domui.RootElement {
	return domui.Tdiv(
		domui.S("hello, world!"),
		counter,
	)
}

type Num int

func (_ Def) Num() Num {
	return 0
}

type CounterElement domui.Spec

func (_ Def) CounterElement(
	num Num,
	update Update,
) CounterElement {
	return Tbutton(

		// style
		Sfont_size("2rem"),
		Scursor("pointer"),

		// text
		S("%d", num),

		// event
		Eclick(func() {
			num++
			update(&num)
		}),
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

var (
	S          = domui.S
	Tdiv       = domui.Tdiv
	Tbutton    = domui.Tbutton
	Eclick     = domui.Eclick
	Sfont_size = domui.Sfont_size
	Scursor    = domui.Scursor
)

type (
	Update = domui.Update
)
