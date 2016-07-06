package env

import (
	"os"
	"strings"

	"github.com/gogolfing/config"
)

const Equal = "="

const UnderscoreSeparatorKeyParser = config.SeparatorKeyParser("_")

type prefixParserLoader struct {
	prefix string
	parser config.KeyParser
}

func NewPrefixLowerUnderscoreLoader(prefix string) config.Loader {
	parser := config.KeyParserFunc(func(k string) config.Key {
		return UnderscoreSeparatorKeyParser.Parse(strings.ToLower(k))
	})
	return NewPrefixParserLoader(prefix, parser)
}

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
	if equalIndex < 0 {
		return config.Key(nil), ""
	}
	key, value := envVar[:equalIndex], envVar[equalIndex+1:]
	if !strings.HasPrefix(key, p.prefix) {
		return config.Key(nil), ""
	}
	key = strings.TrimPrefix(key, p.prefix)
	return p.parser.Parse(key), value
}
