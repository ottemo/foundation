// Package usps is a USPS implementation of shipping method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package usps

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstShippingCode = "usps"
	ConstShippingName = "USPS"

	ConstRemoveRateNameTags = true

	ConstHTTPEndpoint  = "http://production.shippingapis.com/ShippingAPI.dll"
	ConstHTTPSEndpoint = "https://secure.shippingapis.com/ShippingAPI.dll"

	ConstConfigPathGroup             = "shipping.usps"
	ConstConfigPathEnabled           = "shipping.usps.enabled"
	ConstConfigPathUser              = "shipping.usps.userid"
	ConstConfigPathOriginZip         = "shipping.usps.zip"
	ConstConfigPathContainer         = "shipping.usps.container"
	ConstConfigPathSize              = "shipping.usps.size"
	ConstConfigPathDefaultDimensions = "shipping.usps.default_dimensions"
	ConstConfigPathDefaultWeight     = "shipping.usps.default_weight"
	ConstConfigPathAllowedMethods    = "shipping.usps.allowed_methods"
	ConstConfigPathAllowCountries    = "shipping.usps.allow_countries"
	ConstConfigPathDebugLog          = "shipping.usps.debug_log"

	ConstErrorModule = "shipping/usps"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// Package global variables
var (
	ConstShippingMethods = map[string]string{
		"1":      "Priority Mail",
		"2":      "Priority Mail Express Hold For Pickup",
		"3":      "Priority Mail Express",
		"4":      "Standard Post",
		"6":      "Media Mail",
		"7":      "Library Mail",
		"13":     "Priority Mail Express Flat Rate Envelope",
		"15":     "First-Class Mail Large Postcards",
		"16":     "Priority Mail Flat Rate Envelope",
		"17":     "Priority Mail Medium Flat Rate Box",
		"22":     "Priority Mail Large Flat Rate Box",
		"23":     "Priority Mail Express Sunday/Holiday Delivery",
		"25":     "Priority Mail Express Sunday/Holiday Delivery Flat Rate Envelope",
		"27":     "Priority Mail Express Flat Rate Envelope Hold For Pickup",
		"28":     "Priority Mail Small Flat Rate Box",
		"29":     "Priority Mail Padded Flat Rate Envelope",
		"30":     "Priority Mail Express Legal Flat Rate Envelope",
		"31":     "Priority Mail Express Legal Flat Rate Envelope Hold For Pickup",
		"32":     "Priority Mail Express Sunday/Holiday Delivery Legal Flat Rate Envelope",
		"33":     "Priority Mail Hold For Pickup",
		"34":     "Priority Mail Large Flat Rate Box Hold For Pickup",
		"35":     "Priority Mail Medium Flat Rate Box Hold For Pickup",
		"36":     "Priority Mail Small Flat Rate Box Hold For Pickup",
		"37":     "Priority Mail Flat Rate Envelope Hold For Pickup",
		"38":     "Priority Mail Gift Card Flat Rate Envelope",
		"39":     "Priority Mail Gift Card Flat Rate Envelope Hold For Pickup",
		"40":     "Priority Mail Window Flat Rate Envelope",
		"41":     "Priority Mail Window Flat Rate Envelope Hold For Pickup",
		"42":     "Priority Mail Small Flat Rate Envelope",
		"43":     "Priority Mail Small Flat Rate Envelope Hold For Pickup",
		"44":     "Priority Mail Legal Flat Rate Envelope",
		"45":     "Priority Mail Legal Flat Rate Envelope Hold For Pickup",
		"46":     "Priority Mail Padded Flat Rate Envelope Hold For Pickup",
		"47":     "Priority Mail Regional Rate Box A",
		"48":     "Priority Mail Regional Rate Box A Hold For Pickup",
		"49":     "Priority Mail Regional Rate Box B",
		"50":     "Priority Mail Regional Rate Box B Hold For Pickup",
		"53":     "First-Class Package Service Hold For Pickup",
		"55":     "Priority Mail Express Flat Rate Boxes",
		"56":     "Priority Mail Express Flat Rate Boxes Hold For Pickup",
		"57":     "Priority Mail Express Sunday/Holiday Delivery Flat Rate Boxes",
		"58":     "Priority Mail Regional Rate Box C",
		"59":     "Priority Mail Regional Rate Box C Hold For Pickup",
		"61":     "First-Class Package Service",
		"62":     "Priority Mail Express Padded Flat Rate Envelope",
		"63":     "Priority Mail Express Padded Flat Rate Envelope Hold For Pickup",
		"64":     "Priority Mail Express Sunday/Holiday Delivery Padded Flat Rate Envelope",
		"0_FCLE": "First-Class Mail Large Envelope",
		"0_FCL":  "First-Class Mail Letter",
		"0_FCP":  "First-Class Mail Parcel",
		"0_FCPC": "First-Class Mail Postcards",
		"INT_1":  "Priority Mail Express International",
		"INT_2":  "Priority Mail International",
		"INT_4":  "Global Express Guaranteed (GXG)",
		"INT_5":  "Global Express Guaranteed Document",
		"INT_6":  "Global Express Guaranteed Non-Document Rectangular",
		"INT_7":  "Global Express Guaranteed Non-Document Non-Rectangular",
		"INT_8":  "Priority Mail International Flat Rate Envelope",
		"INT_9":  "Priority Mail International Medium Flat Rate Box",
		"INT_10": "Priority Mail Express International Flat Rate Envelope",
		"INT_11": "Priority Mail International Large Flat Rate Box",
		"INT_12": "USPS GXG Envelopes",
		"INT_13": "First-Class Mail International Letter",
		"INT_14": "First-Class Mail International Large Envelope",
		"INT_15": "First-Class Package International Service",
		"INT_16": "Priority Mail International Small Flat Rate Box",
		"INT_17": "Priority Mail Express International Legal Flat Rate Envelope",
		"INT_18": "Priority Mail International Gift Card Flat Rate Envelope",
		"INT_19": "Priority Mail International Window Flat Rate Envelope",
		"INT_20": "Priority Mail International Small Flat Rate Envelope",
		"INT_21": "First-Class Mail International Postcard",
		"INT_22": "Priority Mail International Legal Flat Rate Envelope",
		"INT_23": "Priority Mail International Padded Flat Rate Envelope",
		"INT_24": "Priority Mail International DVD Flat Rate priced box",
		"INT_25": "Priority Mail International Large Video Flat Rate priced box",
		"INT_26": "Priority Mail Express International Flat Rate Boxes",
		"INT_27": "Priority Mail Express International Padded Flat Rate Envelope",
	}
)

// USPS is a implementer of InterfaceShippingMethod for a "USPS" shipping method
type USPS struct{}
