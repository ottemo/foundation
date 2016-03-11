package utils

import "testing"

// TestSortMapByKeysAscend validates sort map by keys implementation in ascending order
func TestSortMapByKeysAscend(t *testing.T) {
	data := []map[string]interface{}{
		{"a": 3, "b": "B"},
		{"a": 3, "b": "A"},
		{"a": 1, "b": "C"},
		{"a": 2, "b": "B"},
		{"a": 2, "c": 33},
	}

	result := SortMapByKeys(data, false, "a", "b")

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

// TestSortMapByKeysDescend validates sort map by keys implementation in decending order
func TestSortMapByKeysDescend(t *testing.T) {
	data := []map[string]interface{}{
		{"a": 3, "b": "B"},
		{"a": 3, "b": "A"},
		{"a": 1, "b": "C"},
		{"a": 2, "b": "B"},
		{"a": 2, "c": 33},
	}

	result := SortMapByKeys(data, true, "a", "b")

	// expecting:
	//		{"a": 3, "b": "B"},
	//		{"a": 3, "b": "A"},
	//		{"a": 2, "b": "B"}
	//		{"a": 2, "c": 33},
	// 		{"a": 1, "b": "C"},
	if result[0]["a"] != 3 || result[1]["b"] != "A" || result[4]["b"] != "C" {
		t.Error("Unexpected sort maps result: ", result)
	}
}

// TestSortByFuncAscend validates sort by function implementation
func TestSortByFuncAscend(t *testing.T) {
	data := []interface{}{"1", 33, "8", "13", 5, true}

	result := SortByFunc(data, false, func(a, b interface{}) bool {
		return InterfaceToInt(a) < InterfaceToInt(b)
	})

	// expecting: [true, "1", 5, "8", "13", 33]
	if result[1] != "1" || result[5] != 33 {
		t.Error("Unexpected sort by func result: ", result)
	}
}

// TestSortByFuncDescend validates sort by function implementation
func TestSortByFuncDescend(t *testing.T) {
	data := []interface{}{"1", 33, "8", "13", 5, true}

	result := SortByFunc(data, true, func(a, b interface{}) bool {
		return InterfaceToInt(a) < InterfaceToInt(b)
	})

	// expecting: [33, "13", "8", 5, "1", true]
	if result[1] != "13" || result[5] != true {
		t.Error("Unexpected sort by func result: ", result)
	}
}

// TestSortMapByIntAscend validates sort by map[string]int implementation in ascending order
func TestSortMapByIntAscend(t *testing.T) {
	data := map[string]int{
		"Blue":   2,
		"red":    7,
		"Green":  21,
		"yellow": 5,
		"black":  0,
		"Violet": 15,
	}

	result := SortByInt(data, false)
	if result[1].Value != 2 || result[3].Value != 7 || result[5].Key != "Green" {
		t.Error("Unexpected sort by func result: ", result)
	}

}

// TestSortMapByIntDescend validates sort by map[string]int implementation in descending order
func TestSortMapByIntDescend(t *testing.T) {
	data := map[string]int{
		"Blue":   2,
		"red":    7,
		"Green":  21,
		"yellow": 5,
		"black":  0,
		"Violet": 15,
	}

	result := SortByInt(data, true)
	if result[1].Value != 15 || result[3].Key != "yellow" || result[5].Key != "black" {
		t.Error("Unexpected sort by func result: ", result)
	}
}
