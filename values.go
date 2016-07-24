package config

import "sync"

//Values provides storage of arbitrary interface{} values referenced by type Key.
//The zero value for Values is not in a valid state, thus Values should be
//created with NewValues().
//Values is safe for use by multiple goroutines.
type Values struct {
	lock *sync.RWMutex
	root *node
}

//NewValues creates an empty *Values.
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

//Merge merges all associations in other into v starting at key.
//To merge other in v at root, use an empty Key for key.
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

//EachKeyValue calls visitor with each set Key value association in v.
func (v *Values) EachKeyValue(visitor func(key Key, value interface{})) {
	v.lock.RLock()
	defer v.lock.RUnlock()

	v.root.eachKeyValue(nil, visitor)
}

//Equal determines whether or not v and other contain the exact same set of
//Keys and associated values.
//Comparison on a value by value basis is done with the == operator.
func (v *Values) Equal(other *Values) bool {
	v.lock.RLock()
	defer v.lock.RUnlock()
	other.lock.RLock()
	defer other.lock.RUnlock()

	return v.root.equal(other.root)
}

//Put adds the key, value association to v.
//changed indicates whether or not this operation changes the set of associations
//in any way.
func (v *Values) Put(key Key, value interface{}) (changed bool) {
	v.lock.Lock()
	defer v.lock.Unlock()

	return v.put(key, value)
}

func (v *Values) put(key Key, value interface{}) bool {
	return v.root.put(key, value)
}

//IsEmpty determines whether or not any associated exist in v.
func (v *Values) IsEmpty() bool {
	v.lock.RLock()
	defer v.lock.RUnlock()

	return v.root.isEmpty()
}

//Get returns the value associated with key or nil if the association does not
//exist.
func (v *Values) Get(key Key) (value interface{}) {
	value, _ = v.GetOk(key)
	return value
}

//GetOk return the value associated with key.
//The return value ok indicates whether or not any value is actually stored at key.
func (v *Values) GetOk(key Key) (value interface{}, ok bool) {
	v.lock.RLock()
	defer v.lock.RUnlock()

	_, found, _ := v.root.findDescendent(nil, key, false, false)
	if found == nil {
		return nil, false
	}
	if found.isSet() {
		return found.value, true
	}
	return newValues(found.clone()), true
}

//Clone creates a new Values with all associations copied into the result.
//The individual values are shallow copied into the result.
func (v *Values) Clone() *Values {
	v.lock.RLock()
	defer v.lock.RUnlock()

	return newValues(v.root)
}

//Remove deletes the association stored at key if one exists.
//v is the removed value or nil if nothing was returned, and ok indicated any value
//was actually removed.
func (v *Values) Remove(key Key) (value interface{}, ok bool) {
	v.lock.Lock()
	defer v.lock.Unlock()

	if key.IsEmpty() {
		if v.root.isSet() {
			result := v.root.value
			v.root = newNode()
			return result, true
		}
		return nil, false
	}

	parent, found, _ := v.root.findDescendent(nil, key, false, false)
	if found == nil {
		return nil, false
	}
	var result interface{}
	if found.isSet() {
		result = found.value
	} else {
		result = newValues(found)
	}
	if parent != nil {
		delete(parent.children, key[key.Len()-1])
	}
	return result, true
}

//node is the internal node type for a Values tree.
//It stores the value stored at its location in the tree and links to children node.
type node struct {
	//value is the value stored at this location in the tree.
	//It may be nil.
	//It is determined to be set by the isSet() method.
	value interface{}

	//children holds the references to this node's child nodes.
	//The keys in children are the parts to the larger Key that references a value.
	children map[string]*node
}

//newNodeValues creates a *node to value set to value and children set to nil.
func newNodeValue(value interface{}) *node {
	n := newNode()
	n.value = value
	n.children = nil
	return n
}

//newNode creates a *node with nil value and empty children.
func newNode() *node {
	return &node{
		value:    nil,
		children: map[string]*node{},
	}
}

//isEmpty determines whether or not n is empty.
//true if not n.isSet() and length of n.children is 0, false otherwise.
func (n *node) isEmpty() bool {
	return !n.isSet() && len(n.children) == 0
}

//isSet determines whether or not n is set (n.values can assumed to be keyed).
//true if n.children is nil, false otherwise.
func (n *node) isSet() bool {
	return n.children == nil
}

//put puts value at key within n's subtree or at n if key is empty.
func (n *node) put(key Key, value interface{}) bool {
	if key.IsEmpty() {
		return n.setValue(value)
	}
	child, changed := n.findChild(key[0], true)
	remainingKey := key[1:]
	return child.put(remainingKey, value) || changed
}

//setValue sets n value to value.
//if value is a *Values, then n.setValues(values.(*Values)) is used.
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

//setValues calls n.put() for each key, value in values.
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

//clone returns a cloned n with a shallow copy of n.value and n.children cloned
//via n.cloneChildren().
func (n *node) clone() *node {
	return &node{
		value:    n.value,
		children: n.cloneChildren(),
	}
}

//cloneChildren returns a new map matching n.children by calling *node.clone()
//on each child in n.children.
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

//eachKeyValues calls visitor for each set value in n's subtree.
//key is a Key onto which to append all descdendent keys.
func (n *node) eachKeyValue(key Key, visitor func(key Key, value interface{})) {
	if n.isSet() {
		visitor(key, n.value)
		return
	}
	for keyPart, childNode := range n.children {
		childNode.eachKeyValue(append(NewKey(key...), keyPart), visitor)
	}
}

//equal determines if n and other are equal by n.value == other.value and n.childrenEqual(other)
func (n *node) equal(other *node) bool {
	return n.value == other.value && n.childrenEqual(other)
}

//childrenEqual determines if n and other's children have the same set of keys
//and all associated child nodes are equal via *node.equal().
func (n *node) childrenEqual(other *node) bool {
	if len(n.children) != len(other.children) {
		return false
	}
	for keyPart, child := range n.children {
		otherChild, ok := other.children[keyPart]
		if !ok || !child.equal(otherChild) {
			return false
		}
	}
	return true
}

//findDescendent finds a desired descendent (node at any lower level) and its possible
//parent.
//The parameter parent should be n's parent. This must be ensured by the caller.
func (n *node) findDescendent(parent *node, key Key, canChange, hasChanged bool) (foundParent *node, found *node, changed bool) {
	//if there is no more key to look into, then return parent, n, and hasChanged.
	if key.IsEmpty() {
		return parent, n, hasChanged
	}
	//find the child at the next part in key, possibly creating or changing with canChange.
	child, changed := n.findChild(key[0], canChange)
	if child == nil {
		return parent, nil, changed
	}
	//recurse into children.
	return child.findDescendent(n, key[1:], canChange, changed || hasChanged)
}

//findChild find a child node of n at keyPart.
//canChange tells findChild if it can modify n.children to modify n.value and or n.children.
//Essentially, canChange being true performs a put, and false performs a get or lookup.
func (n *node) findChild(keyPart string, canChange bool) (*node, bool) {
	changed := false
	if n.isSet() {
		if !canChange {
			return nil, false
		}
		n.value, changed = nil, true
		n.children = map[string]*node{}
	}
	child, ok := n.children[keyPart]
	if !ok {
		if !canChange {
			return nil, false
		}
		n.children[keyPart] = newNode()
		child = n.children[keyPart]
		changed = true
	}
	return child, changed
}
