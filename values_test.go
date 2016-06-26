package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewValues(t *testing.T) {
	v := NewValues()
	if v.lock == nil || v.root == nil {
		t.Fail()
	}
	testNode(t, v.root, nil, false, true)
}

func TestValues_Merge(t *testing.T) {
	//merging single value into non single value
	//merging non single value into single value
	//merging non single value into non single value with no change
	//merging non single value into non single value with change
	//merging not at root for all above

	tests := []struct {
		values  *Values
		toMerge *Values
		mergeAt Key
		changed bool
		result  *Values
	}{
		//merging empty into empty
		{
			NewValues(),
			NewValues(),
			nil,
			false,
			NewValues(),
		},
		//merging empty into non empty
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
			NewValues(),
			nil,
			false,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
		},
		//merging non empty into empty
		{
			NewValues(),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
			nil,
			true,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
		},
		//merging single value into single value with no change
		//merging single value into single value with change
		{},
	}
	for index, test := range tests {
		changed := test.values.Merge(test.mergeAt, test.toMerge)
		if changed != test.changed {
			t.Errorf("%v, *Values.Merge() changed = %v WANT %v", index, changed, test.changed)
		}
		if !reflect.DeepEqual(test.values, test.result) {
			t.Errorf("%v, *Values.Merge() = %v WANT %v", index, test.values, test.result)
		}
	}
}

func TestValues_EachKeyValue(t *testing.T) {
	tests := []struct {
		values *Values
		result map[string]interface{}
	}{
		//empty values should have empty result
		{
			func() *Values {
				v := NewValues()
				return v
			}(),
			map[string]interface{}{},
		},
		//single value should have single result with empty string key as a result of strings.Join()
		{
			func() *Values {
				v := NewValues()
				v.Put(nil, "something")
				return v
			}(),
			map[string]interface{}{
				"": "something",
			},
		},
		//all nested values should have correct key, value pair
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a", "b"), 0)
				v.Put(NewKey("a", "c"), 1)
				v.Put(NewKey("d"), 2)
				return v
			}(),
			map[string]interface{}{
				"a.b": 0,
				"a.c": 1,
				"d":   2,
			},
		},
	}
	for _, test := range tests {
		result := map[string]interface{}{}
		test.values.EachKeyValue(func(key Key, value interface{}) {
			result[strings.Join(key, ".")] = value
		})
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("%v, *Values.EachKeyValue() = %v WANT %v")
		}
	}
}

