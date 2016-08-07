package config

import "strings"

//Key is the entity that allows access to values stored within a Values instance.
type Key []string

//NewKey creates a Key with all strings in parts in the returned Key.
//It essentially casts the string slice to a Key.
func NewKey(parts ...string) Key {
	return Key(parts)
}

//NewKeySep returns a Key that is the result of strings.Split(source, sep).
func NewKeySep(source, sep string) Key {
	return NewKey(strings.Split(source, sep)...)
}

//IsEmpty determines whether or not the length of k is 0.
func (k Key) IsEmpty() bool {
	return k.Len() == 0
}

//Len returns the length of k.
func (k Key) Len() int {
	return len(k)
}

//Equal determines whether or not k and other are the same length and all individual
//strings are identical at their respective indices.
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

//StartsWith determines whether or not k is at least the same length as other
//and all strings in other appear at the first consecutive indices of k.
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

//EndsWith determines whether or not k is at least the same length as other
//and all strings in other appear at the last consecutive indices of k.
func (k Key) EndsWith(other Key) bool {
	if other.Len() > k.Len() {
		return false
	}
	for i := range other {
		part := other[other.Len()-1-i]
		if k[k.Len()-1-i] != part {
			return false
		}
	}
	return true
}

//Append returns a new Key with all strings from k and other.
func (k Key) Append(others ...Key) Key {
	result := NewKey(k...)
	for _, other := range others {
		result = append(result, other...)
	}
	return result
}

//AppendStrings returns a new Key with all strings from k and others.
func (k Key) AppendStrings(others ...string) Key {
	return k.Append(NewKey(others...))
}

//KeyParser defines an entity that can parse a string and turn it into a Key.
type KeyParser interface {
	Parse(k string) Key
}

//KeyParserFunc is a func implementation of KeyParser that takes in a single string
//and returns a Key.
type KeyParserFunc func(k string) Key

//Parse simply calls pf(k).
func (pf KeyParserFunc) Parse(k string) Key {
	return pf(k)
}

//SeparatorKeyParser is a KeyParser that creates Keys from the result of calling
//strings.Split() with k and string(SeparatorKeyParser).
type SeparatorKeyParser string

//Parse returns NewKeySep(k, string(p)).
func (p SeparatorKeyParser) Parse(k string) Key {
	return NewKeySep(k, string(p))
}

//PeriodSeparatorKeyParser is the default KeyParser set to c.KeyParser in New().
//See SeparatorKeyParser.
const PeriodSeparatorKeyParser = SeparatorKeyParser(".")
