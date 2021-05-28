package domui

import "sync"

func NewNodeMap() (
	get func(key any, fn func() *Node) *Node,
) {
	var m sync.Map
	get = func(key any, fn func() *Node) *Node {
		if v, ok := m.Load(key); ok {
			return v.(*Node)
		}
		node := fn()
		m.Store(key, node)
		return node
	}
	return
}
