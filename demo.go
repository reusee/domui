// +build ignore

package main

import (
	"syscall/js"
	"time"

	"github.com/reusee/domui"
)

// Def is a type to hold all element and state definitions
type Def struct{}

// RootElement is the element to be rendered to RenderElement
func (_ Def) RootElement(
	// use CounterElement as depenency
	counter CounterElement,
) domui.RootElement {
	return Div(
		P(
			S("hello, world!"),
		),
		counter,
	)
}

// Num is an integer state
type Num int

// Num initial value
func (_ Def) Num() Num {
	return 0
}

// CounterElement is a button displaying and mutating Num
type CounterElement domui.Spec

func (_ Def) CounterElement(
	// use Num as depenency
	num Num,
	// use Update function
	update Update,
) CounterElement {
	return Button(

		// label
		S("%d", num),

		// style
		FontSize("%.2frem", float64(num)*0.5+1),
		Color("#09C"),

		// increase Num on click
		OnClick(func() {
			num++
			update(&num)
		}),
	)
}

// RenderElement is the HTMLElement to render on
func (_ Def) RenderElement() domui.RenderElement {
	return domui.RenderElement(
		js.Global().Get("document").Call("getElementById", "app"),
	)
}

func main() {
	domui.NewApp(new(Def))
	time.Sleep(time.Hour * 24 * 365 * 200)
}

// aliases
var (
	Div      = domui.Tag("div")
	P        = domui.Tag("p")
	S        = domui.S
	Button   = domui.Tag("button")
	OnClick  = domui.On("click")
	FontSize = domui.Style("font-size")
	Color    = domui.Style("color")
)

type (
	Update = domui.Update
)
