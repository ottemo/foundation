package shipstation

import (
	"testing"

	"github.com/ottemo/foundation/app/actors/order"
	_ "github.com/ottemo/foundation/app/actors/visitor" // required to initialize Visitor Address Model
)

// TestBuildItemReturnsNoAdjustments tests buildItem function.
// In following example go 1.4.3 and 1.6.3 in unfixed version returns adjustment
// item with price -5.684341886080802e-14. After fix there should not be any adjustment items.
func TestBuildItemReturnsNoAdjustments(t *testing.T) {
	const orderID = "5846cbbf7720ae713751e355"

	var orderObj = &order.DefaultOrder{}
	if err := orderObj.SetID(orderID); err != nil {
		t.Error(err)
	}
	orderObj.TaxAmount = 37.34
	orderObj.ShippingAmount = 0
	orderObj.GrandTotal = 426.34
	orderObj.CustomInfo = map[string]interface{}{
		"calculation": map[string]interface{}{
			"0": map[string]interface{}{"SP": 0, "ST": 389, "T": 37.34, "GT": 426.34},
			"1": map[string]interface{}{"GT": 280, "ST": 280},
			"2": map[string]interface{}{"GT": 32, "ST": 32},
			"3": map[string]interface{}{"GT": 21, "ST": 21},
			"4": map[string]interface{}{"GT": 21, "ST": 21},
			"5": map[string]interface{}{"GT": 35, "ST": 35},
		},
	}

	var allOrderItems []map[string]interface{}
	allOrderItems = append(allOrderItems, map[string]interface{}{
		"idx":      1,
		"order_id": orderID,
		"price":    140,
		"qty":      2,
		"sku":      "sku-01",
		"name":     "name-01",
	})
	allOrderItems = append(allOrderItems, map[string]interface{}{
		"idx":      2,
		"order_id": orderID,
		"price":    16,
		"qty":      2,
		"sku":      "sku-02",
		"name":     "name-02",
	})
	allOrderItems = append(allOrderItems, map[string]interface{}{
		"idx":      3,
		"order_id": orderID,
		"price":    21,
		"qty":      1,
		"sku":      "sku-03",
		"name":     "name-03",
	})
	allOrderItems = append(allOrderItems, map[string]interface{}{
		"idx":      4,
		"order_id": orderID,
		"price":    21,
		"qty":      1,
		"sku":      "sku-03",
		"name":     "name-03",
	})
	allOrderItems = append(allOrderItems, map[string]interface{}{
		"idx":      5,
		"order_id": orderID,
		"price":    35,
		"qty":      1,
		"sku":      "sku-04",
		"name":     "name-04",
	})

	var builtOrder = buildItem(orderObj, allOrderItems)

	for _, item := range builtOrder.Items {
		if item.Adjustment {
			t.Error("Should be no adjustments: UnitPrice =", item.UnitPrice)
		}
	}
}
