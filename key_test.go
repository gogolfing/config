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

/*
func TestNewKeyValue(t *testing.T) {
	key := NewKeySep("key", ".")
	keyValue := NewKeyValue(key, "value")
	if !reflect.DeepEqual(keyValue.Key, key) || keyValue.Value != "value" {
		t.Fail()
	}
}
*/
