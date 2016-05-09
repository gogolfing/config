package config

import "sync"

type Config struct {
	Separator string

	lock *sync.RWMutex
	set  Set
}

func New() *Config {
	return &Config{
		Separator: ".",

		lock: &sync.RWMutex{},
		set:  NewSet(),
	}
}

func (c *Config) Get(key string) interface{} {
	v, _ := c.GetOk(key)
	return v
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
