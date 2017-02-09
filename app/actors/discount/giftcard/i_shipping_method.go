package giftcard

import (
	"strings"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetName returns name of shipping method
func (it *Shipping) GetName() string {
	return "No Shipping"
}

// GetCode returns code of shipping method
func (it *Shipping) GetCode() string {
	return "giftcards"
}

// IsAllowed checks for method applicability
func (it *Shipping) IsAllowed(checkout checkout.InterfaceCheckout) bool {
	return true
}

// GetRates returns rates allowed by shipping method for a given checkout
func (it *Shipping) GetRates(currentCheckout checkout.InterfaceCheckout) []checkout.StructShippingRate {

	result := []checkout.StructShippingRate{}

	giftCardSkuElement := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftCardSKU))

	if cart := currentCheckout.GetCart(); cart != nil {
		for _, cartItem := range cart.GetItems() {

			cartProduct := cartItem.GetProduct()
			if cartProduct == nil {
				continue
			}

			if err := cartProduct.ApplyOptions(cartItem.GetOptions()); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5fcd3354-17b5-43ad-946b-0146a6e0017a", err.Error())
			}
			if !strings.Contains(cartProduct.GetSku(), giftCardSkuElement) {
				return result
			}
		}
	}

	result = []checkout.StructShippingRate{
		checkout.StructShippingRate{
			Code:  "freeshipping",
			Name:  "GiftCards",
			Price: 0,
		}}

	return result
}

// GetAllRates returns all the shippmeng method rates available in the system.
func (it Shipping) GetAllRates() []checkout.StructShippingRate {

	result := []checkout.StructShippingRate{
		checkout.StructShippingRate{
			Code:  "freeshipping",
			Name:  "GiftCards",
			Price: 0,
		},
	}

	return result
}
