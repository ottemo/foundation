package composer

type InterfaceComposeUnit interface {
	GetName() string

	ListItems() []string

	GetType(item string) string
	ValidateType(item string, inType string) bool

	GetLabel(item string) string
	GetDescription(item string) string

	Process(args map[string]interface{}, composer InterfaceComposer) (map[string]interface{}, error)
}

type InterfaceComposer interface {
	RegisterUnit(unit InterfaceComposeUnit) error
	UnRegisterUnit(unit InterfaceComposeUnit) error
	ListUnits() []InterfaceComposeUnit

	GetUnit(name string) InterfaceComposeUnit
	SearchUnits(namePattern string, typeFilter map[string]interface{}) []InterfaceComposeUnit

	Validate(item interface{}, value interface{}) (bool, error)
}
