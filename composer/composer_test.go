package composer

import (
	"testing"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/utils"
)

func TestOperations(tst *testing.T) {

	productItem, err := product.GetProductModel()
	if err != nil && productItem != nil {
		productItem.New()
		productItem.Set("sku", "test_product")
		productItem.Set("name", "Test Product")
		productItem.Set("price", 1.1)

		tst.Errorf("Product model undefined %v", err)
	}

	input := map[string]interface{} {
		"a": 10,
		"b": "test",
		"c": 3.14,
		"d": productItem,
	}

	rules, err := utils.DecodeJSONToStringKeyMap(`{
		"a": 10
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