package config

import (
	"fmt"
	"sync"
)

type Values struct {
	lock *sync.RWMutex
	root *node
}

func NewValues() *Values {
	return &Values{
		lock: &sync.RWMutex{},
		root: newNode(),
	}
}

func (v *Values) Put(key Key, value interface{}) bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.root.put(key, value)
}

func (v *Values) EachKeyValue(visitor func(key Key, value interface{}) bool) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	v.root.eachKeyValue(nil, visitor)
}

type node struct {
	value    interface{}
	set      bool
	children map[string]*node
}

func (n *node) String() string {
	return fmt.Sprintf("&%v", *n)
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
		return n.setValue(value)
	}
	child, changed := n.getChild(key[0])
	remainingKey := key[1:]
	return child.put(remainingKey, value) || changed
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
	if values, ok := value.(*Values); ok {
		// fmt.Println("setValue() called with *Values ***********************************")
		return n.setValues(values)
	}
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

func (n *node) setValues(values *Values) bool {
	changed := false
	n.eachKeyValue(nil, func(key Key, value interface{}) bool {
		changed = n.put(key, value) || changed
		return false
	})
	return changed
}

func (n *node) eachKeyValue(key Key, visitor func(key Key, value interface{}) bool) {
	if n.set {
		visitor(key, n.value)
		return
	}
	for keyPart, childNode := range n.children {
		childNode.eachKeyValue(append(append(Key(nil), key...), keyPart), visitor)
	}
}
