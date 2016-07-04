package config

type Loader interface {
	Load() (*Values, error)
}
