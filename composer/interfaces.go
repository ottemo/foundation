package composer

type InterfaceComposeUnit interface {
	GetName() string

	ListItems() []string

	GetType(item string) string
	ValidateType(item string, inType string) bool

	IsRequired(item string) bool

	GetLabel(item string) string
	GetDescription(item string) string

	Process(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error)
}

type InterfaceComposer interface {
	GetName() string

	RegisterUnit(unit InterfaceComposeUnit) error
	UnRegisterUnit(unit InterfaceComposeUnit) error
	ListUnits() []InterfaceComposeUnit

	GetUnit(name string) InterfaceComposeUnit
	SearchUnits(namePattern string, typeFilter map[string]string) []InterfaceComposeUnit

	Check(in interface{}, rule interface{}) (bool, error)
}
