# DomUI: A Declarative UI Framework for Go and WebAssembly

DomUI is a UI framework written purely in Go, designed to be compiled to WebAssembly (WASM). It enables building web interfaces using Go's strong typing and reflection capabilities, leveraging a reactive dependency system for managing both UI components and application state in a unified manner.

## Core Concepts

*   **Go + WebAssembly:** Write your entire frontend application in Go. Benefit from Go's tooling, static typing, and compile it to WASM to run in any modern browser or Electron.
*   **Unified Reactive System:** DomUI employs a dependency injection system (`dscope`) to automatically track relationships between state values and UI components. When a value changes, only the dependent parts of the UI are re-evaluated and updated.
*   **Declarative UI:** Define your UI structure using Go functions and types. Components are simply functions returning `Spec` types or specific state types.
*   **Efficient DOM Updates:** Uses a virtual DOM diffing and patching mechanism to minimize direct manipulations of the actual browser DOM, ensuring efficient updates.
*   **Simplified State Management:** The reactive dependency system inherently manages state. There's no need for separate state management libraries like Redux, Recoil, or MobX, or complex hooks like in React. State and UI components are part of the same dependency graph.

## Prerequisites

*   Go compiler version 1.16 or newer.
*   Familiarity with Go's WebAssembly compilation and execution model. Refer to the [official Go WebAssembly Wiki](https://github.com/golang/go/wiki/WebAssembly) if needed.
*   The standard `wasm_exec.js` file (provided in this repository or from your Go installation) is required to load and run the compiled WASM module in the browser.

## Table of Contents

