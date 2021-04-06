domui: frontend framework for Go

### Features

* pure go code compiled to wasm
* unified view declaration and state management
* reactive view and state transition

### Usage

#### Minimal demo

```go
// demo.go

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

Then open `http://localhost:46789/demo.html` in browser 

