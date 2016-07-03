package config

import (
	"reflect"
	"testing"
)

func TestNewKey(t *testing.T) {
	key := NewKey()
	if !key.IsEmpty() {
		t.Fail()
	}

	key = Key(nil)
	if !key.IsEmpty() {
		t.Fail()
	}

	key = NewKey("a", "", "b")
	if key.Len() != 3 || key[0] != "a" || key[1] != "" || key[2] != "b" {
		t.Fail()
	}
}

func TestNewKeySep(t *testing.T) {
	tests := []struct {
		value  string
		sep    string
		result []string
	}{
		{"", "", []string{}},
		{"", ".", []string{""}},
		{"hello", "", []string{"h", "e", "l", "l", "o"}},
		{"one . two ", ".", []string{"one ", " two "}},
		{"one.two.three.four", ".", []string{"one", "two", "three", "four"}},
		{"one, two, three", ", ", []string{"one", "two", "three"}},
	}
	for _, test := range tests {
		result := NewKeySep(test.value, test.sep)
		if result == nil {
			t.Error("result cannot be nil")
		}
		if len(result) == 0 && len(test.result) == 0 {
			continue
		}
		if !reflect.DeepEqual(result, Key(test.result)) {
			t.Error(result, test.result)
		}
	}
}

func TestKey_IsEmtpy(t *testing.T) {
	tests := []struct {
		key    Key
		result bool
	}{
		{NewKey(), true},
		{Key(nil), true},
		{Key([]string{}), true},
		{NewKey([]string{}...), true},
		{NewKey(""), false},
		{NewKey([]string{"hello", "world"}...), false},
		{NewKey("a", "b", "c"), false},
	}
	for _, test := range tests {
		result := test.key.IsEmpty()
		if result != test.result {
			t.Errorf("%v.IsEmpty() = %v WANT %v", test.key, result, test.result)
		}
	}
}

func TestKey_Len(t *testing.T) {
	tests := []struct {
		key    Key
		result int
	}{
		{NewKey(), 0},
		{Key(nil), 0},
		{Key([]string{}), 0},
		{NewKey([]string{}...), 0},
		{NewKey(""), 1},
		{NewKey([]string{"hello", "world"}...), 2},
		{NewKey("a", "b", "c"), 3},
	}
	for _, test := range tests {
		result := test.key.Len()
		if result != test.result {
			t.Errorf("%v.Len() = %v WANT %v", test.key, result, test.result)
		}
	}
}

func TestKey_Equals(t *testing.T) {
	tests := []struct {
		key    Key
		other  Key
		result bool
	}{
		{NewKey(), NewKey(), true},
		{NewKey(), Key(nil), true},
		{NewKey(), Key([]string{}), true},
		{NewKey(), NewKey([]string{}...), true},
		{NewKey(""), NewKey(), false},
		{NewKey(""), Key(nil), false},
		{NewKey(""), Key([]string{}), false},
		{NewKey(""), NewKey([]string{}...), false},
		{NewKey("a", "b"), NewKey("a", "b"), true},
		{NewKey("a", "b"), NewKey("a"), false},
		{NewKey("a", "b"), NewKey("a", "b", "c"), false},
	}
	for _, test := range tests {
		result := test.key.Equals(test.other)
		if result != test.result {
			t.Errorf("%v.Equals(%v) = %v WANT %v", test.key, test.other, result, test.result)
		}
	}
}

func TestKey_StartsWith(t *testing.T) {
	tests := []struct {
		key    Key
		other  Key
		result bool
	}{
		{NewKey(), NewKey(), true},
		{NewKey(""), NewKey(""), true},
		{NewKey("a", "b"), NewKey("a"), true},
		{NewKey("a", "b"), NewKey("a", "b"), true},
		{NewKey(""), NewKey(), true},
		{NewKey("a"), NewKey("b"), false},
		{NewKey("a"), NewKey("a", "b"), false},
		{NewKey("a", "b"), NewKey(""), false},
		{NewKey(), NewKey(""), false},
	}
	for _, test := range tests {
		result := test.key.StartsWith(test.other)
		if result != test.result {
			t.Errorf("%v.StartsWith(%v) = %v WANT %v", test.key, test.other, result, test.result)
		}
	}
}

func TestKey_EndWith(t *testing.T) {
	tests := []struct {
		key    Key
		other  Key
		result bool
	}{
		{NewKey(), NewKey(), true},
		{NewKey(""), NewKey(""), true},
		{NewKey("a", "b"), NewKey("b"), true},
		{NewKey("a", "b"), NewKey("a", "b"), true},
		{NewKey(""), NewKey(), true},
		{NewKey("a"), NewKey("b"), false},
		{NewKey("a"), NewKey("b", "a"), false},
		{NewKey("a", "b"), NewKey(""), false},
		{NewKey(), NewKey(""), false},
	}
	for _, test := range tests {
		result := test.key.EndsWith(test.other)
		if result != test.result {
			t.Errorf("%v.EndsWith(%v) = %v WANT %v", test.key, test.other, result, test.result)
		}
	}
}

func TestKey_Append(t *testing.T) {
	tests := []struct {
		first  Key
		others []Key
		result Key
	}{
		{Key(nil), []Key{}, Key(nil)},
		{NewKey(), []Key{NewKey()}, NewKey()},
		{Key(nil), []Key{NewKey("a")}, NewKey("a")},
		{NewKey("a"), []Key{}, NewKey("a")},
		{NewKey("a"), []Key{NewKey()}, NewKey("a")},
		{NewKey("a"), []Key{NewKey("b", "c")}, NewKey("a", "b", "c")},
	}
	for _, test := range tests {
		result := test.first.Append(test.others...)
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("%v.Append(%v) = %v WANT %v", test.first, test.others, result, test.result)
		}
	}
}
