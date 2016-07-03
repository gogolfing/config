package config

import "sync"

type Config struct {
	KeyParser KeyParser

	values *Values

	lock    *sync.Mutex
	loaders []Loader
}

func New() *Config {
	return &Config{
		KeyParser: PeriodSeparatorKeyParser,

		values: NewValues(),

		lock:    &sync.Mutex{},
		loaders: []Loader{},
	}
}

func (c *Config) AddLoaders(loaders ...Loader) *Config {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.loaders = append(c.loaders, loaders...)
	return c
}

func (c *Config) LoadAll() (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.PutLoaders(c.loaders...)
}

func (c *Config) GetInt(key string) int {
	i, _ := c.GetIntOk(key)
	return i
}

func (c *Config) GetIntOk(key string) (int, bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return 0, false
	}
	result, ok := 0, false
	switch i := v.(type) {
	case uint8:
		result, ok = int(i), true
	case int8:
		result, ok = int(i), true
	case uint16:
		result, ok = int(i), true
	case int16:
		result, ok = int(i), true
	case uint32:
		result, ok = int(i), true
	case int32:
		result, ok = int(i), true
	case int:
		result, ok = i, true
	case uint64:
		result, ok = int(i), true
	case int64:
		result, ok = int(i), true
	}
	return result, ok
}

func (c *Config) GetInt64(key string) int64 {
	i64, _ := c.GetInt64Ok(key)
	return i64
}

func (c *Config) GetInt64Ok(key string) (int64, bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return 0, false
	}
	result, ok := int64(0), false
	switch i := v.(type) {
	case uint8:
		result, ok = int64(i), true
	case int8:
		result, ok = int64(i), true
	case uint16:
		result, ok = int64(i), true
	case int16:
		result, ok = int64(i), true
	case uint32:
		result, ok = int64(i), true
	case int32:
		result, ok = int64(i), true
	case int:
		result, ok = int64(i), true
	case uint64:
		result, ok = int64(i), true
	case int64:
		result, ok = i, true
	}
	return result, ok
}

func (c *Config) GetBool(key string) bool {
	b, _ := c.GetBoolOk(key)
	return b
}

func (c *Config) GetBoolOk(key string) (bool, bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

func (c *Config) GetString(key string) string {
	s, _ := c.GetStringOk(key)
	return s
}

func (c *Config) GetStringOk(key string) (string, bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func (c *Config) GetFloat64(key string) float64 {
	f, _ := c.GetFloat64Ok(key)
	return f
}

func (c *Config) GetFloat64Ok(key string) (float64, bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return 0, false
	}
	result, ok := float64(0), false
	switch f := v.(type) {
	case float32:
		result, ok = float64(f), true
	case float64:
		result, ok = f, true
	}
	return result, ok
}

func (c *Config) GetValues(key string) *Values {
	v, _ := c.GetValuesOk(key)
	return v
}

func (c *Config) GetValuesOk(key string) (*Values, bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return nil, false
	}
	values, ok := v.(*Values)
	return values, ok
}

func (c *Config) Get(key string) interface{} {
	return c.GetKey(c.KeyParser.Parse(key))
}

func (c *Config) GetOk(key string) (interface{}, bool) {
	return c.GetKeyOk(c.KeyParser.Parse(key))
}

func (c *Config) GetKey(key Key) interface{} {
	return c.values.Get(key)
}

func (c *Config) GetKeyOk(key Key) (interface{}, bool) {
	return c.values.GetOk(key)
}

func (c *Config) Put(key string, value interface{}) bool {
	return c.PutKey(c.NewKey(key), value)
}

func (c *Config) PutKey(key Key, value interface{}) bool {
	return c.values.Put(key, value)
}

func (c *Config) PutLoaders(loaders ...Loader) (bool, error) {
	temp := NewValues()
	changed := false
	for _, loader := range loaders {
		loaderValues, err := loader.Load()
		if err != nil {
			return false, err
		}
		changed = temp.Merge(nil, loaderValues) || changed
	}
	c.values.Merge(nil, temp)
	return changed, nil
}

func (c *Config) NewKey(k string) Key {
	return c.KeyParser.Parse(k)
}
