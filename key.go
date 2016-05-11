package config

import "strings"

type Key []string

func NewKey(source, sep string) Key {
	return Key(strings.Split(source, sep))
}

type KeyValue struct {
	Key
	Value interface{}
}

func NewKeyValue(key Key, value interface{}) KeyValue {
	return KeyValue{
		Key:   key,
		Value: value,
	}
}
