package taxcloud_test

import (
	"testing"

	"github.com/ottemo/foundation/app/actors/tax/taxcloud"
)

func TestCaptured(t *testing.T) {
	var unique = getUniqueStr()
	var cartID = "AUTHORIZE-CART-UUID" + unique
	var customerID = "AUTHORIZE-CUSTOMER-UUID" + unique
	var orderID = "AUTHORIZE-ORDER-UUID" + unique

	var err error

	_, err = doCorrectLookup(cartID, customerID)
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	_, err = doAuthorized(cartID, customerID, orderID)
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	_, err = testGateway.Captured(taxcloud.CapturedRequestType{
		OrderID: orderID,
	})
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
}

func TestCapturedWithError(t *testing.T) {
	result, err := testGateway.Captured(taxcloud.CapturedRequestType{
		OrderID: "ORDER-UUID"+getUniqueStr(),
	})
	if err == nil {
		t.Fatalf("unknown transaction should raise error, but '%s'", result)
	}
}

