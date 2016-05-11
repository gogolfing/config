package config

import "io"

type Loader interface {
	Load() ([]KeyValue, error)
}

type ReaderParserLoader struct {
	io.Reader
	Parser
}

func (r ReaderParserLoader) Load() ([]KeyValue, error) {
	return r.Parser.Parse(r.Reader)
}