*   [Tutorial](#tutorial)
    *   [Minimal Example](#minimal)
    *   [The Dependency System](#dependent)
    *   [The Reactive System](#reactive)
    *   [Defining DOM Elements](#dom)
    *   [Parameterized Components](#parameterized)
    *   [Event Handling](#events)
    *   [Conditional & Loop Rendering](#conditionals-loops)
    *   [Component Caching](#caching)
    *   [Initialization Hook](#init)
*   [Comparison with ReactJS](#reactjs)
*   [Running the Demo](#running-demo)

<a name="tutorial" />

## Tutorial

This tutorial guides you through the core features of DomUI.

<a name="minimal" />

### Minimal Example

This is the simplest runnable DomUI application.

```go
package main

import (
	"github.com/reusee/domui"
	"syscall/js"
	"time"
)

// Define convenient aliases for common DomUI functions
var (
	Div = domui.Tag("div") // Creates a <div> tag node
	T   = domui.Text     // Creates a text node
)

// Def is a struct used to group our definitions (component functions).
type Def struct{}

// RootElement defines the main UI component to be rendered.
// It returns a domui.RootElement, which is essentially a domui.Spec.
func (_ Def) RootElement() domui.RootElement {
	// Equivalent to <div>Hello, world!</div>
	return Div(
		T("Hello, world!"),
	)
}

func main() {
	// Initialize the DomUI application
	domui.NewApp(
		// The target DOM element where the UI will be rendered.
		// Here, it's <div id="app"></div> in your HTML.
		js.Global().Get("document").Call("getElementById", "app"),

		// Pass the methods of our Def struct as definitions.
		// DomUI's dependency system will discover and use them.
		domui.Methods(new(Def))...,
	)

	// Keep the Go program running (WASM needs this)
	select {} // A better way than sleeping indefinitely
}

```
*Save this as `main.go`, create an `index.html` with `<div id="app"></div>`, compile to WASM, and serve.*

<a name="dependent" />

### The Dependency System

DomUI uses a dependency injection system (`dscope`) to manage relationships between components and state. Define functions that return specific types, and other functions can accept these types as arguments to declare dependencies.

```go
package main

import (
	"fmt"
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
	Spec = domui.Spec // Alias for domui.Spec interface
)

// Define a custom type for our state
type Greetings string

// Define a function that provides the initial value for Greetings.
// DomUI's system will find this function based on its return type.
func (_ Def) ProvideGreetings() Greetings {
	return "Hello, dependent world!"
}

// Define a UI component type. It's an alias for domui.Spec.
type GreetingsElement Spec

// Define a function that creates the GreetingsElement.
// It declares a dependency on the Greetings type by accepting it as an argument.
// DomUI will automatically provide the value from ProvideGreetings.
func (_ Def) CreateGreetingsElement(
	greetings Greetings, // Dependency injection
) GreetingsElement {
	// Use the injected greetings value
	return Div(
		T(string(greetings)), // Cast Greetings to string for Text
	)
}

// Define the RootElement, which now depends on GreetingsElement.
func (_ Def) RootElement(
	// Declare dependency on GreetingsElement
	greetingsElem GreetingsElement,
) domui.RootElement {
	return Div(
		greetingsElem, // Use the injected element
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		// Provide all definition methods
		domui.Methods(new(Def))...,
	)
	select {}
}
```
*DomUI automatically wires `ProvideGreetings` -> `CreateGreetingsElement` -> `RootElement` based on the types.*

<a name="reactive" />

### The Reactive System

State and UI updates are handled reactively. When a definition is updated using the `Update` function, all dependent definitions are automatically re-evaluated, ultimately triggering a re-render of the affected parts of the `RootElement`.

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
	Button  = domui.Tag("button")
	OnClick = domui.On("click") // Event handler spec
)

type (
	Def    struct{}
	Spec   = domui.Spec
	Update = domui.Update // Type for the update function
)

// State type
type Message string

// Initial state definition
func (_ Def) InitialMessage() Message {
	return "Click me!"
}

// UI component type
type MessageElement Spec

// Component definition, depends on Message
func (_ Def) CreateMessageElement(msg Message) MessageElement {
	return Div(T(string(msg)))
}

// Root element definition
func (_ Def) RootElement(
	messageElem MessageElement,
	// Declare dependency on the Update function provided by DomUI
	update Update,
) domui.RootElement {

	count := 0 // Local counter for demonstration

	return Div(
		messageElem, // Display the current message

		Button(
			T("Update Message"),
			// Add a click event handler
			OnClick(func() {
				// When the button is clicked, call update
				count++
				newMessage := Message("Updated message! #" + time.Now().Format("15:04:05"))
				// Provide a new definition for the Message type.
				// This can be a function that returns the new value...
				update(func() Message {
					return newMessage
				})
				// ... or for simple types without dependencies, a pointer to the new value.
				// update(&newMessage) // Equivalent for this case
			}),
		),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		domui.Methods(new(Def))...,
	)
	select {}
}
```
*Clicking the button calls `update` with a new definition for `Message`. DomUI detects this change, re-runs `CreateMessageElement` (because it depends on `Message`), then re-runs `RootElement` (because it depends on `MessageElement`), and finally patches the DOM.*

<a name="dom" />

### Defining DOM Elements

DomUI provides functions to specify various aspects of DOM elements like tags, text content, attributes, styles, classes, and IDs.

*   **Tags:** `domui.Tag(name string) func(...Spec) *Node` (e.g., `Div`, `Button`, `Input`)
*   **Text:** `domui.Text(format string, args ...any) *Node` (e.g., `T("Hello")`)
*   **Attributes:** `domui.Attr(name string) func(value any) AttrSpec` (e.g., `Ahref("http://...")`) or `domui.Attrs(keyvals ...any)`
*   **Styles:** `domui.Style(name string) func(format string, args ...any) StyleSpec` (e.g., `SfontSize("1.2em")`) or `domui.Styles(keyvals ...any)`
*   **Classes:** `domui.Class(names ...string) ClassesSpec` (e.g., `Class("active", "highlight")`)
*   **ID:** `domui.ID(id string) IDSpec` (e.g., `ID("main-content")`)

```go
package main

import (
	"github.com/reusee/domui"
	"syscall/js"
	"time"
)

// Aliases for spec functions
var (
	Div        = domui.Tag("div")
	A          = domui.Tag("a") // Anchor tag
	T          = domui.Text
	ID         = domui.ID         // Set element ID
	Class      = domui.Class      // Add CSS classes
	Ahref      = domui.Attr("href") // Specific attribute helper for href
	SfontSize  = domui.Style("font-size") // Specific style helper for font-size
	Styles     = domui.Styles     // Set multiple styles
	Attrs      = domui.Attrs      // Set multiple attributes
)

type Def struct{}

func (_ Def) RootElement() domui.RootElement {
	return Div(
		ID("container"), // Set ID attribute
		Class("main", "content"), // Add 'main' and 'content' classes

		A( // Create an anchor element <a>
			T("Visit GitHub"), // Text content
			ID("github-link"),
			Class("external-link"),
			Ahref("https://github.com/reusee/domui"), // Set href attribute
			Styles( // Set multiple inline styles
				"color", "blue",
				"text-decoration", "none",
			),
			Attrs( // Set multiple attributes
				"target", "_blank",
				"rel", "noopener noreferrer",
			),
			SfontSize("1.1rem"), // Set font-size style
		),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		domui.Methods(new(Def))...,
	)
	select {}
}
```

<a name="parameterized" />

### Parameterized Components

To create reusable UI elements, define them as functions that accept parameters and return a `Spec`.

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
	Specs  = domui.Specs // A slice of Specs
)

// Define a component function type. It takes a name and returns a Spec.
type Greeter func(name any) Spec

// Definition function for the Greeter component.
func (_ Def) CreateGreeter(update Update) Greeter {
	// Return the actual component function (closure)
	return func(name any) Spec {
		// Return a Specs slice containing the element and its behavior
		return Specs{
			T("Hello, %s!", name), // Display the greeting
			OnClick(func() {
				// Update the name *passed by the caller* to uppercase.
				// Reflection is used here to modify the original variable
				// that was passed as the 'name' argument.
				// This works for basic types like string.
				// For complex state, manage it via the dependency system.
				nameType := reflect.TypeOf(name)
				if nameType.Kind() == reflect.String {
					upperName := strings.ToUpper(fmt.Sprintf("%s", name))
					// Create a pointer to the new value of the correct type
					nameValuePtr := reflect.New(nameType)
					nameValuePtr.Elem().SetString(upperName)
					// Update the dependency system with the new value for the original state variable
					update(nameValuePtr.Interface())
				}
			}),
		}
	}
}

// Define some string-typed states
type (
	Name1 string
	Name2 string
)

// Define initial values for the states
func (_ Def) ProvideNames() (Name1, Name2) {
	return "World", "DomUI"
}

// Root element uses the Greeter component multiple times
func (_ Def) RootElement(
	greeter Greeter, // Depend on the Greeter component function
	name1   Name1,   // Depend on the state Name1
	name2   Name2,   // Depend on the state Name2
) domui.RootElement {
	return Div(
		// Use the Greeter component with different state variables
		Div(greeter(name1)), // Pass Name1 state
		Div(greeter(name2)), // Pass Name2 state
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		domui.Methods(new(Def))...,
	)
	select {}
}
```
*The `Greeter` function acts as a reusable component. Clicking on "Hello, World!" will update the `Name1` state via `update`, triggering a re-render for that specific greeting.*

<a name="events" />

### Event Handling

Use `domui.On(eventName)(handlerFunc)` to attach event listeners. The handler function can optionally accept a `js.Value` argument to access the target DOM element.

```go
package main

import (
	"github.com/reusee/domui"
	"syscall/js"
	"time"
)

var (
	Div     = domui.Tag("div")
	Input   = domui.Tag("input")
	T       = domui.Text
	OnInput = domui.On("input") // Input event
	OnClick = domui.On("click") // Click event
	Avalue  = domui.Attr("value") // Value attribute
	Atype   = domui.Attr("type")  // Type attribute
)

type (
	Def    struct{}
	Spec   = domui.Spec
	Update = domui.Update
)

type UserInput string

func (_ Def) InitialInput() UserInput {
	return ""
}

func (_ Def) RootElement(
	update Update,
	userInput UserInput,
) domui.RootElement {
	return Div(
		T("Enter text: "),
		Input(
			Atype("text"),
			Avalue(string(userInput)), // Bind input value to state
			// Update UserInput state on every input event
			OnInput(func(elem js.Value) { // Handler accepts js.Value
				newValue := UserInput(elem.Get("value").String())
				update(&newValue) // Update state using pointer shortcut
			}),
		),
		Div(
			T("You entered: %s", userInput),
		),
		Div(
			OnClick(func(elem js.Value) { // Access element in handler
				domui.Log("Div clicked! Tag: %s", elem.Get("tagName").String())
			}),
			T("Click this div (check console)"),
		),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		domui.Methods(new(Def))...,
	)
	select {}
}
```

<a name="conditionals-loops" />

### Conditional & Loop Rendering

DomUI provides helpers for conditional rendering and rendering lists or slices.

*   `domui.If(condition bool, specs ...Spec) Spec`: Renders `specs` only if `condition` is true.
*   `domui.Alt(condition bool, specIfTrue Spec, specIfFalse Spec) Spec`: Renders `specIfTrue` or `specIfFalse`.
*   `domui.For(slice any, func(item T) Spec) Specs`: Renders a spec for each item in the slice.
*   `domui.Range(slice any, func(index int, item T) Spec) Specs`: Renders a spec for each item, providing both index and item.

```go
package main

import (
	"github.com/reusee/domui"
	"syscall/js"
	"time"
)

var (
	Div        = domui.Tag("div")
	P          = domui.Tag("p")
	Button     = domui.Tag("button")
	T          = domui.Text
	OnClick    = domui.On("click")
	SfontWeight = domui.Style("font-weight")
)

type (
	Def    struct{}
	Spec   = domui.Spec
	Update = domui.Update
	Specs  = domui.Specs
)

type ShowDetails bool
type Items []string

func (_ Def) InitialState() (ShowDetails, Items) {
	return false, []string{"Apple", "Banana", "Cherry"}
}

func (_ Def) RootElement(
	update Update,
	show ShowDetails,
	items Items,
) domui.RootElement {
	return Div(

		Button(
			OnClick(func() {
				show = !show
				update(&show)
			}),
			T(domui.Alt(bool(show), T("Hide Details"), T("Show Details")).(*domui.Node).Text), // Use Alt for button text
		),

		// Conditional rendering with If
		domui.If(bool(show),
			P(T("Showing secret details!")),
			P(SfontWeight("bold"), T("This is important.")),
		),

		// Loop rendering with Range
		P(T("Items:")),
		Div(
			domui.Range(items, func(i int, item string) Spec {
				return P(T("%d: %s", i+1, item))
			}),
		),

		// Loop rendering with For
		P(T("Items (again):")),
		Div(
			domui.For(items, func(item string) Spec {
				// Use Alt inside a loop
				return domui.Alt(item == "Banana",
					P(Class("highlight"), T("* %s *", item)), // If true
					P(T("- %s", item)),                     // If false
				)
			}),
		),
	)
}

func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		domui.Methods(new(Def))...,
	)
	select {}
}
```

<a name="caching" />

### Component Caching

For potentially expensive components or components that are frequently rendered with the same inputs, you can use `domui.NewSpecMap()` to cache the resulting `Spec`.

```go
package main

// ... imports and aliases ...
var (
	Div = domui.Tag("div")
	T   = domui.Text
)

type (
	Def  struct{}
	Spec = domui.Spec
)

// Define a potentially expensive component function type
type Article func(title string, content string) Spec

// Define the Article component provider, using caching
func (_ Def) CreateArticle() Article {
	// Create a cache specific to this component type
	cache := domui.NewSpecMap()

	// Return the component function
	return func(title string, content string) Spec {
		// Use the cache:
		// Provide a key (must be comparable) and a function to generate the Spec if not cached.
		return cache(
			// Key: Use input parameters
			[2]string{title, content},
			// Value generator function: Executed only if key is not in cache
			func() Spec {
				domui.Log("Generating Article Spec for title: %s", title) // Log cache miss
				// Simulate expensive computation or complex structure
				return Div(
					domui.Tag("h3")(T(title)),
					domui.Tag("p")(T(content)),
				)
			},
		)
	}
}

func (_ Def) RootElement(article Article) domui.RootElement {
	return Div(
		article("Cached Title 1", "Content 1..."),
		article("Cached Title 2", "Content 2..."),
		article("Cached Title 1", "Content 1..."), // This call will hit the cache
	)
}

// ... main function ...
func main() {
	domui.NewApp(
		js.Global().Get("document").Call("getElementById", "app"),
		domui.Methods(new(Def))...,
	)
	select {}
}
```
*The `NewSpecMap` acts like a memoization cache. The generator function is only called once for each unique key.*

<a name="running-demo" />

## Running the Demo

The repository includes a `demo.go` file and a helper script `run-demo.go` to build and serve it.

1.  **Ensure Go is installed.**
2.  **Navigate to the `domui` directory.**
3.  **Run the helper script:**
    ```bash
    go run run-demo.go
    ```
    *(If `run-demo.go` is marked `//go:build ignore`, you might need to remove that line or run it explicitly: `go run ./run-demo.go`)*
4.  **Open your browser:** The script will output a URL (usually `http://127.0.0.1:46789/demo.html`). Visit this URL to see the demo application.

The `run-demo.go` script performs the following steps:
*   Compiles `demo.go` (or another specified Go file) to `demo.wasm` with `GOOS=js` and `GOARCH=wasm`.
*   Starts a simple HTTP server.
*   Serves the required `wasm_exec.js`.
*   Serves a basic `demo.html` that loads and runs the compiled `demo.wasm`.

