package composer

type InterfaceComposeUnit interface {
	GetName() string

	GetType(item string) string
	GetLabel(item string) string
	GetDescription(item string) string

	Process(in map[string]interface{}) (map[string]interface{}, error)
}

type InterfaceComposer interface {
	RegisterUnit(unit InterfaceComposeUnit) error
	UnRegisterUnit(unit InterfaceComposeUnit) error

	Process(unit InterfaceComposeUnit, in map[string]interface{}) (interface{}, error)
}