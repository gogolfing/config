package config

type Set map[string]interface{}

func NewSet() Set {
	return Set(map[string]interface{}{})
}

func (s Set) clone() Set {
	result := NewSet()
	for key, value := range s {
		valueSet, ok := value.(Set)
		if ok {
			value = valueSet.clone()
		}
		result[key] = value
	}
	return result
}
