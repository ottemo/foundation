package composer

func (it *BasicUnit) GetName() string {
	return it.Name
}

func (it *BasicUnit) GetType(item string) string {
	switch item {
		case "":
			return
	}
	return ""
}

func (it *BasicUnit) GetLabel(item string) string {
	if value, present := it.Label[item]; present {
		return value
	}
	return ""
}

func (it *BasicUnit) GetDescription(item string) string {
	if value, present := it.Description[item]; present {
		return value
	}
	return ""
}

func (it *BasicUnit) Process(in map[string]interface{}) (map[string]interface{}, error) {
	if it.Action != nil {
		return it.Action(in)
	}
	return nil, nil
}