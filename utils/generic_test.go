package utils

import (
	"testing"
	"fmt"
)

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

func TestStrToSnakeCase(t *testing.T) {

	str := "Product Size "
	if StrToSnakeCase(str) != "product_size" {
		t.Error("case 1 fail")
	}

	str = "ProductSize 1"
	if StrToSnakeCase(str) != "product_size_1" {
		t.Error("case 2 fail")
	}

	str = "-101"
	if StrToSnakeCase(str) != "-101" {
		t.Error("case 3 fail")
	}

	str = " - 101"
	if StrToSnakeCase(str) != "-_101" {
		t.Error("case 4 fail")
	}

	str = " - 101 Discount Amount"

	if StrToSnakeCase(str) != "-_101_discount_amount" {
		t.Error("case 5 fail")
	}

	str = "subtract - 101 from Discount Amount"
	if StrToSnakeCase(str) != "subtract_-_101_from_discount_amount" {
		t.Error("case 6 fail")
	}

	str = ";LARGE"
	if StrToSnakeCase(str) != "large" {
		t.Error("case 8 fail")
	}

	str = "XLarge"
	if StrToSnakeCase(str) != "x_large" {
		t.Error("case 9 fail")
	}

	str = "X-Large"
	if StrToSnakeCase(str) != "x-large" {
		t.Error("case 10 fail")
	}

	str = "$20 X-LARGE"
	if StrToSnakeCase(str) != "$20_x-large" {
		t.Error("case 11 fail")
	}

	str = "Size (*^*%^@XLARGE"
	if StrToSnakeCase(str) != "size_xlarge" {
		t.Error("case 12 fail")
	}

	str = "     Size: XLARGE + @'-3'Num of  *&*&&^^^^()(##A   ; "
	if StrToSnakeCase(str) != "size_xlarge_-3_num_of_a" {
		t.Error("case 13 fail")
	}
}

func TestStrToCamelCase(t *testing.T) {

	str := "product_size_xlarge"
	if StrToCamelCase(str) != "productSizeXlarge" {
		t.Error("case 1 fail")
	}

	str = "subtract_-_101_discount_amount"
	if StrToCamelCase(str) != "subtract-101DiscountAmount" {
		t.Error("case 2 fail")
	}
}

func TestMapSetPathValue(t *testing.T) {
	x := map[string]interface{} { }

	if e := MapSetPathValue(x, "a.b.c.d.e.f", "something", false); e != nil {
		t.Error(e)
	}

	if e := MapSetPathValue(x, "a.g", "something else", false); e != nil {
		t.Error(e)
	}


	if value, _ := MapGetPathValue(x, "a.b.c.d.e.f"); value != "something" {
		t.Error(fmt.Sprintf("test case 1 fail: %v != 'something'", value))
	}

	if e := MapSetPathValue(x, "a.b.c.d.e.f", nil, true); e != nil {
		t.Error(e)
	}

	if value, _ := MapGetPathValue(x, "a.b.c.d.e.f"); value != nil {
		t.Error(fmt.Sprintf("test case 2 fail: %v != nil", value))
	}

	if value, _ := MapGetPathValue(x, "a.g"); value != "something else" {
		t.Error(fmt.Sprintf("test case 3 fail: %v != 'something else'", value))
	}
}
