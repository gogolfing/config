package config

import (
	"fmt"
	"testing"
)

func TestPackageDocConfig(t *testing.T) {
	c := New()
	c.Put("foo", "bar") //true
	c.Put("foo", "bar") //false
	c.Put("foo", 1024)  //true
}

type sliceLoader []interface{}

func (sl sliceLoader) Load() (*Values, error) {
	values := NewValues()
	for i, v := range sl {
		values.Put(NewKey("slice", fmt.Sprint(i)), v)
	}
	return values, nil
}

func TestPackageDocConfigLoader(t *testing.T) {
	loader := sliceLoader([]interface{}{"hello", "world", 234, true})
	c := New()
	if changed, err := c.PutLoaders(loader); !changed || err != nil {
		t.Fail()
	}

	if v := c.Get("slice.0"); v != "hello" {
		t.Fail()
	}
	c.Get("slice.1") //world
	c.Get("slice.2") //234
	c.Get("slice.3") //true
}
