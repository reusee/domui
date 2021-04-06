package domui

import "fmt"

func Text(format string, args ...any) *Node {
	if len(args) == 0 {
		return &Node{
			Kind: TextNode,
			Text: format,
		}
	}
	return &Node{
		Kind: TextNode,
		Text: fmt.Sprintf(format, args...),
	}
}

var S = Text
