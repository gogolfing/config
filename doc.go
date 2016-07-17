//Package config provides types that allow loading, storing, retrieving, and
//removing arbitrary values that are referenced by keys.
//
//The Config Type
//
//We can create instances of Config with New().
//The zero valued &Config{} is not in a valid state and will likely cause panics
//if used.
//The New() func automatically sets Config's KeyParser to PeriodSeparatorKeyParser.
//This means that all string key parameters to Config methods will be converted
//to Key types in the manner of "a.b.c.d" -> Key([]string{"a", "b", "c", "d"}).
//A Config's KeyParser field can be changed before use to override this functionality.
//
//Inserting values into a Config is simple and type agnostic.
//The Put*() methods all return whether or not the internal set of values changed.
//
//	c := New()
//	c.Put("foo", "bar") //true
//	c.Put("foo", "bar") //false
//	c.Put("foo", 1024)  //true
//
//We can additionally use the Loader type to insert multiple values at once
//and from varying sources instead of calling individual Put*() methods.
//
//	type sliceLoader []interface{}
//
//	func (sl sliceLoader) Load() (*Values, error) {
//		values := NewValues()
//		for i, v := range sl {
//			key := NewKey("slice", fmt.Sprint(i))
//			values.Put(key, v)
//		}
//		return values, nil
//	}
//
//	loader := sliceLoader([]interface{}{"hello", "world", 234, true})
//	c := New()
//	c.MergeLoaders(loader) //true, nil
//
//	c.Get("slice.0") //hello
//	c.Get("slice.1") //world
//	c.Get("slice.2") //234
//	c.Get("slice.3") //true
//
//See the Loader documentation and the loader subdirectory for information
//about pre-written loaders for common use cases.
//
//The Values and Key Types
//
//Type Values provides the actaul storage and retrieval of interface{} values.
//The values stored in a Values type are referenced by the Key type.
//Values is implemented as a tree with possible multiple children at each node.
//Each individual string in a Key is the "pointer" to the subtree of possibly
//more values.
//
//For example:
//
//	//the following code...
//
//	v := NewValues()
//	v.Put(NewKey("a", "b"), 1)
//	v.Put(NewKey("a", "c"), 2)
//	v.Put(NewKey("d"), 3)
//
//	//results in this structure.
//	//       __root__
//	//       |      |
//	//       a      d
//	//       |      |
//	//    b--+--c   3
//	//    |     |
//	//    1     2
//
//Continuing from this example, if we were to call v.Put(NewKey("a"), true),
//that would completely remove the [a b] -> 1 and [a c] -> 2 associations from v.
//
//If we were to call v.Put(NewKey("a", "e"), "foobar"), that would result in the
//[a b] and [a c] associations remaining and the new [a e] -> "foobar" within v.
//
package config
