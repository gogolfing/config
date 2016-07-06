package config

import "io"

type Loader interface {
	Load() (*Values, error)
}

type readCloserFuncLoader struct {
	r io.Reader
	f func(io.Reader) (*Values, error)
}

func NewReaderFuncLoader(r io.Reader, f func(io.Reader) (*Values, error)) Loader {
	return &readCloserFuncLoader{
		r: r,
		f: f,
	}
}

func (l *readCloserFuncLoader) Load() (*Values, error) {
	return l.f(l.r)
}
