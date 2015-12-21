package giftcard

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// GetName returns name of shipping method
func (it *GiftcardShipping) GetName() string {
	return "Free Shipping"
}

// GetCode returns code of shipping method
func (it *GiftcardShipping) GetCode() string {
	return "giftcards"
}

// IsAllowed checks for method applicability
func (it *GiftcardShipping) IsAllowed(checkout checkout.InterfaceCheckout) bool {
	return true
}

// GetRates returns rates allowed by shipping method for a given checkout
func (it *GiftcardShipping) GetRates(currentCheckout checkout.InterfaceCheckout) []checkout.StructShippingRate {

	result := []checkout.StructShippingRate{}

	giftCardSkuElement := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftCardSKU))

	if cart := currentCheckout.GetCart(); cart != nil {
		for _, cartItem := range cart.GetItems() {

			cartProduct := cartItem.GetProduct()
			if cartProduct == nil {
				continue
			}

			cartProduct.ApplyOptions(cartItem.GetOptions())
			if !strings.Contains(cartProduct.GetSku(), giftCardSkuElement) {
				return result
			}
		}
	}

	result = []checkout.StructShippingRate{
		checkout.StructShippingRate{
			Code:  "freeshipping",
			Name:  "Freeshipping",
			Price: 0,
		}}

	return result
}
