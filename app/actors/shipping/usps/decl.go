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
	ConstConfigPathDebugLog          = "shipping.usps.debug_log"

	ConstErrorModule = "shipping/usps"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// USPS is a implementer of InterfaceShippingMethod for a "USPS" shipping method
type USPS struct{}
