package domui

import (
	"reflect"
	"syscall/js"
	"time"

	"github.com/reusee/dscope"
)

type RootElement Spec

type RenderElement js.Value

type App struct {
	wrapElement js.Value
	element     js.Value
	getScope    dscope.Get
	derive      dscope.Derive
	dirty       chan struct{}
	rootNode    *Node
	fns         chan any
}

func NewApp(
	defObjects ...any,
) *App {

	app := &App{
		dirty: make(chan struct{}, 1),
		fns:   make(chan any),
	}

	scope := dscope.NewDeriving(
		func() Update {
			return app.Update
		},
		func() *App {
			return app
		},
	)
	scope.Assign(&app.getScope, &app.derive)

	defs := dscope.Methods(new(Def))
	for _, obj := range defObjects {
		defs = append(defs, dscope.Methods(obj)...)
	}
	app.derive(defs...)

	var onInit OnAppInit
	var renderElement RenderElement
	app.getScope().Assign(&onInit, &renderElement)

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
	scope := a.derive(decls...)
	select {
	case a.dirty <- struct{}{}:
	default:
	}
	return scope
}

var rootElementType = reflect.TypeOf((*RootElement)(nil)).Elem()

func (a *App) Render() {
	t0 := time.Now()
	defer func() {
		e := time.Since(t0)
		if e > time.Millisecond*100 {
			log("slow render in %v", time.Since(t0))
		}
	}()
	scope := a.getScope()
	v, err := scope.Get(rootElementType)
	ce(err)
	newNode := v.Interface().(*Node)
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
