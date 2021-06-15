package domui

import (
	"fmt"
	"sync/atomic"
)

func Text(format string, args ...any) *Node {
	if len(args) == 0 {
		return &Node{
			serial: atomic.AddInt64(&nodeSerial, 1),
			Kind:   TextNode,
			Text:   format,
		}
	}
	return &Node{
		serial: atomic.AddInt64(&nodeSerial, 1),
		Kind:   TextNode,
		Text:   fmt.Sprintf(format, args...),
	}
}

var S = Text
