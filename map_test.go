package config

import (
	"reflect"
	"testing"
)

func TestNewMap(t *testing.T) {
	result := NewMap()
	if result == nil || len(result) != 0 {
		t.Fail()
	}
}

func TestMap_Clone_empty(t *testing.T) {
	result := NewMap().Clone()
	if result == nil || len(result) != 0 {
		t.Fail()
	}
}

func TestMap_Clone_primitive(t *testing.T) {
	m := NewMap()
	m["one"] = 1
	m["bool"] = false
	inner := map[string]interface{}{
		"hello": "world",
	}
	m["inner"] = inner
	p := &struct{}{}
	m["p"] = p

	cloned := m.Clone()
	if len(cloned) != 4 {
		t.Fail()
	}
	if cloned["one"] != 1 || cloned["bool"] != false || cloned["p"] != p {
		t.Fail()
	}
	if !reflect.DeepEqual(cloned["inner"], inner) {
		t.Fail()
	}
}

func TestMap_Clone_slice(t *testing.T) {
	m := NewMap()
	value := []interface{}{}
	value = append(value, 1)
	value = append(value, NewMap())
	key := NewKey("slice", ".")
	m.Put(key, value)

	if value, ok := m.GetOk(key); !ok || !reflect.DeepEqual(value, []interface{}{1, NewMap()}) {
		t.Fail()
	}
}

func TestMap_Clone_Map(t *testing.T) {
	m := NewMap()
	m["one"] = "one"

	inner2 := NewMap()
	inner := NewMap()
	inner["inner2"] = inner2
	m["inner"] = inner

	cloned := m.Clone()
	if len(cloned) != 2 {
		t.Fail()
	}
	if !reflect.DeepEqual(m, cloned) {
		t.Fail()
	}

	clonedInner := cloned["inner"].(Map)
	clonedInner2 := clonedInner["inner2"].(Map)
	if &inner == &clonedInner {
		t.Fail()
	}
	if &inner2 == &clonedInner2 {
		t.Fail()
	}

	clonedInner["inner2"] = "inner2"
	if inner["inner2"] == "inner2" {
		t.Fail()
	}
}
