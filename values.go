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

func newValues(root *node) *Values {
	return &Values{
		lock: &sync.RWMutex{},
		root: root,
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
	value, ok := v.root.getValueOrNodeOk(key)
	if !ok {
		return nil, false
	}
	if node, ok := value.(*node); ok {
		return newValues(node.clone()), true
	}
	return value, true
}

func (v *Values) Clone() *Values {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return newValues(v.root)
}

type node struct {
	value    interface{}
	children map[string]*node
}

func newNodeValue(value interface{}) *node {
	n := newNode()
	n.value = value
	n.children = nil
	return n
}

func newNode() *node {
	return &node{
		value:    nil,
		children: map[string]*node{},
	}
}

func (n *node) isEmpty() bool {
	return !n.isSet() && len(n.children) == 0
}

func (n *node) isSet() bool {
	return n.children == nil
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
	if n.isSet() {
		n.value = nil
		changed = true
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
	if n.isSet() {
		changed = value != n.value
	} else {
		changed = true
	}
	n.value = value
	n.children = nil
	return changed
}

func (n *node) setValues(values *Values) bool {
	if values.IsEmpty() {
		if n.isEmpty() {
			return false
		}
		n.value = nil
		n.children = nil
		return true
	}
	changed := false
	values.EachKeyValue(func(key Key, value interface{}) {
		changed = n.put(key, value) || changed
	})
	return changed
}

func (n *node) getValueOrNodeOk(key Key) (interface{}, bool) {
	if key.IsEmpty() {
		if n.isSet() {
			return n.value, true
		}
		return n, true
	}
	child, ok := n.children[key[0]]
	if !ok {
		return nil, false
	}
	return child.getValueOrNodeOk(key[1:])
}

func (n *node) clone() *node {
	return &node{
		value:    n.value,
		children: n.cloneChildren(),
	}
}

func (n *node) cloneChildren() map[string]*node {
	if n.children == nil {
		return nil
	}
	result := map[string]*node{}
	for key, child := range n.children {
		result[key] = child.clone()
	}
	return result
}

func (n *node) eachKeyValue(key Key, visitor func(key Key, value interface{})) {
	if n.isSet() {
		visitor(key, n.value)
		return
	}
	for keyPart, childNode := range n.children {
		childNode.eachKeyValue(append(NewKey(key...), keyPart), visitor)
	}
}