func TestValues_Put(t *testing.T) {
	tests := []struct {
		before  []KeyValue
		key     Key
		value   interface{}
		changed bool
		root    *node
	}{
		//nil key with nil value
		{
			nil,
			nil,
			nil,
			true,
			&node{
				value:    nil,
				set:      true,
				children: nil,
			},
		},
		//nil key with nil value with changed false
		{
			[]KeyValue{
				{nil, nil},
			},
			nil,
			nil,
			false,
			&node{
				value:    nil,
				set:      true,
				children: nil,
			},
		},
		//empty key with nil value
		{
			nil,
			NewKey(),
			nil,
			true,
			&node{
				value:    nil,
				set:      true,
				children: nil,
			},
		},
		//empty key with nil value with changed false
		{
			[]KeyValue{
				{NewKey(), nil},
			},
			NewKey(),
			nil,
			false,
			&node{
				value:    nil,
				set:      true,
				children: nil,
			},
		},
		//nil key with actual value
		{
			nil,
			nil,
			true,
			true,
			&node{
				value:    true,
				set:      true,
				children: nil,
			},
		},
		//nil key with actual value with changed false
		{
			[]KeyValue{
				{nil, false},
			},
			nil,
			false,
			false,
			&node{
				value:    false,
				set:      true,
				children: nil,
			},
		},
		//empty key with actual value
		{
			nil,
			NewKey(),
			true,
			true,
			&node{
				value:    true,
				set:      true,
				children: nil,
			},
		},
		//empty key with actual value with changed false
		{
			[]KeyValue{
				{nil, false},
			},
			NewKey(),
			false,
			false,
			&node{
				value:    false,
				set:      true,
				children: nil,
			},
		},
		//one-level key with nil value
		{
			nil,
			NewKey("hello"),
			nil,
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"hello": &node{
						value:    nil,
						set:      true,
						children: nil,
					},
				},
			},
		},
		//one-level key with nil value with changed false
		{
			[]KeyValue{
				{NewKey("hello"), nil},
			},
			NewKey("hello"),
			nil,
			false,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"hello": &node{
						value:    nil,
						set:      true,
						children: nil,
					},
				},
			},
		},
		//one-level key with actual value
		{
			nil,
			NewKey("hello"),
			2,
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"hello": &node{
						value:    2,
						set:      true,
						children: nil,
					},
				},
			},
		},
		//one-level key with actual value with changed false
		{
			[]KeyValue{
				{NewKey("hello"), 2},
			},
			NewKey("hello"),
			2,
			false,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"hello": &node{
						value:    2,
						set:      true,
						children: nil,
					},
				},
			},
		},
		//one-level key with actual value with changed true
		{
			[]KeyValue{
				{NewKey("hello"), 2},
			},
			NewKey("hello"),
			"two",
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"hello": &node{
						value:    "two",
						set:      true,
						children: nil,
					},
				},
			},
		},
		//two-level key with nil value
		{
			nil,
			NewKey("hello", "world"),
			nil,
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"hello": &node{
						value: nil,
						set:   false,
						children: map[string]*node{
							"world": &node{
								value:    nil,
								set:      true,
								children: nil,
							},
						},
					},
				},
			},
		},
		//two-level key with nil value with changed false
		{
			[]KeyValue{
				{NewKey("hello", "world"), nil},
			},
			NewKey("hello", "world"),
			nil,
			false,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"hello": &node{
						value: nil,
						set:   false,
						children: map[string]*node{
							"world": &node{
								value:    nil,
								set:      true,
								children: nil,
							},
						},
					},
				},
			},
		},
		//two-level key with actual value
		{
			nil,
			NewKey("hello", "world"),
			"something",
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"hello": &node{
						value: nil,
						set:   false,
						children: map[string]*node{
							"world": &node{
								value:    "something",
								set:      true,
								children: nil,
							},
						},
					},
				},
			},
		},
		//two-level key with actual value with changed false
		{
			[]KeyValue{
				{NewKey("hello", "world"), "something"},
			},
			NewKey("hello", "world"),
			"something",
			false,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"hello": &node{
						value: nil,
						set:   false,
						children: map[string]*node{
							"world": &node{
								value:    "something",
								set:      true,
								children: nil,
							},
						},
					},
				},
			},
		},
		//overwrite existing value for change
		{
			[]KeyValue{
				{NewKey("a"), "a"},
			},
			NewKey("a", "b"),
			"b",
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"a": &node{
						value: nil,
						set:   false,
						children: map[string]*node{
							"b": &node{
								value:    "b",
								set:      true,
								children: nil,
							},
						},
					},
				},
			},
		},
		//parallel write in "object" with change
		{
			[]KeyValue{
				{NewKey("a"), "a"},
			},
			NewKey("b"),
			"b",
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"a": newNodeValue("a"),
					"b": newNodeValue("b"),
				},
			},
		},
		//parallel write in "object" with no change
		{
			[]KeyValue{
				{NewKey("a"), "a"},
				{NewKey("b"), "b"},
			},
			NewKey("b"),
			"b",
			false,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"a": newNodeValue("a"),
					"b": newNodeValue("b"),
				},
			},
		},
		//underwrite with change
		{
			[]KeyValue{
				{NewKey("a", "b"), "b"},
				{NewKey("a", "c"), "c"},
			},
			NewKey("a"),
			"a",
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"a": &node{
						value:    "a",
						set:      true,
						children: nil,
					},
				},
			},
		},
		//underwrite in "onject" with change
		{
			[]KeyValue{
				NewKeyValue(NewKey("else"), "else"),
				NewKeyValue(NewKey("a", "b"), "b"),
				NewKeyValue(NewKey("a", "c"), "c"),
			},
			NewKey("a"),
			"a",
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"else": &node{
						value:    "else",
						set:      true,
						children: nil,
					},
					"a": &node{
						value:    "a",
						set:      true,
						children: nil,
					},
				},
			},
		},
		//empty *Values at empty root
		{
			nil,
			nil,
			NewValues(),
			false,
			&node{
				value:    nil,
				set:      false,
				children: nil,
			},
		},
		//single value *Values at empty root
		{
			nil,
			nil,
			"value",
			true,
			&node{
				value:    "value",
				set:      true,
				children: nil,
			},
		},
		//single value *Values at non empty root
		{
			[]KeyValue{
				NewKeyValue(NewKey("a"), "a"),
				NewKeyValue(NewKey("b"), "b"),
			},
			nil,
			"value",
			true,
			&node{
				value:    "value",
				set:      true,
				children: nil,
			},
		},
		//single value *Values at single value root with no change
		{
			[]KeyValue{
				NewKeyValue(nil, "value"),
			},
			nil,
			"value",
			false,
			&node{
				value:    "value",
				set:      true,
				children: nil,
			},
		},
		//single value *Values at single value root with change
		{
			[]KeyValue{
				NewKeyValue(nil, 2),
			},
			nil,
			"value",
			true,
			&node{
				value:    "value",
				set:      true,
				children: nil,
			},
		},
		//empty *Values at single value root
		{
			[]KeyValue{
				NewKeyValue(NewKey("a"), "a"),
			},
			nil,
			NewValues(),
			true,
			&node{
				value:    nil,
				set:      false,
				children: nil,
			},
		},
		//empty *Values at non empty root
		{
			[]KeyValue{
				NewKeyValue(NewKey("a"), "a"),
				NewKeyValue(NewKey("b"), "b"),
			},
			nil,
			NewValues(),
			true,
			&node{
				value:    nil,
				set:      false,
				children: nil,
			},
		},
		//non empty *Values at empty root
		{
			nil,
			nil,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"a": newNodeValue("a"),
					"b": newNodeValue("b"),
				},
			},
		},
		//empty *Value at deep "object"
		{
			[]KeyValue{
				NewKeyValue(NewKey("a", "b"), "b"),
				NewKeyValue(NewKey("a", "c"), "c"),
			},
			NewKey("a"),
			NewValues(),
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"a": &node{
						value:    nil,
						set:      false,
						children: nil,
					},
				},
			},
		},
		//deep *Values overwrites single *Values
		{
			[]KeyValue{
				NewKeyValue(NewKey("a"), "a"),
			},
			NewKey("a"),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("b"), "b")
				v.Put(NewKey("c"), "c")
				return v
			}(),
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"a": &node{
						value: nil,
						set:   false,
						children: map[string]*node{
							"b": newNodeValue("b"),
							"c": newNodeValue("c"),
						},
					},
				},
			},
		},
		//deep *Values overwrites deep *Values
		{
			[]KeyValue{
				NewKeyValue(NewKey("a", "b"), "b"),
				NewKeyValue(NewKey("a", "c"), "c"),
			},
			NewKey("a"),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("c"), "new c")
				v.Put(NewKey("d"), "d")
				v.Put(NewKey("e", "f"), "f")
				return v
			}(),
			true,
			&node{
				value: nil,
				set:   false,
				children: map[string]*node{
					"a": &node{
						value: nil,
						set:   false,
						children: map[string]*node{
							"b": newNodeValue("b"),
							"c": newNodeValue("new c"),
							"d": newNodeValue("d"),
							"e": &node{
								value: nil,
								set:   false,
								children: map[string]*node{
									"f": newNodeValue("f"),
								},
							},
						},
					},
				},
			},
		},
	}
	for index, test := range tests {
		v := NewValues()
		for _, kv := range test.before {
			v.Put(kv.Key, kv.Value)
		}
		changed := v.Put(test.key, test.value)
		if changed != test.changed {
			t.Errorf("%v, v.Put(%v) changed = %v WANT %v", index, test.key, changed, test.changed)
		}
		if !reflect.DeepEqual(v.root, test.root) {
			t.Errorf("%v, v.Put(%v) root = %v WANT %v", index, test.key, v.root, test.root)
		}
	}
}

