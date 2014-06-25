package conf

type DefaultConfigItem struct {
	Name      string
	Validator func(interface{}) (interface{}, bool)
	Default   interface{}
	Value     interface{}
}

type DefaultConfig struct {
	configValues map[string]*DefaultConfigItem
}
