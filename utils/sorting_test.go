package utils

import "testing"

// TestSortMapByKeys validates sort map by keys implementation
func TestSortMapByKeys(t *testing.T) {
	data := []map[string]interface{}{
		{"a": 3, "b": "B"},
		{"a": 3, "b": "A"},
		{"a": 1, "b": "C"},
		{"a": 2, "b": "B"},
		{"a": 2, "c": 33},
	}

	result := SortMapByKeys(data, "a", "b")

	// expecting:
	// 		{"a": 1, "b": "C"},
	//		{"a": 2, "c": 33},
	//		{"a": 2, "b": "B"}
	//		{"a": 3, "b": "A"},
	//		{"a": 3, "b": "B"},
	if result[0]["a"] != 1 || result[1]["c"] != 33 || result[4]["b"] != "B" {
		t.Error("Unexpected sort maps result: ", result)
	}
}

// TestSortByFunc validates sort by function implementation
func TestSortByFunc(t *testing.T) {
	data := []interface{}{"1", 33, "8", "13", 5, true}

	result := SortByFunc(data, func(a, b interface{}) bool {
		return InterfaceToInt(a) < InterfaceToInt(b)
	})

	// expecting: [true, "1", 5, "8", "13", 33]
	if result[1] != "1" || result[5] != 33 {
		t.Error("Unexpected sort by func result: ", result)
	}
}
