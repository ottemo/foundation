package composer

import (
	"testing"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/utils"
)

type testIObject struct { data map[string]interface{} }

func (it *testIObject) Get(attribute string) interface{} {
	if value, present := it.data[attribute]; present {
		return value
	}
	return nil
}

func (it *testIObject) Set(attribute string, value interface{}) error {
	it.data[attribute] = value
	return nil
}

func (it *testIObject) FromHashMap(hashMap map[string]interface{}) error {
	return nil
}

func (it *testIObject) ToHashMap() map[string]interface{} {
	return it.data
}

func (it *testIObject) GetAttributesInfo() []models.StructAttributeInfo {
	return []models.StructAttributeInfo{}
}

func TestOperations(tst *testing.T) {

	var object models.InterfaceObject = &testIObject{ data: map[string]interface{} {
		"sku":   "test_product",
		"name":  "Test Product",
		"price": 1.1,
	}}


	tst.Log(object.Get("sku"))
	input := map[string]interface{} {
		"a": 10,
		"b": "test",
		"c": 3.14,
		"d": object,
	}

	rules, err := utils.DecodeJSONToStringKeyMap(`{
		"a": "10",
		"b": "test",
		"c": 3.140,
		"d": {
			"price": 1.1,
			"sku": "test_product"
		}
	}`)
	if err != nil {
		tst.Errorf("JSON decode fail: %v", err)
	}

	result, err := GetComposer().Validate(input, rules)
	if err != nil {
		tst.Errorf("Validation fail", err)
	} else if !result {
		tst.Errorf("Validation fail")
	}
}