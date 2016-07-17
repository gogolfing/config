package config

import (
	"io"
	"os"
)

//Loader defines an entity that can generate a new Values instance.
type Loader interface {
	Load() (*Values, error)
}

//ReaderFuncLoader is a func definition that takes in an io.Reader and returns
//a new Values instance and possible error.
type ReaderFuncLoader func(io.Reader) (*Values, error)

type fileFuncLoader struct {
	rfl   ReaderFuncLoader
	paths []string
}

//NewFileFuncLoader creates a Loader that uses rfl to load and merge Values from
//from each file existing at each path in paths.
//If rfl returns an error for any path in paths then that error is immediately
//returned from Loader.Load() and Values will be nil.
func NewFileFuncLoader(rfl ReaderFuncLoader, paths ...string) Loader {
	return &fileFuncLoader{
		rfl:   rfl,
		paths: paths,
	}
}

func (l *fileFuncLoader) Load() (*Values, error) {
	values := NewValues()
	for _, path := range l.paths {
		temp, err := l.loadPath(path)
		if err != nil {
			return nil, err
		}
		values.Merge(NewKey(), temp)
	}
	return values, nil
}

func (l *fileFuncLoader) loadPath(path string) (*Values, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	values, err := l.rfl(file)
	if err != nil {
		return nil, err
	}
	return values, file.Close()
}

type readerFuncLoader struct {
	rfl ReaderFuncLoader
	r   io.Reader
}

//NewReaderFuncLoader creates a Loader that return rfl(r) in its Load() method.
func NewReaderFuncLoader(rfl ReaderFuncLoader, r io.Reader) Loader {
	return &readerFuncLoader{
		rfl: rfl,
		r:   r,
	}
}

func (l *readerFuncLoader) Load() (*Values, error) {
	return l.rfl(l.r)
}
