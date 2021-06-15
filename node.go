package domui

import (
	"fmt"
	"reflect"
	"sort"
	"syscall/js"
)

type NodeKind uint8

const (
	TagNode NodeKind = iota
	TextNode
)

var (
	nodeSerial int64
)

type Node struct {
	serial           int64
	Kind             NodeKind
	Text             string
	ID               string
	Style            string
	Styles           SortedMap // string: string
	Classes          SortedMap // string: struct{}
	Attributes       SortedMap // string: any
	Events           map[string][]EventSpec
	childNodes       []*Node
	childNodesSorted []*Node
	Focus            bool
	args             []reflect.Value
}

func (_ *Node) IsSpec() {}

func (n *Node) ToElement(scope Scope) (_ js.Value, err error) {
	defer he(&err)

	switch n.Kind {

	case TagNode:
		element := document.Call(
			"createElement",
			n.Text,
		)

		if len(n.childNodes) > 0 {
			fragment := document.Call("createDocumentFragment")
			for _, childNode := range n.childNodes {
				childElement, err := childNode.ToElement(scope)
				ce(err)
				fragment.Call("append", childElement)
			}
			element.Call("appendChild", fragment)
		}

		if n.ID != "" {
			element.Set("id", n.ID)
		}

		if n.Style != "" {
			element.Set("style", n.Style)
		}

		if len(n.Styles) > 0 {
			style := element.Get("style")
			for _, item := range n.Styles {
				style.Set(item.Key, item.Value)
			}
		}

		if len(n.Classes) > 0 {
			list := element.Get("classList")
			for _, item := range n.Classes {
				list.Call("add", item.Key)
			}
		}

		if len(n.Attributes) > 0 {
			for _, item := range n.Attributes {
				element.Call("setAttribute", item.Key, item.Value)
				element.Set(item.Key, item.Value)
			}
		}

		// events
		if len(n.Events) > 0 {
			var app *App
			scope.Assign(&app)
			setEventSpecs(app.wrapElement, element, n.Events)
		}

		if n.Focus {
			element.Call("focus")
		}

		return element, nil

	case TextNode:
		element := document.Call(
			"createTextNode",
			n.Text,
		)
		return element, nil

	}

	panic("bad kind")
}

func (node *Node) ApplySpec(spec Spec) {
	if spec == nil {
		return
	}

	switch spec := spec.(type) {

	case IDSpec:
		node.ID = spec.Value

	case StyleString:
		node.Style = string(spec)

	case StyleSpec:
		node.Styles.Set(spec.Name, spec.Value)

	case StylesSpec:
		for k, v := range spec.Styles {
			node.Styles.Set(k, v)
		}

	case ClassesSpec:
		for k := range spec.Classes {
			node.Classes.Set(k, struct{}{})
		}

	case AttrsSpec:
		for k, v := range spec.Attrs {
			node.Attributes.Set(k, v)
		}

	case AttrSpec:
		node.Attributes.Set(spec.Name, spec.Value)

	case EventSpec:
		if node.Events == nil {
			node.Events = make(map[string][]EventSpec)
		}
		node.Events[spec.Event] = append(
			node.Events[spec.Event],
			spec,
		)

	case *Node:
		node.appendChild(spec)

	case Specs:
		for _, s := range spec {
			node.ApplySpec(s)
		}

	case FocusSpec:
		node.Focus = true

	case Lazy:
		s := spec()
		node.ApplySpec(s)

	default:
		panic(fmt.Errorf("unknown spec: %#v", spec))

	}
}

func (n *Node) appendChild(child *Node) {
	n.insertChild(len(n.childNodes), child)
}

func (n *Node) insertChild(pos int, child *Node) {
	if l := len(n.childNodes); pos < l {
		n.childNodes = append(n.childNodes[:pos],
			append([]*Node{child}, n.childNodes[pos:]...)...)
	} else if pos == l {
		n.childNodes = append(n.childNodes, child)
	} else {
		panic("bad position")
	}
	i := sort.Search(len(n.childNodesSorted), func(i int) bool {
		return n.childNodesSorted[i].serial > child.serial
	})
	if i < len(n.childNodesSorted) {
		newSlice := make([]*Node, 0, len(n.childNodesSorted)+1)
		newSlice = append(newSlice, n.childNodesSorted[:i]...)
		newSlice = append(newSlice, child)
		newSlice = append(newSlice, n.childNodesSorted[i:]...)
		n.childNodesSorted = newSlice
	} else {
		n.childNodesSorted = append(n.childNodesSorted, child)
	}
}

func (n *Node) popChild() {
	n.deleteChild(len(n.childNodes) - 1)
}

func (n *Node) deleteChild(pos int) {
	var deleted *Node
	if l := len(n.childNodes); pos >= l {
		panic("bad position")
	} else if pos == l-1 {
		deleted = n.childNodes[l-1]
		n.childNodes = n.childNodes[:l-1]
	} else {
		deleted = n.childNodes[pos]
		n.childNodes = append(
			n.childNodes[:pos],
			n.childNodes[:pos+1]...,
		)
	}
	i := sort.Search(len(n.childNodesSorted), func(i int) bool {
		return n.childNodesSorted[i].serial >= deleted.serial
	})
	if i >= len(n.childNodesSorted) {
		panic("impossible")
	}
	if n.childNodesSorted[i].serial != deleted.serial {
		panic("impossible")
	}
	n.childNodesSorted = append(
		n.childNodesSorted[:i],
		n.childNodesSorted[i+1:]...,
	)
}

func (n *Node) replaceChild(pos int, child *Node) {
	deleted := n.childNodes[pos]
	n.childNodes[pos] = child
	i := sort.Search(len(n.childNodesSorted), func(i int) bool {
		return n.childNodesSorted[i].serial >= deleted.serial
	})
	if i >= len(n.childNodesSorted) {
		panic("impossible")
	}
	if n.childNodesSorted[i].serial != deleted.serial {
		panic("impossible")
	}
	n.childNodesSorted = append(
		n.childNodesSorted[:i],
		n.childNodesSorted[i+1:]...,
	)
	i = sort.Search(len(n.childNodesSorted), func(i int) bool {
		return n.childNodesSorted[i].serial > child.serial
	})
	if i < len(n.childNodesSorted) {
		newSlice := make([]*Node, 0, len(n.childNodesSorted)+1)
		newSlice = append(newSlice, n.childNodesSorted[:i]...)
		newSlice = append(newSlice, child)
		newSlice = append(newSlice, n.childNodesSorted[i:]...)
		n.childNodesSorted = newSlice
	} else {
		n.childNodesSorted = append(n.childNodesSorted, child)
	}
}
