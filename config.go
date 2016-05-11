package config

import "sync"

type Config struct {
	Separator string

	lock    *sync.RWMutex
	set     Set
	loaders []Loader
}

func New() *Config {
	return &Config{
		Separator: ".",

		lock:    &sync.RWMutex{},
		set:     NewSet(),
		loaders: []Loader{},
	}
}

func (c *Config) Get(key string) interface{} {
	v, _ := c.GetOk(key)
	return v
}

func (c *Config) GetBool(key string) bool {
	b, _ := c.GetBoolOk(key)
	return b
}

func (c *Config) GetBoolOk(key string) (bool, bool) {
	v, ok := c.GetOk(key)
	if b, typeOk := v.(bool); ok && typeOk {
		return b, true
	}
	return false, false
}

func (c *Config) GetFloat64(key string) float64 {
	f, _ := c.GetFloat64Ok(key)
	return f
}

func (c *Config) GetFloat64Ok(key string) (float64, bool) {
	v, ok := c.GetOk(key)
	if f, typeOk := v.(float64); ok && typeOk {
		return f, true
	}
	return 0, false
}

func (c *Config) GetInt(key string) int {
	i, _ := c.GetIntOk(key)
	return i
}

func (c *Config) GetIntOk(key string) (int, bool) {
	v, _ := c.GetOk(key)
	result, ok := 0, false
	switch v.(type) {
	case int8:
		result, ok = int(v.(int8)), true
	case uint8:
		result, ok = int(v.(uint8)), true
	case int16:
		result, ok = int(v.(int16)), true
	case uint16:
		result, ok = int(v.(uint16)), true
	case int32:
		result, ok = int(v.(int32)), true
	case int:
		result, ok = v.(int), true
	}
	return result, ok
}

func (c *Config) GetInt64(key string) int64 {
	i, _ := c.GetInt64Ok(key)
	return i
}

func (c *Config) GetInt64Ok(key string) (int64, bool) {
	v, _ := c.GetOk(key)
	result, ok := int64(0), false
	switch v.(type) {
	case int8:
		result, ok = int64(v.(int8)), true
	case uint8:
		result, ok = int64(v.(uint8)), true
	case int16:
		result, ok = int64(v.(int16)), true
	case uint16:
		result, ok = int64(v.(uint16)), true
	case int32:
		result, ok = int64(v.(int32)), true
	case uint32:
		result, ok = int64(v.(uint32)), true
	case int:
		result, ok = int64(v.(int)), true
	case uint:
		result, ok = int64(v.(uint)), true
	case int64:
		result, ok = v.(int64), true
	}
	return result, ok
}

func (c *Config) GetString(key string) string {
	s, _ := c.GetStringOk(key)
	return s
}

func (c *Config) GetStringOk(key string) (string, bool) {
	v, ok := c.GetOk(key)
	if s, typeOk := v.(string); ok && typeOk {
		return s, true
	}
	return "", false
}

func (c *Config) GetOk(key string) (interface{}, bool) {
	return c.GetKeyOk(c.NewKey(key))
}

func (c *Config) GetKeyOk(key Key) (interface{}, bool) {
	if len(key) == 0 {
		return nil, false
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	var value interface{} = c.set
	ok := false
	for _, keyPart := range key {
		_, ok = value.(Set)
		if !ok {
			return nil, false
		}
		value, ok = value.(Set)[keyPart]
	}

	return externalValue(value), ok
}

func externalValue(value interface{}) interface{} {
	if valueSet, ok := value.(Set); ok {
		return valueSet.clone()
	}
	return value
}

func (c *Config) Put(key string, value interface{}) bool {
	return c.PutKey(c.NewKey(key), value)
}

func (c *Config) PutKey(key Key, value interface{}) bool {
	if len(key) == 0 {
		return false
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	lastSet := c.set
	changed := false
	for i := 0; i < len(key)-1; i++ {
		keyPart := key[i]
		tempValue := lastSet[keyPart]
		tempSet, tempSetOk := tempValue.(Set)
		if !tempSetOk {
			tempSet = NewSet()
			lastSet[keyPart] = tempSet
			changed = true
		}
		lastSet = tempSet
	}

	lastPart := key[len(key)-1]
	oldValue := lastSet[lastPart]
	changed = changed || (oldValue != value)
	lastSet[lastPart] = value

	return changed
}

func (c *Config) NewKey(source string) Key {
	return NewKey(source, c.Separator)
}
