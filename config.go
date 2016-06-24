package config

import "sync"

type Config struct {
	Separator string

	lock *sync.RWMutex
	m    Map
	// loaders []Loader
}

func New() *Config {
	return &Config{
		Separator: ".",

		lock: &sync.RWMutex{},
		m:    NewMap(),
		// loaders: []Loader{},
	}
}

/*
func (c *Config) AddLoader(loader Loader) {
	c.loaders = append(c.loaders, loader)
}
*/

/*
func (c *Config) LoadAll() error {
	return c.PutLoaders(c.loaders...)
}
*/

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
	value, _ := c.GetOk(key)
	switch v := value.(type) {
	case int8:
		return int(v), true
	case uint8:
		return int(v), true
	case int16:
		return int(v), true
	case uint16:
		return int(v), true
	case int32:
		return int(v), true
	case int:
		return v, true
	}
	return 0, false
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
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.m.GetOk(key)
}

/*
func (c *Config) PutLoaders(loaders ...Loader) error {
	result := []KeyValue{}
	for _, loader := range loaders {
		keyValues, err := loader.Load()
		if err != nil {
			return err
		}
		for _, keyValue := range keyValues {
			result = append(result, keyValue)
		}
	}
	for _, keyValue := range result {
		c.PutKey(keyValue.Key, keyValue.Value)
	}
	return nil
}
*/

func (c *Config) Put(key string, value interface{}) bool {
	return c.PutKey(c.NewKey(key), value)
}

func (c *Config) PutKey(key Key, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.m.Put(c.convertKeyValue(key, value))
}

func (c *Config) convertKeyValue(key Key, value interface{}) (Key, interface{}) {
	if len(key) == 0 {
		return key, value
	}
	return key, value
}

func (c *Config) NewKey(source string) Key {
	return NewKey(source, c.Separator)
}
