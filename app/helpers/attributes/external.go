package attributes


// Init initializes helper instance before usage
func (it *ExternalAttributes) Init(model string, collection string) (*ExternalAttributes, error) {
	it.model = model
	it.values = make(map[string]interface{})

	return it, nil
}