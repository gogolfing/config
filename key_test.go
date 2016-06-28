package config

import (
	"reflect"
	"testing"
)

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
