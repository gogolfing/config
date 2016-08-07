package config

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c.KeyParser == nil ||
		c.values == nil ||
		c.lock == nil ||
		c.loaders == nil {
		t.Fail()
	}
}

func TestConfig_AddLoader(t *testing.T) {
	c := New()

	result := c.AddLoaders(intLoader(1)).AddLoaders(errorLoader("error"))

	if result != c || !reflect.DeepEqual(c.loaders, []Loader{intLoader(1), errorLoader("error")}) {
		t.Fail()
	}
}

func TestConfig_LoadAll_success(t *testing.T) {
	c := New()
	c.AddLoaders(intLoader(2))

	_, err := c.LoadAll()
	if err != nil {
		t.Fail()
	}

	want := NewValues()
	want.Put(NewKey("2"), 2)

	if !reflect.DeepEqual(c.values, want) {
		t.Fail()
	}
}

func TestConfig_LoadAll_error(t *testing.T) {
	c := New()
	c.AddLoaders(intLoader(2), errorLoader("error loading"))

	_, err := c.LoadAll()
	if err.Error() != "error loading" {
		t.Fail()
	}
}

func TestConfig_Values(t *testing.T) {
	c := New()
	if c.Values() != c.values {
		t.Fail()
	}
}

func TestConfig_EqualValues(t *testing.T) {
	c := New()
	if !c.EqualValues(New()) {
		t.Fail()
	}
}

func TestConfig_Clone(t *testing.T) {
	c := New()
	c.AddLoaders(intLoader(1))

	result := c.Clone()

	if !reflect.DeepEqual(c, result) {
		t.Fail()
	}
}

func TestConfig_GetInt64(t *testing.T) {
	c := New()
	c.Put("int64", 8)
	if c.GetInt64("int64") != int64(8) {
		t.Fail()
	}
	if c.GetInt64("") != int64(0) {
		t.Fail()
	}
}

func TestConfig_GetInt64Ok(t *testing.T) {
	c := New()
	c.Put("zero", 0)
	c.Put("zero string", "0")
	c.Put("uint8", uint8(math.MaxUint8))
	c.Put("int8", int8(math.MaxInt8))
	c.Put("uint16", uint16(math.MaxUint16))
	c.Put("int16", int16(math.MaxInt16))
	c.Put("uint32", uint32(math.MaxUint32))
	c.Put("int32", int32(math.MaxInt32))
	c.Put("uint", uint(math.MaxUint32))
	c.Put("int", int(math.MaxInt32))
	c.Put("uint64", uint64(math.MaxUint64))
	c.Put("int64", int64(math.MaxInt64))
	tests := []struct {
		key    string
		result int64
		ok     bool
	}{
		{"something", 0, false},
		{"zero", 0, true},
		{"zero string", 0, false},
		{"uint8", math.MaxUint8, true},
		{"int8", math.MaxInt8, true},
		{"uint16", math.MaxUint16, true},
		{"int16", math.MaxInt16, true},
		{"uint32", math.MaxUint32, true},
		{"int32", math.MaxInt32, true},
		{"int", math.MaxInt32, true},
		{"uint", math.MaxUint32, true},
		{"uint64", -1, true},
		{"int64", math.MaxInt64, true},
	}
	for _, test := range tests {
		result, ok := c.GetInt64Ok(test.key)
		if result != test.result || ok != test.ok {
			t.Errorf("c.GetInt64Ok(%v) = %v, %v WANT %v, %v", test.key, result, ok, test.result, test.ok)
		}
	}
}

func TestConfig_GetBool(t *testing.T) {
	c := New()
	c.Put("true", true)
	c.Put("false", false)
	if c.GetBool("true") != true {
		t.Fail()
	}
	if c.GetBool("false") != false {
		t.Fail()
	}
	if c.GetBool("") != false {
		t.Fail()
	}
}

func TestConfig_GetBoolOk(t *testing.T) {
	c := New()
	c.Put("true", true)
	c.Put("false", false)
	c.Put("true string", "true")
	if b, ok := c.GetBoolOk("true"); b != true || !ok {
		t.Fail()
	}
	if b, ok := c.GetBoolOk("false"); b != false || !ok {
		t.Fail()
	}
	if b, ok := c.GetBoolOk("true string"); b != false || ok {
		t.Fail()
	}
}

func TestConfig_GetString(t *testing.T) {
	c := New()
	c.Put("a", "a")
	if c.GetString("a") != "a" {
		t.Fail()
	}
	if c.GetString("") != "" {
		t.Fail()
	}
}

func TestConfig_GetStringOk(t *testing.T) {
	c := New()
	c.Put("a", "a")
	c.Put("empty", "")
	if s, ok := c.GetStringOk("a"); s != "a" || !ok {
		t.Fail()
	}
	if s, ok := c.GetStringOk("empty"); s != "" || !ok {
		t.Fail()
	}
	if s, ok := c.GetStringOk(""); s != "" || ok {
		t.Fail()
	}
}

