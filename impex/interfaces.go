package impex

// InterfaceImpexImportCmd is an interface used to work with registered import commands
type InterfaceImpexImportCmd interface {
	Init(args []string, exchange map[string]interface{}) error
	Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error)
	Test(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error)
}

// InterfaceImpexModel is an interface model should implement to make it work with impex service
type InterfaceImpexModel interface {
	Import(item map[string]interface{}, testMode bool) (map[string]interface{}, error)
	Export(func(map[string]interface{}) bool) error
}
