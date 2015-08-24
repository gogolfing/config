package config

import "testing"

func TestMain(m *testing.M) {
	New("test.json")
	Parse()
	m.Run()
}

func TestConfigFileCanHandleNestedObjects(t *testing.T) {
	value := GetString("object.nested.nested_object", "fail")
	if value == "fail" {
		t.Error("Expected 'Nested', got ", value)
	}
}

func TestBooleansCanBeRetrievedFromConfigFile(t *testing.T) {
	value := GetBool("bool_test", false)
	if value == false {
		t.Error("Expected true, got ", value)
	}
}

func TestIntegersCanBeRetrievedFromConfigFile(t *testing.T) {
	value := GetInt("int_test", 0)
	if value != 1234 {
		t.Error("Expected 1234, got ", value)
	}
}

func TestFloatsCanBeRetrievedFromConfigFile(t *testing.T) {
	value := GetFloat64("float_test", 0.1)
	if value != 3.14 {
		t.Error("Expected 3.14, got ", value)
	}
}
