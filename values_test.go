package config

import (
	"reflect"
	"strings"
	"testing"
)

type KeyValue struct {
	Key   Key
	Value interface{}
}

func NewKeyValue(key Key, value interface{}) KeyValue {
	return KeyValue{key, value}
}

func TestNewValues(t *testing.T) {
	v := NewValues()
	if v.lock == nil || v.root == nil {
		t.Fail()
	}
	testNode(t, v.root, nil, false)
}

func TestValues_Merge(t *testing.T) {
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
		{
			func() *Values {
				v := NewValues()
				v.Put(nil, "a")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(nil, "a")
				return v
			}(),
			nil,
			false,
			func() *Values {
				v := NewValues()
				v.Put(nil, "a")
				return v
			}(),
		},
		//merging single value into single value with change
		{
			func() *Values {
				v := NewValues()
				v.Put(nil, "a")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(nil, "b")
				return v
			}(),
			nil,
			true,
			func() *Values {
				v := NewValues()
				v.Put(nil, "b")
				return v
			}(),
		},
		//merging single value into non single value
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(nil, "c")
				return v
			}(),
			nil,
			true,
			func() *Values {
				v := NewValues()
				v.Put(nil, "c")
				return v
			}(),
		},
		//merging non single value into single value
		{
			func() *Values {
				v := NewValues()
				v.Put(nil, "a")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
			nil,
			true,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
		},
		//merging non single value into non single value with no change
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
			nil,
			false,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
		},
		//merging non single value into non single value with change
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("b"), "new b")
				v.Put(NewKey("c"), "c")
				v.Put(NewKey("d"), NewValues())
				return v
			}(),
			nil,
			true,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "new b")
				v.Put(NewKey("c"), "c")
				return v
			}(),
		},
		//merging empty into empty not at root
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), NewValues())
				return v
			}(),
			NewValues(),
			NewKey("a"),
			false,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), NewValues())
				return v
			}(),
		},
		//merging empty into non empty not at root
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
			NewValues(),
			NewKey("a"),
			false,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
		},
		//merging non empty into empty not at root
		{
			NewValues(),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
			NewKey("notroot"),
			true,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("notroot", "a"), "a")
				return v
			}(),
		},
		//merging single value into single value with no change not at root
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(nil, "a")
				return v
			}(),
			NewKey("a"),
			false,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				return v
			}(),
		},
		//merging single value into single value with change not at root
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k"), "a")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(nil, "b")
				return v
			}(),
			NewKey("k"),
			true,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k"), "b")
				return v
			}(),
		},
		//merging single value into non single value not at root
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k", "a"), "a")
				v.Put(NewKey("k", "b"), "b")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(nil, "c")
				return v
			}(),
			NewKey("k"),
			true,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k"), "c")
				return v
			}(),
		},
		//merging non single value into single value not at root
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k"), "a")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
			NewKey("k"),
			true,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k", "a"), "a")
				v.Put(NewKey("k", "b"), "b")
				return v
			}(),
		},
		//merging non single value into non single value with no change not at root
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k", "a"), "a")
				v.Put(NewKey("k", "b"), "b")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a"), "a")
				v.Put(NewKey("b"), "b")
				return v
			}(),
			NewKey("k"),
			false,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k", "a"), "a")
				v.Put(NewKey("k", "b"), "b")
				return v
			}(),
		},
		//merging non single value into non single value with change not at root
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k", "a"), "a")
				v.Put(NewKey("k", "b"), "b")
				return v
			}(),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("b"), "new b")
				v.Put(NewKey("c"), "c")
				v.Put(NewKey("d"), NewValues())
				return v
			}(),
			NewKey("k"),
			true,
			func() *Values {
				v := NewValues()
				v.Put(NewKey("k", "a"), "a")
				v.Put(NewKey("k", "b"), "new b")
				v.Put(NewKey("k", "c"), "c")
				return v
			}(),
		},
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
				children: map[string]*node{
					"hello": &node{
						value:    nil,
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
				children: map[string]*node{
					"hello": &node{
						value:    nil,
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
				children: map[string]*node{
					"hello": &node{
						value:    2,
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
				children: map[string]*node{
					"hello": &node{
						value:    2,
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
				children: map[string]*node{
					"hello": &node{
						value:    "two",
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
				children: map[string]*node{
					"hello": &node{
						value: nil,
						children: map[string]*node{
							"world": &node{
								value:    nil,
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
				children: map[string]*node{
					"hello": &node{
						value: nil,
						children: map[string]*node{
							"world": &node{
								value:    nil,
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
				children: map[string]*node{
					"hello": &node{
						value: nil,
						children: map[string]*node{
							"world": &node{
								value:    "something",
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
				children: map[string]*node{
					"hello": &node{
						value: nil,
						children: map[string]*node{
							"world": &node{
								value:    "something",
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
				children: map[string]*node{
					"a": &node{
						value: nil,
						children: map[string]*node{
							"b": &node{
								value:    "b",
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
				children: map[string]*node{
					"a": &node{
						value:    "a",
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
				children: map[string]*node{
					"else": &node{
						value:    "else",
						children: nil,
					},
					"a": &node{
						value:    "a",
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
			newNode(),
		},
		//single value *Values at empty root
		{
			nil,
			nil,
			"value",
			true,
			&node{
				value:    "value",
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
				children: map[string]*node{
					"a": &node{
						value:    nil,
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
				children: map[string]*node{
					"a": &node{
						value: nil,
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
				children: map[string]*node{
					"a": &node{
						value: nil,
						children: map[string]*node{
							"b": newNodeValue("b"),
							"c": newNodeValue("new c"),
							"d": newNodeValue("d"),
							"e": &node{
								value: nil,
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

func areChildrenEqual(a, b map[string]*node) bool {
	if a == nil || b == nil {
		return (a == nil) == (b == nil)
	}
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
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
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey("a", "b", "c"), "c")
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

func TestValues_Get(t *testing.T) {
	values := NewValues()
	values.Put(NewKey(""), "")
	values.Put(NewKey("nil"), nil)
	values.Put(NewKey("a"), "a")
	values.Put(NewKey("b", "c"), "c")
	tests := []struct {
		key    Key
		result interface{}
	}{
		{NewKey(), values},
		{NewKey(""), ""},
		{NewKey("nil"), nil},
		{NewKey("a"), "a"},
		{
			NewKey("b"),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("c"), "c")
				return v
			}(),
		},
		{NewKey("b", "c"), "c"},
		{NewKey("d"), nil},
	}
	for _, test := range tests {
		result := values.Get(test.key)
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("%v, reflect.DeepEqual(%v, %v) should be true", test.key, result, test.result)
		}
	}
}

func TestValues_GetOk(t *testing.T) {
	values := NewValues()
	values.Put(NewKey(""), "")
	values.Put(NewKey("nil"), nil)
	values.Put(NewKey("a"), "a")
	values.Put(NewKey("b", "c"), "c")
	tests := []struct {
		key    Key
		result interface{}
		ok     bool
	}{
		{NewKey(), values, true},
		{NewKey(""), "", true},
		{NewKey("nil"), nil, true},
		{NewKey("a"), "a", true},
		{
			NewKey("b"),
			func() *Values {
				v := NewValues()
				v.Put(NewKey("c"), "c")
				return v
			}(),
			true,
		},
		{NewKey("b", "c"), "c", true},
		{NewKey("d"), nil, false},
	}
	for _, test := range tests {
		result, ok := values.GetOk(test.key)
		if ok != test.ok {
			t.Errorf("%v, ok = %v WANT %v", ok, test.ok)
		}
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("%v, reflect.DeepEqual(%v, %v) should be true", test.key, result, test.result)
		}
	}
}

func TestValues_Clone(t *testing.T) {
	tests := []struct {
		values *Values
	}{
		{NewValues()},
		{
			func() *Values {
				v := NewValues()
				v.Put(NewKey(), "value at emtpy key")
				return v
			}(),
		},
		{
			func() *Values {
				valueValues := NewValues()
				valueValues.Put(NewKey("a"), "a")

				v := NewValues()
				v.Put(NewKey(""), "")
				v.Put(NewKey("nil"), nil)
				v.Put(NewKey("one", "slice"), []string{"hello", "world"})
				v.Put(NewKey("one", "int"), 1234)
				v.Put(NewKey("one", "bool"), false)
				v.Put(NewKey("one", "valueValues"), valueValues)
				return v
			}(),
		},
	}
	for index, test := range tests {
		clone := test.values.Clone()
		if !reflect.DeepEqual(clone, test.values) {
			t.Errorf("%v, reflect.DeepEqual(%v, %v) should be true", index, clone, test.values)
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
		testNode(t, n, test.value, true)
	}
}

func TestNewNode(t *testing.T) {
	n := newNode()
	if n.value != nil || n.children == nil || len(n.children) != 0 {
		t.Fail()
	}
}

func testNode(t *testing.T, n *node, value interface{}, childrenNil bool) {
	if n == nil {
		t.Error("*node should not be nil")
	}
	if !reflect.DeepEqual(n.value, value) {
		t.Error("n.value should reflect.DeepEqual() value", n.value, value)
	}
	if n.children == nil != childrenNil {
		t.Errorf("n.children should have nil value %v, got %v", childrenNil, n.children == nil)
	}
}
