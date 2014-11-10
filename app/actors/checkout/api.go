package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("checkout", "GET", "info", restCheckoutInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout", "GET", "payment/methods", restCheckoutPaymentMethods)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout", "GET", "shipping/methods", restCheckoutShippingMethods)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/shipping/address", restCheckoutSetShippingAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/billing/address", restCheckoutSetBillingAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/payment/method/:method", restCheckoutSetPaymentMethod)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/shipping/method/:method/:rate", restCheckoutSetShippingMethod)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout", "GET", "submit", restSubmit)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "submit", restSubmit)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to get current checkout process status
func restCheckoutInfo(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := map[string]interface{}{
		"billing_address":  nil,
		"shipping_address": nil,

		"payment_method_name": nil,
		"payment_method_code": nil,

		"shipping_method_name": nil,
		"shipping_method_code": nil,

		"shipping_rate":   nil,
		"shipping_amount": nil,

		"discounts":       nil,
		"discount_amount": nil,

		"taxes":      nil,
		"tax_amount": nil,

		"subtotal":   nil,
		"grandtotal": nil,
	}

	if billingAddress := currentCheckout.GetBillingAddress(); billingAddress != nil {
		result["billing_address"] = billingAddress.ToHashMap()
	}

	if shippingAddress := currentCheckout.GetShippingAddress(); shippingAddress != nil {
		result["shipping_address"] = shippingAddress.ToHashMap()
	}

	if paymentMethod := currentCheckout.GetPaymentMethod(); paymentMethod != nil {
		result["payment_method_name"] = paymentMethod.GetName()
		result["payment_method_code"] = paymentMethod.GetCode()
	}

	if shippingMethod := currentCheckout.GetShippingMethod(); shippingMethod != nil {
		result["shipping_method_name"] = shippingMethod.GetName()
		result["shipping_method_code"] = shippingMethod.GetCode()
	}

	if shippingRate := currentCheckout.GetShippingRate(); shippingRate != nil {
		result["shipping_rate"] = shippingRate
		result["shipping_amount"] = shippingRate.Price
	}

	if checkoutCart := currentCheckout.GetCart(); checkoutCart != nil {
		result["subtotal"] = checkoutCart.GetSubtotal()
	}

	result["discount_amount"], result["discounts"] = currentCheckout.GetDiscounts()

	result["tax_amount"], result["taxes"] = currentCheckout.GetTaxes()

	result["grandtotal"] = currentCheckout.GetGrandTotal()

	return result, nil
}

// WEB REST API function to get possible payment methods for checkout
func restCheckoutPaymentMethods(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	type ResultValue struct {
		Name string
		Code string
		Type string
	}
	result := make([]ResultValue, 0)

	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.IsAllowed(currentCheckout) {
			result = append(result, ResultValue{Name: paymentMethod.GetName(), Code: paymentMethod.GetCode(), Type: paymentMethod.GetType()})
		}
	}

	return result, nil
}

// WEB REST API function to get possible shipping methods for checkout
func restCheckoutShippingMethods(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	type ResultValue struct {
		Name  string
		Code  string
		Rates []checkout.T_ShippingRate
	}
	result := make([]ResultValue, 0)

	for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
		if shippingMethod.IsAllowed(currentCheckout) {
			result = append(result, ResultValue{Name: shippingMethod.GetName(), Code: shippingMethod.GetCode(), Rates: shippingMethod.GetRates(currentCheckout)})
		}
	}

	return result, nil
}

// internal function for  restCheckoutSetShippingAddress() and restCheckoutSetBillingAddress()
func checkoutObtainAddress(params *api.T_APIHandlerParams) (visitor.I_VisitorAddress, error) {

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if addressId, present := reqData["id"]; present {

		// Address id was specified - trying to load
		visitorAddress, err := visitor.LoadVisitorAddressByID(utils.InterfaceToString(addressId))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		currentVisitorId := utils.InterfaceToString(params.Session.Get(visitor.SessionKeyVisitorID))
		if visitorAddress.GetVisitorId() != currentVisitorId {
			return nil, env.ErrorNew("wrong address id")
		}

		return visitorAddress, nil
	} else {

		// supposedly address data was specified
		visitorAddressModel, err := visitor.GetVisitorAddressModel()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		for attribute, value := range reqData {
			err := visitorAddressModel.Set(attribute, value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}

		visitorId := utils.InterfaceToString(params.Session.Get(visitor.SessionKeyVisitorID))
		visitorAddressModel.Set("visitor_id", visitorId)

		err = visitorAddressModel.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		return visitorAddressModel, nil
	}
}

// WEB REST API function to set shipping address
func restCheckoutSetShippingAddress(params *api.T_APIHandlerParams) (interface{}, error) {
	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	address, err := checkoutObtainAddress(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = currentCheckout.SetShippingAddress(address)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return address.ToHashMap(), nil
}

// WEB REST API function to set billing address
func restCheckoutSetBillingAddress(params *api.T_APIHandlerParams) (interface{}, error) {
	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	address, err := checkoutObtainAddress(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = currentCheckout.SetBillingAddress(address)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return address.ToHashMap(), nil
}

// WEB REST API function to set payment method
func restCheckoutSetPaymentMethod(params *api.T_APIHandlerParams) (interface{}, error) {

	eventData := make(map[string]interface{})
	eventData["sessionId"] = params.Session.GetId()
	env.Event("api.reachedCheckout", eventData)

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// looking for payment method
	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.GetCode() == params.RequestURLParams["method"] {
			if paymentMethod.IsAllowed(currentCheckout) {

				// updating checkout payment method
				err := currentCheckout.SetPaymentMethod(paymentMethod)
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}

				// checking for additional info
				contentValues, _ := api.GetRequestContentAsMap(params)
				for key, value := range contentValues {
					currentCheckout.SetInfo(key, value)
				}

				return "ok", nil
			} else {
				return nil, env.ErrorNew("payment method not allowed")
			}
		}
	}

	return nil, env.ErrorNew("payment method not found")
}

// WEB REST API function to set payment method
func restCheckoutSetShippingMethod(params *api.T_APIHandlerParams) (interface{}, error) {
	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// looking for shipping method
	for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
		if shippingMethod.GetCode() == params.RequestURLParams["method"] {
			if shippingMethod.IsAllowed(currentCheckout) {

				// looking for shipping rate
				for _, shippingRate := range shippingMethod.GetRates(currentCheckout) {
					if shippingRate.Code == params.RequestURLParams["rate"] {

						err := currentCheckout.SetShippingMethod(shippingMethod)
						if err != nil {
							return nil, env.ErrorDispatch(err)
						}

						err = currentCheckout.SetShippingRate(shippingRate)
						if err != nil {
							return nil, env.ErrorDispatch(err)
						}

						return "ok", nil
					}
				}

			} else {
				return nil, env.ErrorNew("shipping method not allowed")
			}
		}
	}

	return nil, env.ErrorNew("shipping method and/or rate were not found")
}

// WEB REST API function to submit checkout information and make order
func restSubmit(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return currentCheckout.Submit()
}
