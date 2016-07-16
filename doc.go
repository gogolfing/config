//Package config provides types that allow loading, storing, retrieving, and
//removing arbitrary values that are referenced by keys.
//
//The Config Type
//
//You can create instances of Config with New().
//The zero valued &Config{} is not in a valid state and will likely cause panics
//if used.
//The New() func automatically sets Config's KeyParser to PeriodSeparatorKeyParser.
//This means that all string key parameters to Config methods will be converted
//to Key types in the manner of "a.b.c.d" -> Key([]string{"a", "b", "c", "d"}).
//See below for examples regarding storage and retrieval with keys.
//
//Inserting values into a Config is simple and type agnostic.
//The Put*() methods all return whether or not the internal set of values changed.
//
//	c := New()
//	c.Put("foo", "bar") //true
//	c.Put("foo", "bar") //false
//	c.Put("foo", 1024)  //true
//
//You can additionally use the Loader type to insert multiple values at once
//and from from varying sources instead of calling individual Put*() methods.
//
//	type sliceLoader []interface{}
//
//	func (sl sliceLoader) Load() (*Values, error) {
//		values := NewValues()
//		for i, v := range sl {
//			values.Put(NewKey("slice", fmt.Sprint(i)), v)
//		}
//		return values, nil
//	}
//
//	loader := sliceLoader([]interface{}{"hello", "world", 234, true})
//	c := New()
//	c.PutLoaders(loader) //true, nil
//
//	c.Get("slice.0") //hello
//	c.Get("slice.1") //world
//	c.Get("slice.2") //234
//	c.Get("slice.3") //true
//
//See the loader subdirectory for information about pre-written loaders for
//common use cases.
package config
