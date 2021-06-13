domui: DOM UI framework for Go

## Features

* Pure go code compiled to wasm
* Dependent reactive system for both view and state management

## Prerequisites

* Go compiler 1.16 or newer
* If you are not familiar with compiling and running WebAssembly program, 
please read [the official wiki](https://github.com/golang/go/wiki/WebAssembly)

## Tutorial

### Minimal runnable program

```go
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
```

### The dependent system

The definition of RootElement can be refactored to multiple dependent components.

```go
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

// A UI element
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
```

### The reactive system

Definitions can be updated. 
All affected definitions will be re-calculated recursively till the RootElement.

```go
package main

import (
	"github.com/reusee/domui"
	"syscall/js"
	"time"
)

var (
	Div     = domui.Tag("div")
	T       = domui.Text
	OnClick = domui.On("click")
)

type (
	Spec   = domui.Spec
	Update = domui.Update
)

type Greetings string

func defGreetings() Greetings {
	return "Hello, world!"
}

type GreetingsElement Spec

func defGreetingsElement(
	greetings Greetings,
) GreetingsElement {
	return Div(T("%s", greetings))
}

func defRootElement(
	greetingsElem GreetingsElement,
	// use the Update function
	update Update,
) domui.RootElement {
	return Div(
		greetingsElem,

		// when clicked, do update
		OnClick(func() {
			// provide a new definition for Greetings
			update(func() Greetings {
				return "Hello, DomUI!"
			})
		}),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		defRootElement,
		defGreetings,
		defGreetingsElement,
	)
	time.Sleep(time.Hour * 24 * 365 * 100)
}
```

### DOM element specifications

The above programs demonstrated tag and event usages.
Attributes, styles, classes can also be specified.

```go
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
```

