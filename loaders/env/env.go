package env

import (
	"os"
	"strings"

	"github.com/gogolfing/config"
)

//Equal is the string around which results from os.Environ() are split into keys
//and their respective values.
const Equal = "="

//UnderscoreSeparatorKeyParser is a config.SeparatorKeyParser that parses keys
//around the "_" character.
const UnderscoreSeparatorKeyParser = config.SeparatorKeyParser("_")

type prefixParserLoader struct {
	prefix string
	parser config.KeyParser
}

//NewPrefixLowerUnderscoreLoader creates a config.Loader that
//reads in all entries from os.Environ() and inserts into the resulting Values all
//key, value associations whose keys start with prefix.
//The key inserted is parsed with UnderscoreSeparatorKeyParser after prefix is removed.
func NewPrefixLowerUnderscoreLoader(prefix string) config.Loader {
	parser := config.KeyParserFunc(func(k string) config.Key {
		return UnderscoreSeparatorKeyParser.Parse(strings.ToLower(k))
	})
	return NewPrefixParserLoader(prefix, parser)
}

//NewPrefixParserLoader creates a config.Loader that
//reads in all entries from os.Environ() and inserts into the resulting Values all
//key, value associations whose keys start with prefix.
//The key inserted is parsed with parser after prefix is removed.
func NewPrefixParserLoader(prefix string, parser config.KeyParser) config.Loader {
	return &prefixParserLoader{
		prefix: prefix,
		parser: parser,
	}
}

func (p *prefixParserLoader) Load() (*config.Values, error) {
	values := config.NewValues()
	environment := os.Environ()
	for _, envVar := range environment {
		key, value := p.loadPossibleEnvironmentVariable(envVar)
		if !key.IsEmpty() {
			values.Put(key, value)
		}
	}
	return values, nil
}

func (p *prefixParserLoader) loadPossibleEnvironmentVariable(envVar string) (config.Key, string) {
	equalIndex := strings.Index(envVar, Equal)
	key, value := envVar[:equalIndex], envVar[equalIndex+1:]
	if !strings.HasPrefix(key, p.prefix) {
		return config.Key(nil), ""
	}
	key = strings.TrimPrefix(key, p.prefix)
	return p.parser.Parse(key), value
}
