package config

import "io"

type Parser interface {
	Parse(r io.Reader) ([]KeyValue, error)
}
