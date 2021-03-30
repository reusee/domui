package domui

import "fmt"

func Text(text string) *Node {
	return &Node{
		Kind: TextNode,
		Text: text,
	}
}

func S(format string, args ...any) *Node {
	return Text(fmt.Sprintf(format, args...))
}
