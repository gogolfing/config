package config

type Map map[string]interface{}

func NewMap(keyValues ...KeyValue) Map {
	result := Map(map[string]interface{}{})
	for _, keyValue := range keyValues {
		result.Put(keyValue.Key, keyValue.Value)
	}
	return result
}

func (m Map) GetOk(key Key) (interface{}, bool) {
	if len(key) == 0 {
		return nil, false
	}
	var value interface{} = m
	ok := false
	for _, keyPart := range key {
		_, ok = value.(Map)
		if !ok {
			return nil, false
		}
		value, ok = value.(Map)[keyPart]
	}
	return cloneValue(value), ok
}

func (m Map) Put(key Key, value interface{}) bool {
	if len(key) == 0 {
		return false
	}

	lastMap := m
	changed := false
	for i := 0; i < len(key)-1; i++ {
		keyPart := key[i]
		tempValue := lastMap[keyPart]
		tempMap, tempMapOk := tempValue.(Map)
		if !tempMapOk {
			tempMap = NewMap()
			lastMap[keyPart] = tempMap
			changed = true
		}
		lastMap = tempMap
	}

	lastPart := key[len(key)-1]
	oldValue := lastMap[lastPart]
	changed = changed || (oldValue != value)
	lastMap[lastPart] = value

	return changed
}

func (m Map) Clone() Map {
	cloned := NewMap()
	for key, value := range m {
		value = cloneValue(value)
		cloned[key] = value
	}
	return cloned
}

func cloneValue(value interface{}) interface{} {
	switch v := value.(type) {
	case Map:
		return v.Clone()
	case []interface{}:
		return cloneSlice(v)
	}
	return value
}

func cloneSlice(value []interface{}) []interface{} {
	cloned := make([]interface{}, len(value))
	for i, v := range value {
		cloned[i] = cloneValue(v)
	}
	return cloned
}
