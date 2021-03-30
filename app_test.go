package domui

import (
	"testing"
)

func WithTestApp(t *testing.T, fn func(*App), initDecls ...any) {
	element := Document.Call("createElement", "div")
	Document.Get("body").Call("appendChild", element)
	app := NewApp(
		element,
		initDecls...,
	)
	fn(app)
}
