package json

import (
	"testing"

	"github.com/gogolfing/config"
)

func TestLoader_LoadString_notAnObject(t *testing.T) {
	in := "foobar"

	v, err := (&Loader{}).LoadString(in)

	if v != nil || err == nil {
		t.Fail()
	}
}

func TestLoader_LoadString_keyPrefix(t *testing.T) {
	in := `{
		"a": { "b": "b" },
		"c": { "d": "d" }
	}`
	l := &Loader{
		KeyPrefix: config.NewKey("a"),
	}
	want := config.NewValues()
	want.Put(config.NewKey("a", "b"), "b")

	testLoadStringWithWantedValues(t, l, in, want)
}

func TestLoader_LoadString_keySuffix(t *testing.T) {
	in := `{
		"a": { "b": "b" },
		"c": { "d": "d" }
	}`
	l := &Loader{
		KeySuffix: config.NewKey("d"),
	}
	want := config.NewValues()
	want.Put(config.NewKey("c", "d"), "d")

	testLoadStringWithWantedValues(t, l, in, want)
}

func TestLoader_LoadString_discardNull(t *testing.T) {
	in := `{
		"a": null
	}`
	l := &Loader{
		DiscardNull: true,
	}
	want := config.NewValues()

	testLoadStringWithWantedValues(t, l, in, want)
}

func TestLoader_LoadString_numberAsString(t *testing.T) {
	in := `{
		"a": 12
	}`
	l := &Loader{
		NumberAsString: true,
	}
	want := config.NewValues()
	want.Put(config.NewKey("a"), "12")

	testLoadStringWithWantedValues(t, l, in, want)
}

func testLoadStringWithWantedValues(t *testing.T, l *Loader, in string, want *config.Values) {
	v, err := l.LoadString(in)
	if err != nil {
		t.Error(err)
	}
	if !v.Equal(want) {
		t.Fail()
	}
}
