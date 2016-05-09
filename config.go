package config

import "strings"

type Config struct {
	Separator string

	set
}

func New() *Config {
	return &Config{
		Separator: ".",

		set: newSet(),
	}
}

func (c *Config) Get(key string) interface{} {
	v, _ := c.GetOk(key)
	return v
}

func (c *Config) GetOk(key string) (interface{}, bool) {
	return nil, false
}

func (c *Config) Put(key string, value interface{}) bool {
	return c.PutKey(c.NewKey(key), value)
}

func (c *Config) PutKey(key Key, value interface{}) bool {
	return false
}

func (c *Config) NewKey(source string) Key {
	return NewKey(source, c.Separator)
}

type Key []string

func NewKey(source, sep string) Key {
	return Key(strings.Split(source, sep))
}

type set map[string]interface{}

func newSet() set {
	return set(map[string]interface{}{})
}
