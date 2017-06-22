// Package fedex is a FedEx implementation of shipping method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package fedex

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	ConstShippingMethods = map[string]string{
		"EUROPE_FIRST_INTERNATIONAL_PRIORITY": "Europe First Priority",
		"FEDEX_1_DAY_FREIGHT":                 "1 Day Freight",
		"FEDEX_2_DAY_FREIGHT":                 "2 Day Freight",
		"FEDEX_2_DAY":                         "2 Day",
		"FEDEX_2_DAY_AM":                      "2 Day AM",
		"FEDEX_3_DAY_FREIGHT":                 "3 Day Freight",
		"FEDEX_EXPRESS_SAVER":                 "Express Saver",
		"FEDEX_GROUND":                        "Ground",
		"FIRST_OVERNIGHT":                     "First Overnight",
		"GROUND_HOME_DELIVERY":                "Home Delivery",
		"INTERNATIONAL_ECONOMY":               "International Economy",
		"INTERNATIONAL_ECONOMY_FREIGHT":       "Intl Economy Freight",
		"INTERNATIONAL_FIRST":                 "International First",
		"INTERNATIONAL_GROUND":                "International Ground",
		"INTERNATIONAL_PRIORITY":              "International Priority",
		"INTERNATIONAL_PRIORITY_FREIGHT":      "Intl Priority Freight",
		"PRIORITY_OVERNIGHT":                  "Priority Overnight",
		"SMART_POST":                          "Smart Post",
		"STANDARD_OVERNIGHT":                  "Standard Overnight",
		"FEDEX_FREIGHT":                       "Freight",
		"FEDEX_NATIONAL_FREIGHT":              "National Freight",
	}

	ConstShippingDropoff = map[string]string{
		"REGULAR_PICKUP":          "Regular Pickup",
		"REQUEST_COURIER":         "Request Courier",
		"DROP_BOX":                "Drop Box",
		"BUSINESS_SERVICE_CENTER": "Business Service Center",
		"STATION":                 "Station",
	}

	ConstShippingPackaging = map[string]string{
		"FEDEX_ENVELOPE": "FedEx Envelope",
		"FEDEX_PAK":      "FedEx Pak",
		"FEDEX_BOX":      "FedEx Box",
		"FEDEX_TUBE":     "FedEx Tube",
		"FEDEX_10KG_BOX": "FedEx 10kg Box",
		"FEDEX_25KG_BOX": "FedEx 25kg Box",
		"YOUR_PACKAGING": "Your Packaging",
	}
)

// Package global constants
const (
	ConstShippingCode = "fedex"
	ConstShippingName = "fedex"

	ConstConfigPathGroup = "shipping.fedex"

	ConstConfigPathEnabled = "shipping.fedex.enabled"
	ConstConfigPathTitle   = "shipping.fedex.title"

	ConstConfigPathGateway  = "shipping.fedex.gateway"
	ConstConfigPathKey      = "shipping.fedex.key"
	ConstConfigPathPassword = "shipping.fedex.password"
	ConstConfigPathNumber   = "shipping.fedex.number"
	ConstConfigPathMeter    = "shipping.fedex.meter"

	ConstConfigPathDefaultWeight  = "shipping.fedex.default_weight"
	ConstConfigPathAllowedMethods = "shipping.fedex.allowed_methods"
	ConstConfigPathPackaging      = "shipping.fedex.packaging"
	ConstConfigPathAllowCountries = "shipping.fedex.allow_countries"
	ConstConfigPathDropoff        = "shipping.fedex.dropoff"
	ConstConfigPathResidential    = "shipping.fedex.residential"

	ConstConfigPathDebugLog = "shipping.fedex.debug_log"

	ConstErrorModule = "shipping/fedex"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// FedEx is a implementer of InterfaceShippingMethod for "FedEx" shipping method
type FedEx struct{}
