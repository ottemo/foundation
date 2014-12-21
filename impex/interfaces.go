package impex

// InterfaceImpexImportCmd is an interface used to work with registered import commands
type InterfaceImpexImportCmd interface {
	Init(args []string, exchange map[string]interface{}) error
	Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error)
}
