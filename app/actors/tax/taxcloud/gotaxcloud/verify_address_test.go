package gotaxcloud_test

import (
	"testing"

	"github.com/ottemo/foundation/app/actors/tax/taxcloud/gotaxcloud"
)

func TestValidateAddressImpossible(t *testing.T) {
	_, err := testGateway.VerifyAddress(gotaxcloud.Address{
		Address1: "string",
		Address2: "string",
		City:     "string",
		State:    "string",
		Zip5:     "string",
		Zip4:     "string",
	})
	if err == nil {
		t.Fatal("expected ErrNumber to be not equal to '0'")
	}
}

func TestValidateAddressCorrection(t *testing.T) {
	result, err := testGateway.VerifyAddress(gotaxcloud.Address{
		Address1: "162 East Avenue",
		Address2: "Third Floor",
		City:     "Wilton",
		State:    "CT",
		Zip5:     "06851",
		Zip4:     "",
	})
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	if result.Address1 != "162 East Ave # 3" {
		t.Errorf("something changed in API; Address1 should be: '%s', got '%s'", "162 East Ave # 3", result.Address1)
	} else if result.Zip4 != "5715" {
		t.Errorf("something changed in API; Zip4 should be: '%s', got '%s'", "5715", result.Zip4)
	}
}
