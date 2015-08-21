package config

import "testing"

var options struct {
	Name        string
	Description string
}

func TestConfigFileCanBeOpenedAndParsed(t *testing.T) {
	New("test.json")
	Parse()
	value := GetString("object.nested.nested_object", "fail")
	if value == "fail" {
		t.Error("Expected 'Nested', got ", value)
	}
}
