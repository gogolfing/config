package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewFileFuncLoader(t *testing.T) {
	l := NewFileFuncLoader(func(_ io.Reader) (*Values, error) { return nil, nil }, "path")
	if l == nil {
		t.Fail()
	}
	ffl, ok := l.(*fileFuncLoader)
	if !ok {
		t.Fail()
	}
	if ffl.rfl == nil || !reflect.DeepEqual(ffl.paths, []string{"path"}) {
		t.Fail()
	}
}

func TestFileFuncLoader_Load(t *testing.T) {
	file, err := ioutil.TempFile("", "gogolfing.config")
	if err != nil {
		t.Fatal(err)
	}
	_, err = fmt.Fprint(file, "foobar")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	defer os.Remove(file.Name())

	l := NewFileFuncLoader(
		func(r io.Reader) (*Values, error) {
			bytes, err := ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}
			if string(bytes) != "foobar" {
				t.Fail()
			}

			v := NewValues()
			v.Put(nil, "foobar")
			return v, nil
		},
		file.Name(),
	)

	v, err := l.Load()

	want := NewValues()
	want.Put(nil, "foobar")

	if !v.Equal(want) || err != nil {
		t.Fail()
	}
}

func TestNewReaderFuncLoader(t *testing.T) {
	l := NewReaderFuncLoader(func(_ io.Reader) (*Values, error) { return nil, nil }, strings.NewReader("foobar"))
	if l == nil {
		t.Fail()
	}
	rfl, ok := l.(*readerFuncLoader)
	if !ok {
		t.Fail()
	}
	if rfl.r == nil || rfl.rfl == nil {
		t.Fail()
	}
}

func TestReaderFuncLoader_Load(t *testing.T) {
	l := NewReaderFuncLoader(
		func(r io.Reader) (*Values, error) {
			bytes, err := ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}
			if string(bytes) != "foobar" {
				t.Fail()
			}

			v := NewValues()
			v.Put(nil, "foobar")
			return v, nil
		},
		strings.NewReader("foobar"),
	)

	v, err := l.Load()

	want := NewValues()
	want.Put(nil, "foobar")

	if !v.Equal(want) || err != nil {
		t.Fail()
	}
}
