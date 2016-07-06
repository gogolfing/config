package json

import (
	"bytes"
	jsonlib "encoding/json"
	"io"

	"github.com/gogolfing/config"
)

type Loader struct {
	KeyPrefix        config.Key
	KeySuffix        config.Key
	DiscardNil       bool
	NumberAsString   bool
	KeyPartTransform func(string) string
}

func (l *Loader) LoadString(in string) (*config.Values, error) {
	return l.LoadBytes([]byte(in))
}

func (l *Loader) LoadBytes(in []byte) (*config.Values, error) {
	return l.LoadReader(bytes.NewReader(in))
}

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
		l.loadSingleIntoValues(key.AppendStrings(keyPart), values, v)
	}
}

func (l *Loader) loadSingleIntoValues(key config.Key, values *config.Values, value interface{}) {
	if !key.StartsWith(l.KeyPrefix) || !key.EndsWith(l.KeySuffix) {
		return
	}
	if value == nil {
		if l.DiscardNil {
			return
		}
		values.Put(key, nil)
	}
	switch v := value.(type) {
	case map[string]interface{}:
		l.loadMapIntoValues(key, values, v)
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
