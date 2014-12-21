package flatrate

import (
	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetName returns name of shipping method
func (it *ShippingMethod) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathName))
}

// GetCode returns code of shipping method
func (it *ShippingMethod) GetCode() string {
	return ConstShippingCode
}

// IsAllowed checks for method applicability
func (it *ShippingMethod) IsAllowed(checkout checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEnabled))
}

// GetRates returns rates allowed by shipping method for a given checkout
func (it *ShippingMethod) GetRates(checkoutObject checkout.InterfaceCheckout) []checkout.StructShippingRate {

	return []checkout.StructShippingRate{
		checkout.StructShippingRate{
			Code:  "default",
			Name:  utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathName)),
			Price: utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathAmount)),
		}}
}
