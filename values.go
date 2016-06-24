package config

import "fmt"

type Values struct {
	root *node
}

func NewValues() *Values {
	return &Values{
		root: newNodeValue(nil),
	}
}

type node struct {
	value    interface{}
	set      bool
	children map[string]*node
}

func newNodeValue(value interface{}) *node {
	n := newNode()
	n.value, n.set = value, true
	return n
}

func newNodeChildren(children map[string]*node) *node {
	n := newNode()
	n.children = children
	return n
}

func newNode() *node {
	return &node{
		value:    nil,
		set:      false,
		children: nil,
	}
}

func (n *node) put(key Key, value interface{}) bool {
	if len(key) == 0 {
		return false
	}
	child, changed := n.getChild(key[0])
	if len(key) == 1 {
		return n.putLastKeyPart(key[0], child, value) || changed
	}
	remainingKey := key[1:]
	return child.put(remainingKey, value) || changed
}

func (n *node) putLastKeyPart(keyPart string, child *node, value interface{}) bool {
	_, ok := isValues(value)
	if ok {
		fmt.Println("************************************* need to implement putting values")
		return true
	}
	return child.setValue(value)
}

func isValues(value interface{}) (*Values, bool) {
	if values, ok := value.(*Values); ok {
		return values, ok
	}
	if values, ok := value.(Values); ok {
		return &values, ok
	}
	return nil, false
}

func (n *node) getChild(key string) (*node, bool) {
	changed := false
	if n.set {
		n.value = nil
		n.set = false
		changed = true
	}
	if n.children == nil {
		n.children = map[string]*node{}
	}
	child, ok := n.children[key]
	if !ok {
		n.children[key] = newNode()
		child = n.children[key]
		changed = true
	}
	return child, changed
}

func (n *node) setValue(value interface{}) bool {
	changed := false
	if n.set {
		changed = value != n.value
	} else {
		changed = true
	}
	n.value = value
	n.set = true
	n.children = nil
	return changed
}
