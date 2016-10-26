package flatweight

import (
	"testing"
)

// validateAndApplyRates
var ValidateAndApplyTests = []struct {
	rawRates        interface{}
	startRateGlobal Rates
	finalRateGlobal Rates
	expectedErr     bool
}{
	// Positive test
	{"", make(Rates, 0), make(Rates, 0), false},
	{"[]", make(Rates, 0), make(Rates, 0), false},
	{"{}", make(Rates, 0), make(Rates, 0), false},
	{
		// Starting with no stored rates
		"[{\"title\": \"Standard Shipping\",\"code\": \"std_1\",\"price\": 1.99,\"weight_from\": 0.0,\"weight_to\": 5.0}]",
		make(Rates, 0),
		[]Rate{Rate{"Standard Shipping", "std_1", 1.99, 0.0, 5.0, "", ""}},
		false,
	},
	{
		// Updating the price of a stored rate
		"[{\"title\": \"Standard Shipping\",\"code\": \"std_1\",\"price\": 5.99,\"weight_from\": 0.0,\"weight_to\": 5.0}]",
		[]Rate{Rate{"Standard Shipping", "std_1", 1.99, 0.0, 5.0, "", ""}},
		[]Rate{Rate{"Standard Shipping", "std_1", 5.99, 0.0, 5.0, "", ""}},
		false,
	},
	{
		// Updating the label of a stored rate
		"[{\"title\": \"Shipping\",\"code\": \"std_1\",\"price\": 1.99,\"weight_from\": 0.0,\"weight_to\": 5.0}]",
		[]Rate{Rate{"Standard Shipping", "std_1", 1.99, 0.0, 5.0, "", ""}},
		[]Rate{Rate{"Shipping", "std_1", 1.99, 0.0, 5.0, "", ""}},
		false,
	},

	// Fail tests
	// 1. broken json
	{
		"[{broken: json\" : 1.23}]",
		make(Rates, 0),
		make(Rates, 0),
		true,
	},

	// 2. missing fields
	{
		"[{\"title\": \"Standard Shipping\",\"price\": 1.99,\"weight_from\": 0.0,\"weight_to\": 5.0}]",
		make(Rates, 0),
		make(Rates, 0),
		true,
	},
}

func TestValidateAndApplyRates(t *testing.T) {
	// Testing a variety of inputs
	for _, vt := range ValidateAndApplyTests {
		// reset rates global so we can perform our comparison
		rates = vt.startRateGlobal

		_, err := validateAndApplyRates(vt.rawRates)

		// error doesn't match expectations
		hasErr := err != nil
		if hasErr != vt.expectedErr {
			t.Errorf("validateAndApplyRates(%v) err: %v, expectedErr: %v",
				vt.rawRates, err, vt.expectedErr)
		}

		// did we set the global rates variable properly?
		if !equalRates(rates, vt.finalRateGlobal) {
			t.Errorf("validateAndApplyRates(%v) expectedRates:%v, actualRates:%v",
				vt.rawRates, vt.finalRateGlobal, rates)
		}
	}
}

func equalRates(a, b Rates) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
