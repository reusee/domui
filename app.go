package domui

import (
	"reflect"
	"sync"
	"syscall/js"
	"time"

	"github.com/reusee/dscope"
)

type RootElement Spec

type RenderElement js.Value

type App struct {
	wrapElement js.Value
	element     js.Value
	scopeLock   sync.RWMutex
	scope       Scope
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

	scope := dscope.New(
		func() Update {
			return app.Update
		},
		func() ScopedCall {
			return app.ScopedCall
		},
		func() *App {
			return app
		},
	)
	app.scope = scope

	defs := dscope.Methods(new(Def))
	for _, obj := range defObjects {
		defs = append(defs, dscope.Methods(obj))
	}
	app.scope = app.scope.Sub(defs...)

	var onInit OnAppInit
	var renderElement RenderElement
	app.scope.Assign(&onInit, &renderElement)

	onInit()

	parentElement := js.Value(renderElement)
	parentElement.Set("innerHTML", "")
	wrap := Document.Call("createElement", "div")
	parentElement.Call("appendChild", wrap)
	element := Document.Call("createElement", "div")
	wrap.Call("appendChild", element)
	app.wrapElement = wrap
	app.element = element

	go func() {
		for {
			select {

			case <-app.dirty:
				app.Render()

			case fn := <-app.fns:
				app.ScopedCall(fn)

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
	a.scopeLock.Lock()
	defer a.scopeLock.Unlock()
	if len(decls) == 0 {
		return a.scope
	}
	a.scope = a.scope.Sub(decls...)
	select {
	case a.dirty <- struct{}{}:
	default:
	}
	return a.scope
}

func (a *App) ScopedCall(fn any) {
	a.scopeLock.RLock()
	s := a.scope
	a.scopeLock.RUnlock()
	s.Call(fn)
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
	a.ScopedCall(func(scope Scope) {
		v, err := scope.Get(rootElementType)
		ce(err)
		newNode := v.Interface().(*Node)
		a.element, err = Patch(scope, newNode, a.element, a.rootNode)
		ce(err)
		a.rootNode = newNode
	})
}

func (a *App) HTML() string {
	if a.element.InstanceOf(HTMLElement) {
		return a.element.Get("outerHTML").String()
	}
	return a.element.Get("nodeValue").String()
}
