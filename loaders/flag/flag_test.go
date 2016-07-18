package flag

import (
	"io/ioutil"
	"testing"

	"github.com/gogolfing/config"
)

func TestLoader_Load(t *testing.T) {
	l := New("")
	l.FlagSet.String("a", "", "")
	l.Args = []string{"-a", "value"}

	want := config.NewValues()
	want.Put(config.NewKey("a"), "value")

	testLoadWithWantedValues(t, l, want)
}

func TestLoader_Load_idempotent(t *testing.T) {
	l := New("")
	l.FlagSet.String("a", "", "")
	l.Args = []string{"-a", "value"}

	want := config.NewValues()
	want.Put(config.NewKey("a"), "value")

	testLoadWithWantedValues(t, l, want)
	testLoadWithWantedValues(t, l, want)
}

func TestLoader_Load_error(t *testing.T) {
	l := New("")
	l.Args = []string{"-a", "value"}
	l.FlagSet.SetOutput(ioutil.Discard)

	v, err := l.Load()
	if v != nil || err == nil {
		t.Fail()
	}
}

func testLoadWithWantedValues(t *testing.T, l *Loader, want *config.Values) {
	v, err := l.Load()
	if err != nil {
		t.Error(err)
	}
	if !v.Equal(want) {
		t.Fail()
	}
}
