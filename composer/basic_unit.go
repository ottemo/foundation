package composer

import (
	"sort"
	"strings"
)

func (it *BasicUnit) GetName() string {
	return it.Name
}

func (it *BasicUnit) ValidateType(item string, inType string) bool {
	if it.Validator != nil {
		return it.Validator(item, inType)
	}

	inType = strings.TrimSpace(inType)
	allowed := it.GetType(item)
	for _, typeName := range strings.Split(allowed, "|") {
		if strings.TrimSpace(typeName) == inType {
			return true
		}
	}

	return false
}

func (it *BasicUnit) ListItems() []string {
	var result []string

	for itemName := range it.Type {
		result = append(result, itemName)
	}

	sort.Strings(result)
	return result
}

func (it *BasicUnit) GetType(item string) string {
	if value, present := it.Type[item]; present {
		return value
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

func (it *BasicUnit) IsRequired(item string) bool {
	if value, present := it.Required[item]; present {
		return value
	}
	return true
}

func (it *BasicUnit) Process(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
	if it.Action != nil {
		return it.Action(in, args, composer)
	}
	return nil, nil
}
