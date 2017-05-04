package taxcloud_test

import (
	"testing"

	"github.com/ottemo/foundation/app/actors/tax/taxcloud"
)

func TestLookup(t *testing.T) {
	result, err := doCorrectLookup("CART-UUID", "CUSTOMER-UUID")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	if result == nil {
		t.Fatal("result is empty")
	} else if result.CartItemsResponse == nil {
		t.Fatal("expected info for cart items")
	} else if len(result.CartItemsResponse) != 2 {
		t.Fatal("expected 2 cart items")
	} else if result.CartItemsResponse[0].TaxAmount != 0.091575 {
		t.Fatalf("something changed in API; expected tax for item [0] == '%s', got '%s'", 0.091575, result.CartItemsResponse[0].TaxAmount)
	}
}

func doCorrectLookup(cartID, customerID string) (*taxcloud.LookupResponse, error) {
	return testGateway.Lookup(taxcloud.LookupParams{
		CartID: cartID,
		CustomerID: customerID,
		DeliveredBySeller: false,
		Origin: taxcloud.Address{
			Address1: "11 Crandall Dr",
			Address2: "",
			City: "East Brunswick",
			State: "NJ",
			Zip5: "08816",
			Zip4: "5613",
		},
		Destination: taxcloud.Address{
			Address1: "9650 Ensworth St",
			Address2: "",
			City: "Las Vegas",
			State: "NV",
			Zip5: "89123",
			Zip4: "6545",
		},
		CartItems: []taxcloud.CartItem{
			{
				Index:0,
				ItemID:"ITEM-UUID-0",
				Price: 1.11,
				Qty: 1,
				TIC: 0,
			},
			{
				Index:1,
				ItemID:"ITEM-UUID-2",
				Price: 2.22,
				Qty: 2,
				TIC: 0,
			},
		},
	})
}

