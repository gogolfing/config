package config

type Loader interface {
	Load() (*Values, error)
}

/*
type EnvLoader struct {
	Prefix    string
	Separator string
	Transform func(string) string
}

func NewUnderscoreLowerEnvLoader(prefix string) Loader {
	return NewEnvLoader(prefix, "_", strings.ToLower)
}

func NewEnvLoader(prefix, sep string, transform func(string) string) Loader {
	return EnvLoader{
		Prefix:    prefix,
		Separator: sep,
		Transform: transform,
	}
}

func (e EnvLoader) Load() ([]KeyValue, error) {
	result := []KeyValue{}
	env := os.Environ()
	for _, envKeyValue := range env {
		keyValue := e.getKeyValueFromEnv(envKeyValue)
		if keyValue != nil {
			result = append(result, *keyValue)
		}
	}
	return result, nil
}

func (e EnvLoader) getKeyValueFromEnv(value string) *KeyValue {
	index := strings.Index(value, "=")
	if index < 0 {
		index = 0
	}
	key, value := value[:index], value[index+1:]
	transform := e.getTransform()
	key = transform(key)
	prefix := transform(e.Prefix)
	if !strings.HasPrefix(key, prefix) {
		return nil
	}
	key = strings.TrimPrefix(key, prefix)
	return &KeyValue{
		Key:   NewKey(key, e.Separator),
		Value: value,
	}
}

func (e EnvLoader) getTransform() func(string) string {
	if e.Transform == nil {
		return NoTransform
	}
	return e.Transform
}

func NoTransform(s string) string {
	return s
}

type FileParserLoader struct {
	Path string
	Parser
}

func NewFileParserLoader(path string, p Parser) Loader {
	return FileParserLoader{
		Path:   path,
		Parser: p,
	}
}

func (f FileParserLoader) Load() ([]KeyValue, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	r := ReaderParserLoader{
		Reader: file,
		Parser: f.Parser,
	}
	return r.Load()
}

type ReaderParserLoader struct {
	io.Reader
	Parser
}

func NewReaderParserLoader(r io.Reader, p Parser) Loader {
	return ReaderParserLoader{
		Reader: r,
		Parser: p,
	}
}

func (r ReaderParserLoader) Load() ([]KeyValue, error) {
	return r.Parser.Parse(r.Reader)
}
*/
