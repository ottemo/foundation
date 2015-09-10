package composer

type InterfaceComposeUnit interface {
	GetName() string

	ListItems() []string

	GetType(item string) string
	GetLabel(item string) string
	GetDescription(item string) string

	Process(args map[string]interface{}, composer InterfaceComposer) (interface{}, error)
}

type InterfaceComposer interface {
	RegisterUnit(unit InterfaceComposeUnit) error
	UnRegisterUnit(unit InterfaceComposeUnit) error
	ListUnits() []InterfaceComposeUnit

	GetUnit(name string) InterfaceComposeUnit
	SearchUnits(namePattern string, typeFilter map[string]interface{}) []InterfaceComposeUnit

	Process(in interface{}, rules map[string]interface{}) (bool, error)
}

func InKey(name string) {
	return ConstInPrefix + name
}

func OutKey(name string) {
	return ConstOutPrefix + name
}