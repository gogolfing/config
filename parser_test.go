package config

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestParserFunc_Parse(t *testing.T) {
	reader := strings.NewReader("hello")
	called := false
	pf := ParserFunc(func(r io.Reader) ([]KeyValue, error) {
		if r != reader {
			t.Fail()
		}
		called = true
		return []KeyValue{{Key: NewKey("key", "."), Value: "value"}}, fmt.Errorf("error")
	})
	pf.Parse(reader)
	if !called {
		t.Fail()
	}
}
