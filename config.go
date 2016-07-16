package config

import "sync"

//Type Config provides an methods to store, retrieve, and remove arbitrary values
//that are referenced by keys.
//
//The keys to a Config type are of type string. These string keys are parsed into
//Key types via a KeyParser.
//
//Key types that are parsed by a Config type are then used to reference into Values.
//Values is the storage type providing all functionality to Config.
//
//Config is safe for use by multiple goroutines.
//Though the KeyParser is not protected from concurrent use, implementations can be
//implemented in a safe manner.
//
//The zero value for *Config is not in a valid state and will likely cause panics if
//used.
type Config struct {
	//KeyParser that turns strings into Keys that are then used with
	//this Config's underlying Values.
	//This enabled easier access to Config value with a simple string as opposed
	//to a Key type.
	KeyParser KeyParser

	values *Values

	lock    *sync.Mutex
	loaders []Loader
}

//New creates a new *Config with an empty Values and Loaders.
//KeyParser is set to PeriodSeparatorKeyParser.
func New() *Config {
	return &Config{
		KeyParser: PeriodSeparatorKeyParser,

		values: NewValues(),

		lock:    &sync.Mutex{},
		loaders: []Loader{},
	}
}

//AddLoaders adds Loaders to an internal slice of loaders.
//There is no check for duplicates or nil Loaders.
//This method servers as a utility to store Loaders associated with a Config
//to be merged into c at a later time by c.LoadAll().
func (c *Config) AddLoaders(loaders ...Loader) *Config {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.loaders = append(c.loaders, loaders...)
	return c
}

//LoadAll is a helper for c.PutLoaders() called with all Loaders
//added previously with c.AddLoaders().
func (c *Config) LoadAll() (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.PutLoaders(c.loaders...)
}

//Values returns the internal Values used for storage.
//The Values type is sage for use by multiple goroutines independent of
//the Config type.
//The internal reference is simply returned, and thus modifications to the
//returned *Values will affect c.
func (c *Config) Values() *Values {
	return c.values
}

//EqualValues is sugar for c.Values().Equal(other.Values()).
func (c *Config) EqualValues(other *Config) bool {
	return c.values.Equal(other.values)
}

//Clone creates and returns a new *Config with KeyParser and added loaders
//shallow copied, and with *Values cloned via *Values.Clone().
func (c *Config) Clone() *Config {
	c.lock.Lock()
	defer c.lock.Unlock()

	return &Config{
		KeyParser: c.KeyParser,

		values: c.values.Clone(),

		lock:    &sync.Mutex{},
		loaders: c.loaders,
	}
}

//GetInt64 returns an int64 casted integer type stored at key.
//The zero value for int64 is returned if an integer type does not exist at key.
func (c *Config) GetInt64(key string) (i int64) {
	i64, _ := c.GetInt64Ok(key)
	return i64
}

//GetInt64Ok returns an int64 casted integer type stored at key.
//The zero value for int64 is returned if an integer type does not exist at key.
//The return value ok indicates whether or not an integer type actually exists at key.
func (c *Config) GetInt64Ok(key string) (i int64, ok bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return 0, false
	}
	i, ok = int64(0), false
	switch iType := v.(type) {
	case uint8:
		i, ok = int64(iType), true
	case int8:
		i, ok = int64(iType), true
	case uint16:
		i, ok = int64(iType), true
	case int16:
		i, ok = int64(iType), true
	case uint32:
		i, ok = int64(iType), true
	case int32:
		i, ok = int64(iType), true
	case uint:
		i, ok = int64(iType), true
	case int:
		i, ok = int64(iType), true
	case uint64:
		i, ok = int64(iType), true
	case int64:
		i, ok = iType, true
	}
	return
}

//GetBool returns a bool stored at key.
//The zero value for bool is returned if a bool does not exist at key.
func (c *Config) GetBool(key string) (b bool) {
	b, _ = c.GetBoolOk(key)
	return
}

