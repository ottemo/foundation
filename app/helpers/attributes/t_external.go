package attributes

// GenericProduct type implements:
// 	- InterfaceExternalAttributes
// 	- InterfaceObject
// 	- InterfaceStorable

// Init initializes helper instance before usage
func (it *ModelExternalAttributes) Init(model string, collection string) (*ModelExternalAttributes, error) {
	it.model = model
	it.values = make(map[string]interface{})

	return it, nil
}