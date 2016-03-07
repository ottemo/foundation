package utils

import (
	"sort"
)

// funcSorter is a sort.Interface implementor for []interface{} type
type funcSorter struct {
	data []interface{}
	less func(interface{}, interface{}) bool
}

// mapSorter is a sort.Interface implementor for []map[string]interface type
type mapSorter struct {
	keys []string
	data []map[string]interface{}
}

// Len returns length of given slice
func (it *funcSorter) Len() int {
	return len(it.data)
}

// Len returns length of given slice
func (it *mapSorter) Len() int {
	return len(it.data)
}

// Swap switches slice values between themselves
func (it *funcSorter) Swap(i, j int) {
	it.data[i], it.data[j] = it.data[j], it.data[i]
}

// Swap switches slice values between themselves
func (it *mapSorter) Swap(i, j int) {
	it.data[i], it.data[j] = it.data[j], it.data[i]
}

// Less compares slice values with a given function
func (it *funcSorter) Less(i, j int) bool {
	return it.less(it.data[i], it.data[j])
}

// Less compares slice values between themselves
func (it *mapSorter) Less(i, j int) bool {
	for _, key := range it.keys {
		a := it.data[i][key]
		b := it.data[j][key]

		// nil values equals, preventing time loose on conversion
		if a == nil && b == nil {
			continue
		}

		// looking for not nil values
		x := a
		if a == nil {
			x = b
		}

		// comparable types are either string or number
		switch x.(type) {
		case string:
			a := InterfaceToString(a)
			b := InterfaceToString(b)

			if a != b {
				return a < b
			}
		default:
			a := InterfaceToFloat64(a)
			b := InterfaceToFloat64(b)

			if a != b {
				return a < b
			}
		}
	}

	return false
}

// SortByFunc sorts slice with a given comparator function
func SortByFunc(data interface{}, less func(a, b interface{}) bool) []interface{} {
	sortable := &funcSorter{data: InterfaceToArray(data), less: less}
	sort.Sort(sortable)
	return sortable.data
}

// SortMapByKeys sorts given map by specified keys values
func SortMapByKeys(data []map[string]interface{}, keys ...string) []map[string]interface{} {
	sortable := &mapSorter{data: data, keys: keys}
	sort.Sort(sortable)
	return sortable.data
}
