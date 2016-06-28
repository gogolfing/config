package config

import "sync"

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

func (v *Values) Merge(key Key, other *Values) bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	changed := false
	other.EachKeyValue(func(otherKey Key, value interface{}) {
		actualKey := key.Append(otherKey)
		changed = v.put(actualKey, value) || changed
	})
	return changed
}

func (v *Values) EachKeyValue(visitor func(key Key, value interface{})) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	v.root.eachKeyValue(nil, visitor)
}

func (v *Values) Put(key Key, value interface{}) bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.put(key, value)
}

func (v *Values) put(key Key, value interface{}) bool {
	return v.root.put(key, value)
}

func (v *Values) IsEmpty() bool {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return v.root.isEmpty()
}

func (v *Values) Get(key Key) interface{} {
	value, _ := v.GetOk(key)
	return value
}

func (v *Values) GetOk(key Key) (interface{}, bool) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return v.root.getOk(key)
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

func (n *node) isEmpty() bool {
	return !n.set && len(n.children) == 0
}

func (n *node) put(key Key, value interface{}) bool {
	if key.IsEmpty() {
		return n.setValue(value)
	}
	child, changed := n.getChild(key[0])
	remainingKey := key[1:]
	return child.put(remainingKey, value) || changed
}

func (n *node) getChild(keyPart string) (*node, bool) {
	changed := false
	if n.set {
		n.value = nil
		n.set = false
		changed = true
	}
	if n.children == nil {
		n.children = map[string]*node{}
	}
	child, ok := n.children[keyPart]
	if !ok {
		n.children[keyPart] = newNode()
		child = n.children[keyPart]
		changed = true
	}
	return child, changed
}

func (n *node) setValue(value interface{}) bool {
	if values, ok := value.(*Values); ok {
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
	if values.IsEmpty() {
		if n.isEmpty() {
			return false
		}
		n.value = nil
		n.set = false
		n.children = nil
		return true
	}
	changed := false
	values.EachKeyValue(func(key Key, value interface{}) {
		changed = n.put(key, value) || changed
	})
	return changed
}

func (n *node) getOk(key Key) (interface{}, bool) {
	if key.IsEmpty() {
		if n.set {
			return n.value, true
		}
		return nil, false
	}
	if n.set {
		return nil, false
	}
	child, ok := n.children[key[0]]
	if !ok {
		return nil, false
	}
	return child.getOk(key[1:])
}

func (n *node) eachKeyValue(key Key, visitor func(key Key, value interface{})) {
	if n.set {
		visitor(key, n.value)
		return
	}
	for keyPart, childNode := range n.children {
		childNode.eachKeyValue(append(NewKey(key...), keyPart), visitor)
	}
}
