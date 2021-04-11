domui: DOM UI framework for Go

## Features

* pure go code compiled to wasm
* unified view and state in dependent declarations

## Demo

```go
// demo.go

package main

import (
	"syscall/js"
	"time"

	"github.com/reusee/dscope"
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
	domui.NewApp(dscope.Methods(new(Def))...)
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
```

build wasm file:

```
env GOOS=js GOARCH=wasm go build -o demo.wasm demo.go
```

entry html:

```html
<!-- demo.html -->
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <script src="wasm_exec.js"></script>
  </head>
  <body>
    <div id="app"></div>
    <script>

      (async function exec() {
        const go = new Go();
        const result = await WebAssembly.instantiateStreaming(fetch("demo.wasm"), go.importObject);
        await go.run(result.instance);
      })()

    </script>
  </body>
</html>

```

run a simple http file server:

```go

package main

import (
	"net/http"
	"os"
)

func main() {
	dirFS := os.DirFS(".")
	http.Handle("/", http.FileServer(http.FS(dirFS)))
	http.ListenAndServe(":46789", nil)
}

```

Open `http://localhost:46789/demo.html` in browser to check the result

