package domui

import "sync"

func NewSpecMap() (
	get func(key any, fn func() Spec) Spec,
) {
	var m sync.Map
	get = func(key any, fn func() Spec) Spec {
		if v, ok := m.Load(key); ok {
			return v.(Spec)
		}
		spec := fn()
		m.Store(key, spec)
		return spec
	}
	return
}
