package config

// IniConfig is an initialization interface for reading INI file values
type IniConfig interface {
	GetValue(Name string) string
	List() []string
}

// Config is an interface for working with configuration entities and values
type Config interface {
	Register(Name string, Validator func(interface{}) (interface{}, bool), Default interface{}) error
	Destroy(Name string) error

	GetValue(Name string) interface{}
	SetValue(Name string, Value interface{}) error

	List() []string

	Load() error
	Save() error
}
