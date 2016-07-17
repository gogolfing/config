//Package json provides a Loader type that can be used in conjunction with
//the parent config package to create a config.Loader to load values from JSON
//encoded objects.
package json

import (
	"bytes"
	jsonlib "encoding/json"
	"io"

	"github.com/gogolfing/config"
)

//Loader is a collection of settings that can be used with config.NewReaderFuncLoader()
//in order to create a config.Loader that parses JSON objects.
//Loader itself is not a config.Loader.
//The empty valued Loader has sane defaults where all key, values found within the JSON
//are included in the resulting Values, JSON nulls are inserted at nil, and JSON numbers
//are converted to float64 and int64 equivalents.
//See the individual fields for overriding this behaviour.
//See the package examples for use with the config package.
type Loader struct {

	//KeyPrefix is a Key that all Keys found in the JSON must start with in order
	//to be included in the resulting config.Values.
	//Notice that an empty KeyPrefix means all Keys are matched.
	KeyPrefix config.Key

	//KeySuffix is a Key that all Keys found in the JSON must end with in order
	//to be included in the resulting config.Values.
	//Notice that an empty KeySuffix means all Keys are matched.
	KeySuffix config.Key

	//DiscardNull tells Loader whether or not to include JSON nulls as nil
	//in the resulting Values.
	//The zero value means all nulls are indeed included.
	DiscardNull bool

	//NumberAsString tells Loader whether to insert JSON numbers as their
	//string representations or to parse them into float64 or int64 representations
	//via encoding/json.Number.
	//The zero value means all numbers will be parsed into float64 and int64 types.
	//
	//If NumberAsString is false, when parsing a JSON number the following logic
	//is used.
	//If encoding/json.Number.Int64() returns an error, then Float64()'s result
	//is inserted.
	//If their is no error, then the int64 value is inserted.
	NumberAsString bool

	//KeyPartTransform is an optional function that is called (if not nil)
	//on each individual key part found in the JSON. The returned response from
	//this function is then used to create the resulting Key.
	//
	//See the examples for use of this function.
	KeyPartTransform func(string) string
}

//LoadString uses l's settings and returns the parsed Values and possible error
//from decoding in.
//It is sugar for l.LoadBytes([]byte(in)).
func (l *Loader) LoadString(in string) (*config.Values, error) {
	return l.LoadBytes([]byte(in))
}

//LoadBytes uses l's settings and returns the parsed Values and possible error
//from decoding in.
//It is sugar for l.LoadReader(bytes.NewReader(in)).
func (l *Loader) LoadBytes(in []byte) (*config.Values, error) {
	return l.LoadReader(bytes.NewReader(in))
}

//LoadReader uses l's settings and a encoding/json.Decoder to parse Values from in.
//If the call to encoding/json.Decoder.Decode() returns an error, then that error
//is returned with nil *Values.
//in must represent a JSON encoded object. Any other type will error.
//
//Notice that LoadReader is a config.ReaderFuncLoader and it is used in this manner
//in the examples.
func (l *Loader) LoadReader(in io.Reader) (*config.Values, error) {
	object, err := parseJson(in)
	if err != nil {
		return nil, err
	}
	values := config.NewValues()
	l.loadMapIntoValues(config.Key(nil), values, object)
	return values, nil
}

func (l *Loader) loadMapIntoValues(key config.Key, values *config.Values, object map[string]interface{}) {
	for keyPart, v := range object {
		if l.KeyPartTransform != nil {
			keyPart = l.KeyPartTransform(keyPart)
		}
		nextKey := key.AppendStrings(keyPart)
		if m, ok := v.(map[string]interface{}); ok {
			l.loadMapIntoValues(nextKey, values, m)
		} else {
			l.loadSingleIntoValues(nextKey, values, v)
		}
	}
}

func (l *Loader) loadSingleIntoValues(key config.Key, values *config.Values, value interface{}) {
	if !key.StartsWith(l.KeyPrefix) || !key.EndsWith(l.KeySuffix) {
		return
	}
	if value == nil {
		if l.DiscardNull {
			return
		}
		values.Put(key, nil)
	}
	switch v := value.(type) {
	case jsonlib.Number:
		l.loadNumberIntoValues(key, values, v)
	default:
		values.Put(key, v)
	}
}

func (l *Loader) loadNumberIntoValues(key config.Key, values *config.Values, num jsonlib.Number) {
	if l.NumberAsString {
		values.Put(key, num.String())
		return
	}
	i64, err := num.Int64()
	if err != nil {
		f64, _ := num.Float64()
		values.Put(key, f64)
	} else {
		values.Put(key, i64)
	}
}

func parseJson(in io.Reader) (map[string]interface{}, error) {
	dec := jsonlib.NewDecoder(in)
	dec.UseNumber()
	result := map[string]interface{}{}
	err := dec.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
