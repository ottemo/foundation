package composer

import (
	"sort"
)

func (it *BasicType) GetName() string {
	return it.Name
}

func (it *BasicType) ListItems() []string {
	var result []string

	for itemName := range it.Type {
		result = append(result, itemName)
	}

	sort.Strings(result)
	return result
}

func (it *BasicType) GetType(item string) string {
	if value, present := it.Type[item]; present {
		return value
	}
	return ""
}

func (it *BasicType) GetLabel(item string) string {
	if value, present := it.Label[item]; present {
		return value
	}
	return ""
}

func (it *BasicType) GetDescription(item string) string {
	if value, present := it.Description[item]; present {
		return value
	}
	return ""
}
