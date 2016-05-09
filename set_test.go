package config

import (
	"reflect"
	"testing"
)

func TestNewSet(t *testing.T) {
	result := NewSet()
	if result == nil || len(result) != 0 {
		t.Fail()
	}
}

func TestSet_clone_empty(t *testing.T) {
	result := NewSet().clone()
	if result == nil || len(result) != 0 {
		t.Fail()
	}
}

func TestSet_clone_primitive(t *testing.T) {
	s := NewSet()
	s["one"] = 1
	s["bool"] = false
	m := map[string]interface{}{
		"hello": "world",
	}
	s["m"] = m
	p := &struct{}{}
	s["p"] = p

	cloned := s.clone()
	if len(cloned) != 4 {
		t.Fail()
	}
	if cloned["one"] != 1 || cloned["bool"] != false || cloned["p"] != p {
		t.Fail()
	}
	if !reflect.DeepEqual(cloned["m"], m) {
		t.Fail()
	}
}

func TestSet_clone_Set(t *testing.T) {
	s := NewSet()
	s["one"] = "one"

	inner2 := NewSet()
	inner := NewSet()
	inner["inner2"] = inner2
	s["inner"] = inner

	cloned := s.clone()
	if len(cloned) != 2 {
		t.Fail()
	}
	if !reflect.DeepEqual(s, cloned) {
		t.Fail()
	}

	clonedInner := cloned["inner"].(Set)
	clonedInner2 := clonedInner["inner2"].(Set)
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
