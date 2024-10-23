package domui

import (
	"reflect"
	"syscall/js"
	"time"

	"github.com/reusee/dscope"
)

type RootElement Spec

type App struct {
	wrapElement js.Value
	element     js.Value
	scope       dscope.Scope
	dirty       chan struct{}
	rootNode    *Node
	fns         chan any
}

func NewApp(
	renderElement js.Value,
	defs ...any,
) *App {

	app := &App{
		dirty: make(chan struct{}, 1),
		fns:   make(chan any),
	}

	app.scope = dscope.New(
		func() Update {
			return app.Update
		},
		func() *App {
			return app
		},
	)

	defs = append(defs, dscope.Methods(new(Def))...)
	app.scope = app.scope.Fork(defs...)

	dscope.Get[OnAppInit](app.scope)()

	parentElement := js.Value(renderElement)
	parentElement.Set("innerHTML", "")
	wrap := document.Call("createElement", "div")
	parentElement.Call("appendChild", wrap)
	element := document.Call("createElement", "div")
	wrap.Call("appendChild", element)
	app.wrapElement = wrap
	app.element = element

	go func() {
		for {
			select {

			case <-app.dirty:
				app.Render()

			case fn := <-app.fns:
				app.scope.Call(fn)

			}
		}
	}()

	app.Render()

	return app
}

type OnAppInit func()

var _ dscope.Reducer = OnAppInit(nil)

func (_ OnAppInit) IsReducer() {
}

func (_ Def) OnAppInit() OnAppInit {
	return func() {}
}

func (a *App) Update(decls ...any) {
	a.scope = a.scope.Fork(decls...)
	select {
	case a.dirty <- struct{}{}:
	default:
	}
}

var rootElementType = reflect.TypeOf((*RootElement)(nil)).Elem()

type SlowRenderThreshold time.Duration

func (_ Def) SlowRenderThreshold() SlowRenderThreshold {
	return SlowRenderThreshold(time.Millisecond) * 50
}

func (a *App) Render() {
	t0 := time.Now()
	var slowThreshold SlowRenderThreshold
	defer func() {
		e := time.Since(t0)
		if e > time.Duration(slowThreshold) {
			log("slow render in %v", time.Since(t0))
		}
	}()
	var rootElement RootElement
	a.scope.Assign(&slowThreshold, &rootElement)
	newNode := rootElement.(*Node)
	var err error
	a.element, err = patch(a.scope, newNode, a.element, a.rootNode)
	ce(err)
	a.rootNode = newNode
}

func (a *App) HTML() string {
	if a.element.InstanceOf(htmlElement) {
		return a.element.Get("outerHTML").String()
	}
	return a.element.Get("nodeValue").String()
}