func TestValues_IsEmpty(t *testing.T) {
	tests := []struct {
		values  *Values
		isEmpty bool
	}{
		{
			NewValues(),
			true,
		},
		{
			func() *Values {
				v := NewValues()
				v.Put(nil, "something")
				return v
			}(),
			false,
		},
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
			false,
		},
	}
	for index, test := range tests {
		isEmpty := test.values.IsEmpty()
		if isEmpty != test.isEmpty {
			t.Errorf("%v, *Values.IsEmpty() = %v WANT %v", index, isEmpty, test.isEmpty)
		}
	}
}

func TestNewNodeValue(t *testing.T) {
	tests := []struct {
		value interface{}
	}{
		{nil},
		{1},
		{"hello"},
		{[]string{"hello"}},
	}
	for _, test := range tests {
		n := newNodeValue(test.value)
		testNode(t, n, test.value, true, true)
	}
}

func TestNewNodeChildren(t *testing.T) {
	tests := []struct {
		children map[string]*node
	}{
		{nil},
		{map[string]*node{}},
	}
	for _, test := range tests {
		n := newNodeChildren(test.children)
		testNode(t, n, nil, false, test.children == nil)
	}
}

func TestNewNode(t *testing.T) {
	n := newNode()
	if n.value != nil || n.set != false || n.children != nil {
		t.Fail()
	}
}

func testNode(t *testing.T, n *node, value interface{}, set, childrenNil bool) {
	if n == nil {
		t.Error("*node should not be nil")
	}
	if !reflect.DeepEqual(n.value, value) {
		t.Error("n.value should reflect.DeepEqual() value", n.value, value)
	}
	if n.set != set {
		t.Error("n.set should equal set", n.set, set)
	}
	if n.children == nil != childrenNil {
		t.Error("n.children should have nil value %v, got %v", childrenNil, n.children == nil)
	}
}
