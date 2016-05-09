package config

import "strings"

type Key []string

func NewKey(source, sep string) Key {
	return Key(strings.Split(source, sep))
}
