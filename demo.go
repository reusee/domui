//go:build ignore

package main

import (
	"syscall/js"
	"time"

	"github.com/reusee/domui"
)

// RootElement is the element to be rendered to RenderElement
func RootElementDecl(
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
func NumDecl() Num {
	return 0
}

// CounterElement is a button displaying and mutating Num
type CounterElement domui.Spec

func CounterElementDecl(
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

func main() {
	domui.NewApp(
		// render element
		js.Global().Get("document").Call("getElementById", "app"),
		// list all declarations
		RootElementDecl,
		NumDecl,
		CounterElementDecl,
	)
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
