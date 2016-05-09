package config

import "testing"

func TestNew(t *testing.T) {
	c := New()
	if c.Separator != "." {
		t.Fail()
	}
	if c.lock == nil || c.set == nil {
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
	if _, ok := c.set["one"].(Set)["two"].(Set); !ok {
		t.Fail()
	}
}

func TestConfig_Put_underwrite(t *testing.T) {
	c := New()
	c.Put("one.two.three", NewSet())
	changed := c.Put("one.two", "two")
	if !changed {
		t.Fail()
	}
	if c.set["one"].(Set)["two"] != "two" {
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
