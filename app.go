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
	getScope    dscope.GetScope
	mutate      dscope.Mutate
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

	scope := dscope.NewMutable(
		func() Update {
			return app.Update
		},
		func() *App {
			return app
		},
	)
	scope.Assign(&app.getScope, &app.mutate)

	defs = append(defs, dscope.Methods(new(Def))...)
	app.mutate(defs...)

	var onInit OnAppInit
	app.getScope().Assign(&onInit)

	onInit()

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
				app.getScope().Call(fn)

			}
		}
	}()

	app.Render()

	return app
}

type OnAppInit func()

var _ dscope.Reducer = OnAppInit(nil)

func (_ OnAppInit) Reduce(_ Scope, vs []reflect.Value) reflect.Value {
	return dscope.Reduce(vs)
}

func (_ Def) OnAppInit() OnAppInit {
	return func() {}
}

func (a *App) Update(decls ...any) Scope {
	scope := a.mutate(decls...)
	select {
	case a.dirty <- struct{}{}:
	default:
	}
	return scope
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
	scope := a.getScope()
	var rootElement RootElement
	scope.Assign(&slowThreshold, &rootElement)
	newNode := rootElement.(*Node)
	var err error
	a.element, err = patch(scope, newNode, a.element, a.rootNode)
	ce(err)
	a.rootNode = newNode
}

func (a *App) HTML() string {
	if a.element.InstanceOf(htmlElement) {
		return a.element.Get("outerHTML").String()
	}
	return a.element.Get("nodeValue").String()
}
