package config

import (
	"io"
	"os"
)

type Loader interface {
	Load() (*Values, error)
}

type ReaderFuncLoader func(io.Reader) (*Values, error)

type fileFuncLoader struct {
	path string
	rfl  ReaderFuncLoader
}

func NewFileFuncLoader(path string, rfl ReaderFuncLoader) Loader {
	return &fileFuncLoader{
		path: path,
		rfl:  rfl,
	}
}

func (l *fileFuncLoader) Load() (*Values, error) {
	file, err := os.Open(l.path)
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
	r   io.Reader
	rfl ReaderFuncLoader
}

func NewReaderFuncLoader(r io.Reader, rfl ReaderFuncLoader) Loader {
	return &readerFuncLoader{
		r:   r,
		rfl: rfl,
	}
}

func (l *readerFuncLoader) Load() (*Values, error) {
	return l.rfl(l.r)
}