func TestConfig_GetFloat64(t *testing.T) {
	c := New()
	c.Put("float32", float32(32.0))
	c.Put("float64", float64(64.0))
	if c.GetFloat64("float32") != float64(32.0) {
		t.Fail()
	}
	if c.GetFloat64("float64") != float64(64.0) {
		t.Fail()
	}
	if c.GetFloat64("") != float64(0) {
		t.Fail()
	}
}

func TestConfig_GetFloat64Ok(t *testing.T) {
	c := New()
	c.Put("float32", float32(32.0))
	c.Put("float64", float64(64.0))
	c.Put("empty", float64(0))
	c.Put("float64 string", "1234.1234")
	if f, ok := c.GetFloat64Ok("float32"); f != float64(32.0) || !ok {
		t.Fail()
	}
	if f, ok := c.GetFloat64Ok("float64"); f != float64(64.0) || !ok {
		t.Fail()
	}
	if f, ok := c.GetFloat64Ok("empty"); f != float64(0) || !ok {
		t.Fail()
	}
	if f, ok := c.GetFloat64Ok("float64 string"); f != float64(0) || ok {
		t.Fail()
	}
	if f, ok := c.GetFloat64Ok(""); f != float64(0) || ok {
		t.Log("float64")
		t.Fail()
	}
}

func TestConfig_GetValues(t *testing.T) {
	c := New()
	c.Put("a.b", "b")
	c.Put("a.c", "c")

	aValues := NewValues()
	aValues.Put(NewKey("b"), "b")
	aValues.Put(NewKey("c"), "c")
	if !reflect.DeepEqual(c.GetValues("a"), aValues) {
		t.Fail()
	}

	if c.GetValues("a.b") != nil {
		t.Fail()
	}
	if c.GetValues("") != nil {
		t.Fail()
	}
}

func TestConfig_GetValuesOk(t *testing.T) {
	c := New()
	c.Put("a.b", "b")
	c.Put("a.c", "c")

	aValues := NewValues()
	aValues.Put(NewKey("b"), "b")
	aValues.Put(NewKey("c"), "c")
	if v, ok := c.GetValuesOk("a"); !reflect.DeepEqual(v, aValues) || !ok {
		t.Fail()
	}

	if v, ok := c.GetValuesOk("a.b"); v != nil || ok {
		t.Fail()
	}
	if v, ok := c.GetValuesOk(""); v != nil || ok {
		t.Fail()
	}
}

func TestConfig_Get(t *testing.T) {
	c := getFullConfig()
	tests := getFullConfigTests()
	for _, test := range tests {
		keyString := strings.Join(test.key, string(c.KeyParser.(SeparatorKeyParser)))
		value := c.Get(keyString)
		if !reflect.DeepEqual(value, test.value) {
			t.Errorf("c.Get(%v) = %v WANT %v", test.key, value, test.value)
		}
	}
}

func TestConfig_GetOk(t *testing.T) {
	c := getFullConfig()
	tests := getFullConfigTests()
	for _, test := range tests {
		keyString := strings.Join(test.key, string(c.KeyParser.(SeparatorKeyParser)))
		value, ok := c.GetOk(keyString)
		if !reflect.DeepEqual(value, test.value) || ok != test.ok {
			t.Errorf("c.GetOk(%v) = %v, %v WANT %v, %v", test.key, value, ok, test.value, test.ok)
		}
	}
}

func TestConfig_GetKey(t *testing.T) {
	c := getFullConfig()
	tests := getFullConfigTests()
	for _, test := range tests {
		value := c.GetKey(test.key)
		if !reflect.DeepEqual(value, test.value) {
			t.Errorf("c.GetKey(%v) = %v WANT %v", test.key, value, test.value)
		}
	}
}

func TestConfig_GetKeyOk(t *testing.T) {
	c := getFullConfig()
	tests := getFullConfigTests()
	for _, test := range tests {
		value, ok := c.GetKeyOk(test.key)
		if !reflect.DeepEqual(value, test.value) || ok != test.ok {
			t.Errorf("c.GetKeyOk(%v) = %v, %v WANT %v, %v", test.key, value, ok, test.value, test.ok)
		}
	}
}

func TestConfig_Merge(t *testing.T) {
	c := New()
	c.Merge(New())
}

