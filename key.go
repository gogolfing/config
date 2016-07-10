package config

import "strings"

type Key []string

func NewKey(parts ...string) Key {
	return Key(parts)
}

func NewKeySep(source, sep string) Key {
	return NewKey(strings.Split(source, sep)...)
}

func (k Key) IsEmpty() bool {
	return k.Len() == 0
}

func (k Key) Len() int {
	return len(k)
}

func (k Key) Equal(other Key) bool {
	if k.IsEmpty() && other.IsEmpty() {
		return true
	}
	if len(k) != len(other) {
		return false
	}
	for i, part := range k {
		if part != other[i] {
			return false
		}
	}
	return true
}

func (k Key) StartsWith(other Key) bool {
	if other.Len() > k.Len() {
		return false
	}
	for i, part := range other {
		if k[i] != part {
			return false
		}
	}
	return true
}

func (k Key) EndsWith(other Key) bool {
	if other.Len() > k.Len() {
		return false
	}
	for i, _ := range other {
		part := other[other.Len()-1-i]
		if k[k.Len()-1-i] != part {
			return false
		}
	}
	return true
}

func (k Key) Append(others ...Key) Key {
	result := NewKey(k...)
	for _, other := range others {
		result = append(result, other...)
	}
	return result
}

func (k Key) AppendStrings(others ...string) Key {
	return k.Append(NewKey(others...))
}

type KeyParser interface {
	Parse(k string) Key
}

type KeyParserFunc func(k string) Key

func (pf KeyParserFunc) Parse(k string) Key {
	return pf(k)
}

type SeparatorKeyParser string

func (p SeparatorKeyParser) Parse(k string) Key {
	return NewKey(strings.Split(k, string(p))...)
}

const PeriodSeparatorKeyParser = SeparatorKeyParser(".")
