package config

import "io"

type Parser interface {
	Parse(r io.Reader) ([]KeyValue, error)
}

type ParserFunc func(io.Reader) ([]KeyValue, error)

func (p ParserFunc) Parse(r io.Reader) ([]KeyValue, error) {
	return p(r)
}
