package utils

import "testing"

func TestMatchMapAValuesToMapB(t *testing.T) {

	B := map[string]interface{}{
		"color":   "red",
		"size":    "XL",
		"enabled": true,
		"price":   11.2,
		"qty":     15,
		"tags":    []interface{}{"cool", "best"},

		"extra": map[string]interface{}{
			"x": "1",
			"y": 2,
			"z": true,
		},
	}

	A := map[string]interface{}{}
	if MatchMapAValuesToMapB(A, B) != true {
		t.Error("case 1 fail")
	}

	A = map[string]interface{}{"color": "red", "size": "XL"}
	if MatchMapAValuesToMapB(A, B) != true {
		t.Error("case 2 fail")
	}

	A = map[string]interface{}{"size": "XL", "qty": 15, "price": 11.2, "enabled": true}
	if MatchMapAValuesToMapB(A, B) != true {
		t.Error("case 3 fail")
	}

	A = map[string]interface{}{"size": "XL", "qty": 15, "price": 11.3, "enabled": true}
	if MatchMapAValuesToMapB(A, B) != false {
		t.Error("case 4 fail")
	}

	A = map[string]interface{}{"size": "XL", "qty": 15, "price": 11.2, "enabled": false}
	if MatchMapAValuesToMapB(A, B) != false {
		t.Error("case 5 fail")
	}

	A = map[string]interface{}{"size": "XL", "extra": map[string]interface{}{"x": "1", "z": true}}
	if MatchMapAValuesToMapB(A, B) != true {
		t.Error("case 6 fail")
	}

	A = map[string]interface{}{"extra": map[string]interface{}{"x": "1", "y": 2, "z": true}}
	if MatchMapAValuesToMapB(A, B) != true {
		t.Error("case 7 fail")
	}

	A = map[string]interface{}{"extra": map[string]interface{}{"z": false}}
	if MatchMapAValuesToMapB(A, B) != false {
		t.Error("case 8 fail")
	}

	A = map[string]interface{}{"tags": []interface{}{"best"}}
	if MatchMapAValuesToMapB(A, B) != true {
		t.Error("case 9 fail")
	}

	A = map[string]interface{}{"tags": []interface{}{"super"}}
	if MatchMapAValuesToMapB(A, B) != false {
		t.Error("case 10 fail")
	}
}
