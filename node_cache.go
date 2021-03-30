package domui

import "sync"

func NewNodeCache() (
	get func(key any, fn func() *Node) *Node,
) {
	var cache sync.Map
	get = func(key any, fn func() *Node) *Node {
		if v, ok := cache.Load(key); ok {
			return v.(*Node)
		}
		node := fn()
		cache.Store(key, node)
		return node
	}
	return
}
