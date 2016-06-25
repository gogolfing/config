package config

import (
	"fmt"
	"reflect"
	"testing"
)

// func TestNew(t *testing.T) {
// 	c := New()
// 	if c.Separator != "." {
// 		t.Fail()
// 	}
// 	if c.lock == nil {
// 		t.Fail()
// 	}
// 	if c.m == nil || c.loaders == nil {
// 		t.Fail()
// 	}
// }

// func TestConfif_AddLoader(t *testing.T) {
// 	c := New()
// 	c.AddLoader(intLoader(0))
// 	if c.loaders[0] != intLoader(0) {
// 		t.Fail()
// 	}
// }

// func TestConfig_LoadAll(t *testing.T) {
// 	c := New()
// 	c.AddLoader(intLoader(1))
// 	c.AddLoader(intLoader(2))
// 	c.LoadAll()

// 	if v, ok := c.GetOk("1"); v != intLoader(1) || !ok {
// 		t.Fail()
// 	}
// 	if v, ok := c.GetOk("2"); v != intLoader(2) || !ok {
// 		t.Fail()
// 	}
// }

type intLoader int

func (i intLoader) Load() ([]KeyValue, error) {
	return []KeyValue{
		{Key: NewKey(fmt.Sprint(int(i)), ""), Value: i},
	}, nil
}

func TestConfig_Get(t *testing.T) {
	c := New()
	c.Put("hello", "world")
	tests := []struct {
		key   string
		value interface{}
	}{
		{"hello", "world"},
		{"anything else", nil},
	}
	for _, test := range tests {
		value := c.Get(test.key)
		if value != test.value {
			t.Fail()
		}
	}
}

func TestConfig_GetBool(t *testing.T) {
	c := New()
	c.Put("yes", true)
	if c.GetBool("yes") != true {
		t.Fail()
	}
	if c.GetBool("no") == true {
		t.Fail()
	}
}

func TestConfig_GetBoolOk(t *testing.T) {
	c := New()
	c.Put("yes", true)
	c.Put("false", false)
	if v, ok := c.GetBoolOk("yes"); v != true || !ok {
		t.Fail()
	}
	if v, ok := c.GetBoolOk("no"); v != false || ok {
		t.Fail()
	}
	if v, ok := c.GetBoolOk("false"); v != false || !ok {
		t.Fail()
	}
}

func TestConfig_GetFloat64(t *testing.T) {
	c := New()
	c.Put("one", 1.0)
	if c.GetFloat64("one") != 1.0 {
		t.Fail()
	}
	if c.GetFloat64("no") != 0 {
		t.Fail()
	}
}

func TestConfig_GetFloat64Ok(t *testing.T) {
	c := New()
	c.Put("one", 1.0)
	c.Put("zero", 0.0)
	if v, ok := c.GetFloat64Ok("one"); v != 1.0 || !ok {
		t.Fail()
	}
	if v, ok := c.GetFloat64Ok("no"); v != 0.0 || ok {
		t.Fail()
	}
	if v, ok := c.GetFloat64Ok("zero"); v != 0.0 || !ok {
		t.Fail()
	}
}

func TestConfig_GetInt(t *testing.T) {
	c := New()
	c.Put("one", 1)
	if c.GetInt("one") != 1 {
		t.Fail()
	}
	if c.GetInt("no") != 0 {
		t.Fail()
	}
}

func TestConfig_GetIntOk(t *testing.T) {
	c := New()
	c.Put("int8", int8(1))
	c.Put("uint8", uint8(1))
	c.Put("int16", int16(1))
	c.Put("uint16", uint16(1))
	c.Put("int32", int32(1))
	c.Put("int", int(1))
	c.Put("zero", 0)
	if v, ok := c.GetIntOk("int8"); v != 1 || !ok {
		t.Fail()
	}
	if v, ok := c.GetIntOk("uint8"); v != 1 || !ok {
		t.Fail()
	}
	if v, ok := c.GetIntOk("int16"); v != 1 || !ok {
		t.Fail()
	}
	if v, ok := c.GetIntOk("uint16"); v != 1 || !ok {
		t.Fail()
	}
	if v, ok := c.GetIntOk("int32"); v != 1 || !ok {
		t.Fail()
	}
	if v, ok := c.GetIntOk("int"); v != 1 || !ok {
		t.Fail()
	}
	if v, ok := c.GetIntOk("no"); v != 0 || ok {
		t.Fail()
	}
	if v, ok := c.GetIntOk("zero"); v != 0 || !ok {
		t.Fail()
	}
}

func TestConfig_GetInt64(t *testing.T) {
	c := New()
	c.Put("one", int64(1))
	if c.GetInt64("one") != int64(1) {
		t.Fail()
	}
	if c.GetInt64("no") != int64(0) {
		t.Fail()
	}
}

