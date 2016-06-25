package config

import "strings"

type Key []string

func NewKey(parts ...string) Key {
	return Key(parts)
}

func NewKeySep(source, sep string) Key {
	return Key(strings.Split(source, sep))
}

func (k Key) IsEmpty() bool {
	return len(k) == 0
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

type Converter interface {
	Convert(value interface{}) interface{}
}

type ConverterFunc func(interface{}) interface{}

func (c ConverterFunc) Convert(value interface{}) interface{} {
	return c(value)
}

var DefaultConverter Converter = ConverterFunc(Convert)

func Convert(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		return ConvertBuiltinMap(v)
	case []interface{}:
		return ConvertSlice(v)
	case KeyValue:
		return ConvertKeyValues([]KeyValue{v})
	case []KeyValue:
		return ConvertKeyValues(v)
	}
	return value
}

func ConvertBuiltinMap(m map[string]interface{}) Map {
	result := NewMap()
	for key, value := range m {
		result[key] = Convert(value)
	}
	return result
}

func ConvertSlice(slice []interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = Convert(v)
	}
	return result
}

func ConvertKeyValues(keyValues []KeyValue) Map {
	result := NewMap()
	for _, kv := range keyValues {
		result.Put(kv.Key, Convert(kv.Value))
	}
	return result
}
