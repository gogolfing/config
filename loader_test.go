package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNewFileFuncLoader(t *testing.T) {
	l := NewFileFuncLoader("path", func(_ io.Reader) (*Values, error) { return nil, nil })
	if l == nil {
		t.Fail()
	}
	ffl, ok := l.(*fileFuncLoader)
	if !ok {
		t.Fail()
	}
	if ffl.path != "path" || ffl.rfl == nil {
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

	l := NewFileFuncLoader(file.Name(), func(r io.Reader) (*Values, error) {
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
	})

	v, err := l.Load()

	want := NewValues()
	want.Put(nil, "foobar")

	if !v.Equal(want) || err != nil {
		t.Fail()
	}
}

func TestNewReaderFuncLoader(t *testing.T) {
	l := NewReaderFuncLoader(strings.NewReader("foobar"), func(_ io.Reader) (*Values, error) { return nil, nil })
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
	l := NewReaderFuncLoader(strings.NewReader("foobar"), func(r io.Reader) (*Values, error) {
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
	})

	v, err := l.Load()

	want := NewValues()
	want.Put(nil, "foobar")

	if !v.Equal(want) || err != nil {
		t.Fail()
	}
}