func TestConfig_GetInt64Ok(t *testing.T) {
	c := New()
	c.Put("int8", int8(1))
	c.Put("uint8", uint8(1))
	c.Put("int16", int16(1))
	c.Put("uint16", uint16(1))
	c.Put("int32", int32(1))
	c.Put("uint32", uint32(1))
	c.Put("int", int(1))
	c.Put("uint", uint(1))
	c.Put("int64", int64(1))
	c.Put("zero", int64(0))
	if v, ok := c.GetInt64Ok("int8"); v != int64(1) || !ok {
		t.Fail()
	}
	if v, ok := c.GetInt64Ok("uint8"); v != int64(1) || !ok {
		t.Fail()
	}
	if v, ok := c.GetInt64Ok("int16"); v != int64(1) || !ok {
		t.Fail()
	}
	if v, ok := c.GetInt64Ok("uint16"); v != int64(1) || !ok {
		t.Fail()
	}
	if v, ok := c.GetInt64Ok("int32"); v != int64(1) || !ok {
		t.Fail()
	}
	if v, ok := c.GetInt64Ok("uint32"); v != int64(1) || !ok {
		t.Fail()
	}
	if v, ok := c.GetInt64Ok("int"); v != int64(1) || !ok {
		t.Fail()
	}
	if v, ok := c.GetInt64Ok("uint"); v != int64(1) || !ok {
		t.Fail()
	}
	if v, ok := c.GetInt64Ok("no"); v != int64(0) || ok {
		t.Fail()
	}
	if v, ok := c.GetInt64Ok("zero"); v != int64(0) || !ok {
		t.Fail()
	}
}

func TestConfig_GetString(t *testing.T) {
	c := New()
	c.Put("yes", "yes")
	if c.GetString("yes") != "yes" {
		t.Fail()
	}
	if c.GetString("no") != "" {
		t.Fail()
	}
}

func TestConfig_GetStringOk(t *testing.T) {
	c := New()
	c.Put("yes", "yes")
	c.Put("zero", "zero")
	if v, ok := c.GetStringOk("yes"); v != "yes" || !ok {
		t.Fail()
	}
	if v, ok := c.GetStringOk("no"); v != "" || ok {
		t.Fail()
	}
	if v, ok := c.GetStringOk("zero"); v != "zero" || !ok {
		t.Fail()
	}
}

func TestConfig_GetOk(t *testing.T) {
	c := New()
	c.Put("empty", "")
	c.Put("bool", false)
	c.Put("one", "one")
	c.Put("two.three", "three")
	c.Put("four.five.six", "six")
	c.Put("four.five.seven", 7)
	tests := []struct {
		key   string
		value interface{}
		ok    bool
	}{
		{"", nil, false},
		{"empty", "", true},
		{"bool", false, true},
		{"one", "one", true},
		{"One", nil, false},
		{"ten", nil, false},
		{"two", Map(map[string]interface{}{"three": "three"}), true},
		{"four.five.six", "six", true},
		{"four.five.seven", 7, true},
		{"four.ten", nil, false},
		{"four.five.ten", nil, false},
		{"four.five.six.todeep", nil, false},
		{"four.five.six.todeep.reallydeep", nil, false},
		{"four.five", Map(map[string]interface{}{"six": "six", "seven": 7}), true},
	}
	for _, test := range tests {
		value, ok := c.GetOk(test.key)
		if !reflect.DeepEqual(value, test.value) || ok != test.ok {
			t.Errorf("c.GetOk(%v) = %v, %v WANT %v, %v", test.key, value, ok, test.value, test.ok)
		}
	}
}

func TestConfig_GetKeyOk_returnsNilForEmptyKeys(t *testing.T) {
	c := New()
	if v, ok := c.GetKeyOk(nil); v != nil || ok {
		t.Fail()
	}
	if v, ok := c.GetKeyOk(Key([]string{})); v != nil || ok {
		t.Fail()
	}
}

func TestConfig_Put_changed(t *testing.T) {
	c := New()
	changed := c.Put("", "new")
	if !changed {
		t.Fail()
	}
}

func TestConfig_Put_unchanged(t *testing.T) {
	c := New()
	c.Put("one", 1)
	changed := c.Put("one", 1)
	if changed {
		t.Fail()
	}
}

func TestConfig_Put_overwrite(t *testing.T) {
	c := New()
	c.Put("one.two", "two")
	changed := c.Put("one.two.three", "three")
	if !changed {
		t.Fail()
	}
	if _, ok := c.m["one"].(Map)["two"].(Map); !ok {
		t.Fail()
	}
}

func TestConfig_Put_underwrite(t *testing.T) {
	c := New()
	c.Put("one.two.three", NewMap())
	changed := c.Put("one.two", "two")
	if !changed {
		t.Fail()
	}
	if c.m["one"].(Map)["two"] != "two" {
		t.Fail()
	}
}

func TestConfig_PutKey_emptyKey(t *testing.T) {
	c := New()
	changed := c.PutKey(nil, "")
	if changed {
		t.Fail()
	}

	changed = c.PutKey(Key([]string{}), "")
	if changed {
		t.Fail()
	}
}
