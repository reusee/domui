package domui

import (
	"strings"
	"syscall/js"
)

func Patch(
	scope Scope,
	node *Node,
	lastElement js.Value,
	lastNode *Node,
) (
	element js.Value,
	err error,
) {
	defer he(&err)

	if lastElement.IsUndefined() {
		panic("bad last element")
	}

	if lastNode != nil && lastNode == node {
		// same node, same element
		return lastElement, nil
	}

	replace := func(lastNode *Node) (err error) {
		defer he(&err)
		// replace element with newly created one
		element, err = node.ToElement(scope)
		ce(err)
		parent := lastElement.Get("parentNode")
		parent.Call("insertBefore", element, lastElement)
		lastElement.Call("remove")
		unsetEventSpecs(lastElement)
		return nil
	}

	if lastNode == nil {
		// not patchable
		ce(replace(lastNode))
		return
	}

	if node.Kind != lastNode.Kind {
		// not patchable
		ce(replace(lastNode))
		return
	}

	switch node.Kind {

	case TagNode:
		if node.Text != lastNode.Text {
			// not patchable
			ce(replace(lastNode))
			return
		}

	case TextNode:
		element = lastElement
		if node.Text != lastNode.Text {
			element.Set("data", node.Text)
		}
		return

	}

	// patchable
	element = lastElement

	//TODO match childNodes and lastChildNodes
	// children
	// child nodes may contains specs, handle first
	elementChildren := element.Get("childNodes")
	childNodes := node.childNodes
	lastChildNodes := lastNode.childNodes
	hasFocus := false
	for node := document.Get("activeElement"); !node.IsNull() && !node.IsUndefined() && !node.Equal(body); node = node.Get("parentNode") {
		if node.Equal(element) {
			hasFocus = true
			break
		}
	}
	for i, childNode := range childNodes {
		if i < len(lastChildNodes) {

			childElement := elementChildren.Index(i)
			hasScrollBar := false
			if childElement.InstanceOf(htmlElement) {
				hasScrollBar = hasScrollBar ||
					childElement.Get("scrollWidth").Int() > childElement.Get("clientWidth").Int()
				hasScrollBar = hasScrollBar ||
					childElement.Get("scrollHeight").Int() > childElement.Get("clientHeight").Int()
			}

			if !hasFocus && !hasScrollBar &&
				len(lastChildNodes) < len(childNodes) {
				// insert
				childElement, err := childNode.ToElement(scope)
				ce(err)
				element.Call(
					"insertBefore",
					childElement,
					elementChildren.Index(i),
				)
				lastChildNodes = append(
					lastChildNodes[:i],
					append([]*Node{nil}, lastChildNodes[i:]...)...,
				) // insert placeholder

			} else {
				// replace
				_, err := Patch(
					scope,
					childNode,
					childElement,
					lastChildNodes[i],
				)
				ce(err)
			}

		} else {
			// append
			childElement, err := childNode.ToElement(scope)
			ce(err)
			element.Call("appendChild", childElement)
		}

	}
	if n := len(lastChildNodes) - len(childNodes); n > 0 {
		for i := 0; i < n; i++ {
			lastChild := element.Get("lastChild")
			lastChild.Call("remove")
			unsetEventSpecs(lastChild)
		}
	}

	// id
	if node.ID != lastNode.ID {
		if node.ID == "" {
			element.Call("removeAttribute", "id")
			element.Delete("id")
		} else {
			element.Set("id", node.ID)
			element.Call("setAttribute", "id", node.ID)
		}
	}

	// style
	if node.Style != lastNode.Style {
		element.Set("style", node.Style)
	}

	// styles
	style := element.Get("style")
	// must do removing before adding, since different attributes may affect the same style
	// for example, adding `padding: 1px` then removing `padding-bottom` results to `1px 1px 0 1px` wrongly
	for _, item := range lastNode.Styles {
		if node.Styles != nil {
			if _, ok := node.Styles.Get(item.Key); !ok {
				style.Set(item.Key, nil)
			}
		} else {
			style.Set(item.Key, nil)
		}
	}
	for _, item := range node.Styles {
		if lastNode.Styles != nil {
			if v, ok := lastNode.Styles.Get(item.Key); !ok ||
				v != item.Value ||
				// always set animation properties
				strings.HasPrefix(item.Key, "animation") {
				key := item.Key
				if strings.HasSuffix(key, "|reset") {
					// reset
					key = strings.TrimSuffix(key, "|reset")
					style.Set(key, nil)
					element.Get("offsetHeight") // trigger reflow
				}
				style.Set(key, item.Value)
			}
		} else {
			style.Set(item.Key, item.Value)
		}
	}

	if node.Style == "" && len(node.Styles) == 0 {
		element.Call("removeAttribute", "style")
		element.Delete("style")
	}

	// classes
	list := element.Get("classList")
	if len(node.Classes) > 0 {
		for _, item := range node.Classes {
			if lastNode.Classes != nil {
				if _, ok := lastNode.Classes.Get(item.Key); !ok {
					list.Call("add", item.Key)
				}
			} else {
				list.Call("add", item.Key)
			}
		}
		for _, item := range lastNode.Classes {
			if node.Classes != nil {
				if _, ok := node.Classes.Get(item.Key); !ok {
					list.Call("remove", item.Key)
				}
			} else {
				list.Call("remove", item.Key)
			}
		}
	} else {
		element.Call("removeAttribute", "class")
		element.Delete("class")
	}

	// attrs
	for _, item := range node.Attributes {
		if lastNode.Attributes != nil {
			if v, ok := lastNode.Attributes.Get(item.Key); !ok || v != item.Value {
				element.Call("setAttribute", item.Key, item.Value)
				element.Set(item.Key, item.Value)
			}
		} else {
			element.Call("setAttribute", item.Key, item.Value)
			element.Set(item.Key, item.Value)
		}
	}
	for _, item := range lastNode.Attributes {
		if node.Attributes != nil {
			if _, ok := node.Attributes.Get(item.Key); !ok {
				element.Call("removeAttribute", item.Key)
				element.Delete(item.Key)
			}
		} else {
			element.Call("removeAttribute", item.Key)
			element.Delete(item.Key)
		}
	}

	// events
	if len(node.Events) > 0 {
		var app *App
		scope.Assign(&app)
		setEventSpecs(app.wrapElement, element, node.Events)
	} else {
		unsetEventSpecs(element)
	}

	// focus
	if node.Focus {
		element.Call("focus")
	}

	return
}
