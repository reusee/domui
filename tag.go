package domui

import (
	"syscall/js"
)

func Tag(name string) func(specs ...Spec) *Node {
	return func(specs ...Spec) *Node {
		node := &Node{
			Kind: TagNode,
			Text: name,
		}

		for _, spec := range specs {
			node.ApplySpec(spec)
		}

		// bind input state to node
		if name == "input" {
			v, ok := node.Attributes.Get("type")
			var typ string
			if ok {
				typ = v.(string)
			} else {
				typ = "text"
			}
			switch typ {

			case "checkbox", "radio":
				node.ApplySpec(On("input")(func(elem js.Value) {
					node.Attributes.Set("checked", elem.Get("checked").String())
				}))

			case "", "color",
				"date",
				"file",
				"datetime-local",
				"email",
				"month",
				"number",
				"password",
				"range",
				"search",
				"tel",
				"text",
				"time",
				"url",
				"week":
				node.ApplySpec(On("input")(func(elem js.Value) {
					node.Attributes.Set("value", elem.Get("value").String())
				}))

			}
		}

		return node
	}
}
