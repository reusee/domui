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
	Def  struct{}
	Spec = domui.Spec
)

// A string-typed state
type Greetings string

// define Greetings
func (_ Def) Greetings() Greetings {
	return "Hello, world!"
}

// An UI element
type GreetingsElement Spec

// define GreetingsElement
func (_ Def) GreetingsElement(
	// use Greetings
	greetings Greetings,
) GreetingsElement {
	return Div(T("%s", greetings))
}

// The root UI element
func (_ Def) RootElement(
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
		domui.Methods(new(Def))...,
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
	Def    struct{}
	Spec   = domui.Spec
	Update = domui.Update
)

type Greetings string

func (_ Def) Greetings() Greetings {
	return "Hello, world!"
}

type GreetingsElement Spec

func (_ Def) GreetingsElement(
	greetings Greetings,
) GreetingsElement {
	return Div(T("%s", greetings))
}

func (_ Def) RootElement(
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
		domui.Methods(new(Def))...,
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

type Def struct{}

func (_ Def) RootElement() domui.RootElement {
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
		domui.Methods(new(Def))...,
	)
	time.Sleep(time.Hour * 24 * 365 * 100)
}
```

### Parameterized element

To make reusable element, define it as a function.

```go
package main

import (
	"fmt"
	"github.com/reusee/domui"
	"reflect"
	"strings"
	"syscall/js"
	"time"
)

var (
	Div     = domui.Tag("div")
	T       = domui.Text
	OnClick = domui.On("click")
)

type (
	Def    struct{}
	any    = interface{}
	Spec   = domui.Spec
	Update = domui.Update
	Specs  = domui.Specs
)

// Greetings with name parameter
type Greetings func(name any) Spec

func (_ Def) Greetings(
	update Update,
) Greetings {
	return func(name any) Spec {
		return Specs{
			T("Hello, %s!", name),
			OnClick(func() {
				// when clicked, update the name argument to upper case
				// use reflect to support all string-typed arguments
				nameValue := reflect.New(reflect.TypeOf(name))
				nameValue.Elem().SetString(
					strings.ToUpper(fmt.Sprintf("%s", name)))
				update(nameValue.Interface())
			}),
		}
	}
}

type TheWorld string

func (_ Def) TheWorld() TheWorld {
	return "world"
}

type TheDomUI string

func (_ Def) TheDomUI() TheDomUI {
	return "DomUI"
}

func (_ Def) RootElement(
	greetings Greetings,
	world TheWorld,
	domUI TheDomUI,
) domui.RootElement {
	return Div(
		// use Greetings
		Div(
			greetings(world),
		),
		Div(
			greetings(domUI),
		),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		domui.Methods(new(Def))...,
	)
	time.Sleep(time.Hour * 24 * 365 * 100)
}
```

