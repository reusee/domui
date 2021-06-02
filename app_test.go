package domui

import (
	"syscall/js"
	"testing"
)

func WithTestApp(t *testing.T, fn func(*App), defs ...any) {
	document := js.Global().Get("document")
	element := document.Call("createElement", "div")
	document.Get("body").Call("appendChild", element)
	app := NewApp(
		element,
		defs...,
	)
	fn(app)
}
