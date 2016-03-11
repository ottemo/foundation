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

// KeyValueSorter is a container for key values
type KeyValueSorter struct {
	Key   string
	Value int
}

// KeyValueList is a sort.Interface implementor for map[string]int
type KeyValueList []KeyValueSorter

// Len returns length of given slice
func (it *funcSorter) Len() int {
	return len(it.data)
}

// Len returns length of given slice
func (it *mapSorter) Len() int {
	return len(it.data)
}

// Len returns the length of Key Value List
func (it KeyValueList) Len() int {
	return len(it)
}

// Swap switches slice values between themselves
func (it *funcSorter) Swap(i, j int) {
	it.data[i], it.data[j] = it.data[j], it.data[i]
}

// Swap switches slice values between themselves
func (it *mapSorter) Swap(i, j int) {
	it.data[i], it.data[j] = it.data[j], it.data[i]
}

// Swap switches the values in the Key Value List
func (it KeyValueList) Swap(i, j int) {
	it[i], it[j] = it[j], it[i]
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

// Less compares the values
func (it KeyValueList) Less(i, j int) bool {
	return it[i].Value < it[j].Value
}

// SortByFunc sorts slice with a given comparator function
// 	- to sort in ascending order pass reverse as false
//      - to sort in descending order pass reverse as true
func SortByFunc(data interface{}, reverse bool, less func(a, b interface{}) bool) []interface{} {
	sortable := &funcSorter{data: InterfaceToArray(data), less: less}
	if reverse {
		sort.Sort(sort.Reverse(sortable))
	} else {
		sort.Sort(sortable)
	}
	return sortable.data
}

// SortMapByKeys sorts given map by specified keys values
// 	- to sort in ascending order pass reverse as false
//      - to sort in descending order pass reverse as true
func SortMapByKeys(data []map[string]interface{}, reverse bool, keys ...string) []map[string]interface{} {
	sortable := &mapSorter{data: data, keys: keys}
	if reverse {
		sort.Sort(sort.Reverse(sortable))
	} else {
		sort.Sort(sortable)
	}
	return sortable.data
}

// SortByInt will sort a map[string]int in ascending order
// 	- to sort in ascending order pass reverse as false
//      - to sort in descending order pass reverse as true
func SortByInt(count map[string]int, reverse bool) KeyValueList {
	list := make(KeyValueList, len(count))
	i := 0
	for k, v := range count {
		list[i] = KeyValueSorter{k, v}
		i++
	}
	if reverse {
		sort.Sort(sort.Reverse(list))
	} else {
		sort.Sort(list)
	}
	return list
}
