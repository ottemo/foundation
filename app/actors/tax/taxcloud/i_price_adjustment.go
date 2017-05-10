// Implements app/models/checkout/InterfacePriceAdjustment
package taxcloud

import (
	"encoding/base64"
	"encoding/json"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/actors/tax/taxcloud/gotaxcloud"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
)

func (it *DefaultTaxCloud) GetName() string {
	return "Tax Cloud"
}

func (it *DefaultTaxCloud) GetCode() string {
	return "tax_cloud"
}

func (it *DefaultTaxCloud) GetPriority() []float64 {
	return []float64{ConstPriorityValue}
}

func (it *DefaultTaxCloud) Calculate(checkoutInstance checkout.InterfaceCheckout, currentPriority float64) []checkout.StructPriceAdjustment {
	result := []checkout.StructPriceAdjustment{}

	if currentPriority != ConstPriorityValue {
		// empty
		return result
	}

	config := env.GetConfig()
	if config == nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "de222130-0683-4b03-b3df-cbea18e819b2", "can't obtain config")
		// empty
		return result
	}

	if !(utils.InterfaceToBool(config.GetValue(ConstConfigPathEnabled))) {
		// empty
		return result
	}

	destinationAddress := getDestinationAddress(checkoutInstance)
	originAddressPtr, err := getOriginAddress()
	if err != nil {
		// empty
		return result
	}

	visitor := checkoutInstance.GetVisitor()
	var visitorID string
	if visitor != nil {
		visitorID = visitor.GetID()
	} else {
		emailBytes, err := json.Marshal(checkoutInstance.GetInfo("customer_email"))
		if err != nil {
			_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "0a3bd6bf-8efa-4cf5-a711-d0c740d15332", "unable to convert email")
			// empty
			return result
		}
		visitorID = base64.StdEncoding.EncodeToString(emailBytes)
	}

	cart := checkoutInstance.GetCart()
	cartID := cart.GetID()
	cartItems := getCartItems(checkoutInstance)

	gateway := gotaxcloud.NewGateway(
		utils.InterfaceToString(config.GetValue(ConstConfigPathAPILoginID)),
		utils.InterfaceToString(config.GetValue(ConstConfigPathAPIKey)))

	verifiedOriginAddressPtr, err := gateway.VerifyAddress(*originAddressPtr)
	if err != nil {
		// empty
		return result
	}

	verifiedDestinationAddressPtr, err := gateway.VerifyAddress(destinationAddress)
	if err != nil {
		// empty
		return result
	}

	lookupParams := gotaxcloud.LookupParams{
		Destination:       (*verifiedDestinationAddressPtr).Address,
		CartID:            cartID,
		CartItems:         cartItems,
		CustomerID:        visitorID,
		DeliveredBySeller: false,
		Origin:            (*verifiedOriginAddressPtr).Address,
	}
	lookupResponse, err := gateway.Lookup(lookupParams)
	if err != nil {
		// empty
		return result
	}

	var amount float64
	for _, cartItemInfo := range lookupResponse.CartItemsResponse {
		amount += cartItemInfo.TaxAmount
	}

	taxRate := checkout.StructPriceAdjustment{
		Code:      "TC",
		Name:      it.GetName(),
		Amount:    amount,
		IsPercent: false,
		Priority:  currentPriority,
		Labels:    []string{checkout.ConstLabelTax},
	}
	result = append(result, taxRate)

	return result
}

func getDestinationAddress(checkoutInstance checkout.InterfaceCheckout) gotaxcloud.Address {
	address := checkoutInstance.GetShippingAddress()

	return gotaxcloud.Address{
		Address1: address.GetAddressLine1(),
		Address2: address.GetAddressLine2(),
		City:     address.GetCity(),
		State:    address.GetState(),
		Zip5:     address.GetZipCode(),
	}
}

func getOriginAddress() (*gotaxcloud.Address, error) {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "de222130-0683-4b03-b3df-cbea18e819b2", "can't obtain config")
		// empty
		return nil, err
	}

	return &gotaxcloud.Address{
		Address1: utils.InterfaceToString(config.GetValue(app.ConstConfigPathStoreAddressline1)),
		Address2: utils.InterfaceToString(config.GetValue(app.ConstConfigPathStoreAddressline2)),
		City:     utils.InterfaceToString(config.GetValue(app.ConstConfigPathStoreCity)),
		State:    utils.InterfaceToString(config.GetValue(app.ConstConfigPathStoreState)),
		Zip5:     utils.InterfaceToString(config.GetValue(app.ConstConfigPathStoreZip)),
	}, nil
}

func getCartItems(checkoutInstance checkout.InterfaceCheckout) []gotaxcloud.CartItem {
	result := []gotaxcloud.CartItem{}

	itemsGrandTotal := 0.0
	cart := checkoutInstance.GetCart()
	cartItems := cart.GetItems()
	for idx, cartItem := range cartItems {
		// Price should have discounts applied
		grandTotal := checkoutInstance.GetItemSpecificTotal(cartItem.GetIdx(), checkout.ConstLabelGrandTotal)
		itemsGrandTotal += grandTotal
		price := grandTotal / float64(cartItem.GetQty())
		result = append(result, gotaxcloud.CartItem{
			Index:  idx, //idx should be 0-based
			ItemID: cartItem.GetProductID(),
			Price:  price,
			Qty:    cartItem.GetQty(),
			TIC:    0,
		})
	}

	discounts := checkoutInstance.GetDiscounts()
	perCartDiscount := 0.0
	for _, discount := range discounts {
		if discount.PerItem == nil {
			// negative value
			perCartDiscount += discount.Amount
		}
	}

	if perCartDiscount != 0 {
		// negative value
		perItemDiscountPercent := perCartDiscount / itemsGrandTotal
		discountedGrandTotal := 0.0

		for i := range result {
			result[i].Price = result[i].Price * (1 + perItemDiscountPercent)
			result[i].Price = utils.RoundPrice(result[i].Price)
			discountedGrandTotal += result[i].Price * float64(result[i].Qty)
		}

		//negative
		grandTotalDiff := (itemsGrandTotal + perCartDiscount) - discountedGrandTotal
		if grandTotalDiff != 0 {
			for i := range result {
				if result[i].Price > -grandTotalDiff {
					if result[i].Qty == 1 {
						result[i].Price = utils.RoundPrice(result[i].Price + grandTotalDiff)
						break
					} else {
						lastIndex := len(result) - 1
						lastItem := result[lastIndex]
						result = append(result, gotaxcloud.CartItem{
							Index:  lastItem.Index + 1,
							ItemID: result[i].ItemID,
							Price:  utils.RoundPrice(result[i].Price + grandTotalDiff),
							Qty:    1,
							TIC:    0,
						})
						result[i].Qty -= 1
						break
					}
				}
			}
		}
	}

	return result
}
