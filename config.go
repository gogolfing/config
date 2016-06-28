package config

type Config struct {
	KeyParser KeyParser

	values *Values
}

func New() *Config {
	return &Config{
		KeyParser: PeriodSeparatorKeyParser,

		values: NewValues(),
	}
}

func (c *Config) Put(key string, value interface{}) bool {
	return c.PutKey(c.NewKey(key), value)
}

func (c *Config) PutKey(key Key, value interface{}) bool {
	return c.values.Put(key, value)
}

// func (c *Config) PutLoaders(loaders ...Loader) (bool, error)

func (c *Config) NewKey(k string) Key {
	return c.KeyParser.Parse(k)
}
