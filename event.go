package domui

import (
	"sync"
	"sync/atomic"
	"syscall/js"

	"github.com/reusee/dscope"
)

var elementID int32 = 42

type EventSpec struct {
	Event string
	Func  any
}

func (_ EventSpec) IsSpec() {}

func On(ev string, fn any) EventSpec {
	return EventSpec{
		Event: ev,
		Func:  fn,
	}
}

func makeEventFunc(event string) func(fn func()) EventSpec {
	return func(fn func()) EventSpec {
		return On(event, fn)
	}
}

var (
	eventRegistryLock sync.RWMutex
	eventRegistry     = make(map[int32]map[string][]EventSpec)
	eventHandlerSet   = make(map[string]bool)
)

var eventHandlerScope = dscope.New()

func setEventSpecs(wrap js.Value, element js.Value, specs map[string][]EventSpec) {

	idValue := element.Get("__element_id__")
	var id int32
	if idValue.IsUndefined() {
		id = atomic.AddInt32(&elementID, 1)
		element.Set("__element_id__", id)
	} else {
		id = int32(idValue.Int())
	}

	for event := range specs {
		if eventHandlerSet[event] {
			continue
		}
		wrap.Call(
			"addEventListener",
			event,
			js.FuncOf(
				func(this js.Value, args []js.Value) any {
					go func() {
						ev := args[0]
						typ := ev.Get("type").String()
						bubbles := ev.Get("bubbles").Bool()
						for node := ev.Get("target"); !node.IsNull() && !node.IsUndefined() && !node.Equal(wrap); node = node.Get("parentNode") {
							idValue := node.Get("__element_id__")
							if idValue.IsUndefined() {
								if !bubbles {
									break
								}
								continue
							}
							id := int32(idValue.Int())
							eventRegistryLock.RLock()
							var specs []EventSpec
							if evs, ok := eventRegistry[id]; ok {
								if ss, ok := evs[typ]; ok {
									specs = append(ss[:0:0], ss...)
								}
							}
							eventRegistryLock.RUnlock()
							for _, spec := range specs {
								eventHandlerScope.Sub(
									func() js.Value {
										return node
									},
								).Call(spec.Func)
							}
							if !bubbles {
								break
							}
						}
					}()
					return nil
				},
			),
			true,
		)
		eventHandlerSet[event] = true
	}

	eventRegistryLock.Lock()
	defer eventRegistryLock.Unlock()
	eventRegistry[id] = specs

}

func unsetEventSpecs(element js.Value) {
	idValue := element.Get("__element_id__")
	var id int32
	if idValue.IsUndefined() {
		return
	} else {
		id = int32(idValue.Int())
	}
	eventRegistryLock.Lock()
	defer eventRegistryLock.Unlock()
	eventRegistry[id] = nil
	childNodes := element.Get("childNodes")
	for i := childNodes.Length() - 1; i >= 0; i-- {
		unsetEventSpecs(childNodes.Index(i))
	}
}
