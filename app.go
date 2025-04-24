package domui

import (
	"reflect"
	"sync"
	"syscall/js"
	"time"

	"github.com/reusee/dscope"
)

type RootElement Spec

type App struct {
	wrapElement js.Value
	element     js.Value
	dirty       chan struct{}
	rootNode    *Node
	scope       dscope.Scope
	scopeLock   sync.Mutex
}

func NewApp(
	renderElement js.Value,
	defs ...any,
) *App {

	app := &App{
		dirty: make(chan struct{}, 1),
	}

	defs = append(
		defs,
		func() Update {
			return app.Update
		},
		func() *App {
			return app
		},
	)

	app.scope = dscope.New(defs...)

	app.scope = app.scope.Fork(
		dscope.Methods(new(Def))...,
	)

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

			}
		}
	}()

	app.Render()

	return app
}

func (a *App) Update(defs ...any) {
	a.scopeLock.Lock()
	defer a.scopeLock.Unlock()
	a.scope = a.scope.Fork(defs...)
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

	a.scopeLock.Lock()
	defer a.scopeLock.Unlock()

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