func getFullConfig() *Config {
	c := New()

	c.Put("bools.true", true)
	c.Put("bools.false", false)
	c.Put("bools.zero", false)

	c.Put("ints.ten", int(10))
	c.Put("ints.negativeOne", int(-1))
	c.Put("ints.zero", int(0))

	c.Put("int64s.ten", int64(10))
	c.Put("int64s.negativeOne", int64(-1))
	c.Put("int64s.zero", int64(0))

	c.Put("strings.foo", "foo")
	c.Put("strings.bar", "bar")
	c.Put("strings.zero", "")

	c.Put("float64s.ten", float64(10.0))
	c.Put("float64s.negativeOne", float64(-1.0))
	c.Put("float64s.zero", float64(0.0))

	c.Put("values.one.a", "a")
	c.Put("values.one.b", "b")

	c.Put("values.two.a", "a")
	c.Put("values.two.b", "b")

	return c
}

func getFullConfigTests() []*getFullConfigTest {
	abValues := NewValues()
	abValues.Put(NewKey("a"), "a")
	abValues.Put(NewKey("b"), "b")

	return []*getFullConfigTest{
		{NewKey("bools", "true"), true, true},
		{NewKey("bools", "false"), false, true},
		{NewKey("bools", "true"), true, true},
		{NewKey("bools", ""), nil, false},

		{NewKey("ints", "ten"), int(10), true},
		{NewKey("ints", "negativeOne"), int(-1), true},
		{NewKey("ints", "zero"), int(0), true},
		{NewKey("ints", ""), nil, false},

		{NewKey("int64s", "ten"), int64(10), true},
		{NewKey("int64s", "negativeOne"), int64(-1), true},
		{NewKey("int64s", "zero"), int64(0), true},
		{NewKey("int64s", ""), nil, false},

		{NewKey("strings", "foo"), "foo", true},
		{NewKey("strings", "bar"), "bar", true},
		{NewKey("strings", "zero"), "", true},
		{NewKey("strings", ""), nil, false},

		{NewKey("float64s", "ten"), float64(10.0), true},
		{NewKey("float64s", "negativeOne"), float64(-1.0), true},
		{NewKey("float64s", "zero"), float64(0.0), true},
		{NewKey("float64s", ""), nil, false},

		{NewKey("values", "one"), abValues, true},
		{NewKey("values", "one", ""), nil, false},

		{NewKey("values", "two"), abValues, true},
		{NewKey("values", "two", ""), nil, false},

		{NewKey("values", "one", "A"), nil, false},
		{NewKey("values", "two", "B"), nil, false},
	}
}

func TestConfig_Put(t *testing.T) {
	c := New()

	tests := []struct {
		keyString string
		value     interface{}
		changed   bool
	}{
		{"a.b.c", "c", true},
		{"a.b.cc", "c", true},
		{"a.b.cc", "cc", true},
		{"a.b.c.d", "d", true},
		{"a.b.c.d", 10, true},
		{"a.b.c.d", 10, false},
		{"a.b.c.d", true, true},
		{"a.b.c.d", false, true},
		{"a", "something not a values", true},
	}

	for _, test := range tests {
		changed := c.Put(test.keyString, test.value)
		if changed != test.changed {
			t.Errorf("c.Put(%v) = %v WANT %v", test.keyString, changed, test.changed)
		}
	}
}

func TestConfig_MergeLoaders_success(t *testing.T) {
	c := New()
	_, err := c.MergeLoaders(intLoader(1), intLoader(2))
	if err != nil {
		t.Error(err)
	}

	result := New()
	result.Put("1", 1)
	result.Put("2", 2)

	if !c.EqualValues(result) {
		t.Fail()
	}
}

func TestConfig_MergeLoaders_error(t *testing.T) {
	c := New()
	_, err := c.MergeLoaders(intLoader(1), errorLoader("error"), intLoader(2))
	if err == nil {
		t.Fail()
	}

	result := New()

	if !c.EqualValues(result) {
		t.Fail()
	}
}

func TestConfig_Remove(t *testing.T) {
	c := New()
	c.Put("a", "a")
	c.Remove("a")
	if result, ok := c.GetOk("a"); result != nil || ok {
		t.Fail()
	}
}

func TestConfig_NewKey(t *testing.T) {
	c := New()
	called := false
	c.KeyParser = KeyParserFunc(func(k string) Key {
		if k != "key" {
			t.Fail()
		}
		called = true
		return NewKey("this", "is", "a", "parsed", "key")
	})
	key := c.NewKey("key")
	if !key.Equal(NewKey("this", "is", "a", "parsed", "key")) {
		t.Fail()
	}
	if !called {
		t.Fail()
	}
}

type getFullConfigTest struct {
	key   Key
	value interface{}
	ok    bool
}

type intLoader int

func (l intLoader) Load() (*Values, error) {
	v := NewValues()
	v.Put(NewKey(fmt.Sprint(int(l))), int(l))
	return v, nil
}

type errorLoader string

func (l errorLoader) Load() (*Values, error) {
	return nil, fmt.Errorf("%v", string(l))
}