//GetBoolOk returns a bool stored at key.
//The zero value for bool is returned if a bool does not exist at key.
//The return value ok indicates whether or not a bool actually exists at key.
func (c *Config) GetBoolOk(key string) (b bool, ok bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return false, false
	}
	b, ok = v.(bool)
	return
}

//GetString returns a string stored at key.
//The zero value for string is returned if a string does not exist at key.
func (c *Config) GetString(key string) (s string) {
	s, _ = c.GetStringOk(key)
	return
}

//GetStringOk returns a string stored at key.
//The zero value for string is returned if a string does not exist at key.
//The return value ok indicates whether or not a string actually exists at key.
func (c *Config) GetStringOk(key string) (s string, ok bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return "", false
	}
	s, ok = v.(string)
	return
}

//GetFloat64 returns a float64 casted floating point type stored at key.
//The zero value for float64 is returned if a floating point type does not exist at key.
func (c *Config) GetFloat64(key string) (f float64) {
	f, _ = c.GetFloat64Ok(key)
	return
}

//GetFloat64 returns a float64 casted floating point type stored at key.
//The zero value for float64 is returned if a floating point type does not exist at key.
//The return value ok indicates whether or not a floating point type actually exists at key.
func (c *Config) GetFloat64Ok(key string) (f float64, ok bool) {
	v, ok := c.GetOk(key)
	if !ok {
		return 0, false
	}
	f, ok = float64(0), false
	switch fType := v.(type) {
	case float32:
		f, ok = float64(fType), true
	case float64:
		f, ok = fType, true
	}
	return
}

//GetValues returns a *Values stored at key.
//This means that there exists some value stored at a longer Key.
//The returned *Values is cloned and thus changes to v do not affect c and vice versa.
//nil is returned if a *Values does not exist at key.
func (c *Config) GetValues(key string) (v *Values) {
	v, _ = c.GetValuesOk(key)
	return
}

//GetValuesOk returns a *Values stored at key.
//This means that there exists some value stored at a longer Key.
//The returned *Values is cloned and thus changes to v do not affect c and vice versa.
//nil is returned if more values do not exist at key.
//The return value ok indicates whether or not there are more values stored at key.
func (c *Config) GetValuesOk(key string) (v *Values, ok bool) {
	vInterface, ok := c.GetOk(key)
	if !ok {
		return nil, false
	}
	v, ok = vInterface.(*Values)
	return
}

//Get is sugar for c.GetKey(c.NewKey(key)).
//It returns a raw interface{} value stored at key or nil if a value does not
//exist at key.
func (c *Config) Get(key string) (v interface{}) {
	return c.GetKey(c.NewKey(key))
}

//GetOk is sugar for c.GetKeyOk(c.NewKey(key)).
//It returns a raw interface{} value stored at key or nil if a value does not
//exist at key.
//The return value ok indicates whether or not any value is actually stored at key.
func (c *Config) GetOk(key string) (v interface{}, ok bool) {
	return c.GetKeyOk(c.NewKey(key))
}

//GetKey returns Get(key) called on c's internal *Values instance.
//It returns a raw interface{} value stored at key or nil if a value does not
//exist at key.
func (c *Config) GetKey(key Key) (v interface{}) {
	return c.values.Get(key)
}

//GetKeyOk returns GetKey(key) called on c's internal *Values instance.
//It returns a raw interface{} value stored at key or nil if a value does not
//exist at key.
//The return value ok indicates whether or not any value is actually stored at key.
func (c *Config) GetKeyOk(key Key) (v interface{}, ok bool) {
	return c.values.GetOk(key)
}

func (c *Config) Merge(other *Config) bool {
	return c.values.Merge(nil, other.values)
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

func (c *Config) Remove(key string) (interface{}, bool) {
	return c.values.Remove(c.NewKey(key))
}

func (c *Config) NewKey(k string) Key {
	return c.KeyParser.Parse(k)
}
