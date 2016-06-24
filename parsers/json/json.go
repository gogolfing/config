package json

import (
	jsonlib "encoding/json"
	"fmt"
	"io"

	"github.com/gogolfing/config"
)

type Parser struct {
	KeepNil bool
}

func NewParser(keepNil bool) config.Parser {
	return Parser{
		KeepNil: keepNil,
	}
}

func (p Parser) Parse(r io.Reader) ([]config.KeyValue, error) {
	dec := jsonlib.NewDecoder(r)
	dec.UseNumber()
	m := map[string]interface{}{}
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	fmt.Println(m)
	return p.convertValueToKeyValues(nil, m), nil
}

func Convert(value interface{}) interface{} {
	if value == nil {
		if p.KeepNil {
			return value
		}
		return 
	}
	value = config.Convert(value)
	if 
}

func (p Parser) convertValueToKeyValues(key config.Key, value interface{}) []config.KeyValue {
	switch v := value.(type) {
	case nil:
		if p.KeepNil {
			return []config.KeyValue{config.NewKeyValue(key, nil)}
		}
		return []config.KeyValue{}

	case map[string]interface{}:
		return p.convertMapToKeyValues(key, v)
	}
	//all other JSON unmarshalled values are fine the way they are.
	return []config.KeyValue{config.NewKeyValue(key, value)}
}

func (p Parser) convertMapToKeyValues(key config.Key, m map[string]interface{}) []config.KeyValue {
	return nil
}
