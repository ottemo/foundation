package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/actors/payment/zeropay"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("checkout", api.ConstRESTOperationGet, APIGetCheckout)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout/payment/methods", api.ConstRESTOperationGet, APIGetPaymentMethods)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout/shipping/methods", api.ConstRESTOperationGet, APIGetShippingMethods)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout/shipping/address", api.ConstRESTOperationUpdate, APISetShippingAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout/billing/address", api.ConstRESTOperationUpdate, APISetBillingAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout/payment/method/:method", api.ConstRESTOperationUpdate, APISetPaymentMethod)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout/shipping/method/:method/:rate", api.ConstRESTOperationUpdate, APISetShippingMethod)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("checkout", api.ConstRESTOperationUpdate, APISetCheckoutInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("checkout/submit", api.ConstRESTOperationCreate, APISubmitCheckout)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIGetCheckout returns information related to current checkkout
func APIGetCheckout(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// record rts event for  checkout
	eventData := map[string]interface{}{"session": context.GetSession(), "checkout": currentCheckout}
	env.Event("api.checkout.visit", eventData)

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
		"info":       nil,
	}

	if billingAddress := currentCheckout.GetBillingAddress(); billingAddress != nil {
		result["billing_address"] = billingAddress.ToHashMap()
	}

	if shippingAddress := currentCheckout.GetShippingAddress(); shippingAddress != nil {
		shippingAddressMap := shippingAddress.ToHashMap()

		if notes := utils.InterfaceToString(currentCheckout.GetInfo("notes")); notes != "" {
			shippingAddressMap["notes"] = notes
		}

		result["shipping_address"] = shippingAddressMap
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
	result["grandtotal"] = currentCheckout.GetGrandTotal()

	result["tax_amount"] = currentCheckout.GetTaxAmount()
	result["taxes"] = currentCheckout.GetTaxes()

	result["discount_amount"] = currentCheckout.GetDiscountAmount()
	result["discounts"] = currentCheckout.GetDiscounts()

	// prevent from showing cc values in info
	infoMap := make(map[string]interface{})
	for key, value := range utils.InterfaceToMap(currentCheckout.GetInfo("*")) {
		if key != "cc" {
			infoMap[key] = value
		}
	}

	result["info"] = infoMap

	return result, nil
}

// APIGetPaymentMethods returns currently available payment methods
func APIGetPaymentMethods(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	type ResultValue struct {
		Name string
		Code string
		Type string
	}
	var result []ResultValue

	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.IsAllowed(currentCheckout) {
			result = append(result, ResultValue{Name: paymentMethod.GetName(), Code: paymentMethod.GetCode(), Type: paymentMethod.GetType()})
		}
	}

	return result, nil
}

// APIGetShippingMethods returns currently available shipping methods
func APIGetShippingMethods(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	type ResultValue struct {
		Name  string
		Code  string
		Rates []checkout.StructShippingRate
	}
	var result []ResultValue

	for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
		if shippingMethod.IsAllowed(currentCheckout) {
			result = append(result, ResultValue{Name: shippingMethod.GetName(), Code: shippingMethod.GetCode(), Rates: shippingMethod.GetRates(currentCheckout)})
		}
	}

	return result, nil
}

// APISetCheckoutInfo allows to specify and assign to checkout extra information
func APISetCheckoutInfo(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for key, value := range requestData {
		err := currentCheckout.SetInfo(key, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// updating session
	checkout.SetCurrentCheckout(context, currentCheckout)

	return "ok", nil
}

// checkoutObtainAddress is an internal usage function used create and validate address
//   - address data supposed to be in request content
func checkoutObtainAddress(data interface{}) (visitor.InterfaceVisitorAddress, error) {

	var err error
	var currentVisitorID string
	var addressData map[string]interface{}

	switch context := data.(type) {
	case api.InterfaceApplicationContext:
		currentVisitorID = utils.InterfaceToString(context.GetSession().Get(visitor.ConstSessionKeyVisitorID))
		addressData, err = api.GetRequestContentAsMap(context)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	case map[string]interface{}:
		addressData = context
	default:
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "995c0d71-5465-4003-9ed2-347bcf6255c5", "unknown address data type")
	}

	// checking for address id was specified, if it was - making sure it correct
	if addressID, present := addressData["id"]; present {

		// loading specified address by id
		visitorAddress, err := visitor.LoadVisitorAddressByID(utils.InterfaceToString(addressID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// checking address owner is current visitor

		if currentVisitorID != "" && visitorAddress.GetVisitorID() != currentVisitorID {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "bef27714-4ac5-4705-b59a-47c8e0bc5aa4", "address id is not related to current visitor")
		}

		return visitorAddress, nil
	}

	// creating new address model instance
	visitorAddressModel, err := visitor.GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filling new instance with request provided data
	for attribute, value := range addressData {
		err := visitorAddressModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// setting address owner to current visitor (for sure)
	if currentVisitorID != "" {
		visitorAddressModel.Set("visitor_id", currentVisitorID)
	}

	// if address id was specified it means that address was changed, so saving it
	// new address we are not saving as if could be temporary address
	if visitorAddressModel.GetID() != "" {
		err = visitorAddressModel.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	return visitorAddressModel, nil
}

// APISetShippingAddress specifies shipping address for a current checkout
func APISetShippingAddress(context api.InterfaceApplicationContext) (interface{}, error) {
	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	address, err := checkoutObtainAddress(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = currentCheckout.SetShippingAddress(address)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestContents, _ := api.GetRequestContentAsMap(context)

	if notes, present := requestContents["notes"]; present {
		currentCheckout.SetInfo("notes", notes)
	}

	// updating session
	checkout.SetCurrentCheckout(context, currentCheckout)

	return address.ToHashMap(), nil
}

// APISetBillingAddress specifies billing address for a current checkout
func APISetBillingAddress(context api.InterfaceApplicationContext) (interface{}, error) {
	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	address, err := checkoutObtainAddress(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = currentCheckout.SetBillingAddress(address)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// updating session
	checkout.SetCurrentCheckout(context, currentCheckout)

	return address.ToHashMap(), nil
}

// APISetPaymentMethod assigns payment method to current checkout
//   - "method" argument specifies requested payment method (it should be available for a meaning time)
func APISetPaymentMethod(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// looking for payment method
	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.GetCode() == context.GetRequestArgument("method") {
			if paymentMethod.IsAllowed(currentCheckout) {

				// updating checkout payment method
				err := currentCheckout.SetPaymentMethod(paymentMethod)
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}

				// checking for additional info
				contentValues, _ := api.GetRequestContentAsMap(context)
				for key, value := range contentValues {
					currentCheckout.SetInfo(key, value)
				}

				// visitor event for setting payment method
				eventData := map[string]interface{}{"session": context.GetSession(), "paymentMethod": paymentMethod, "checkout": currentCheckout}
				env.Event("api.checkout.setPayment", eventData)

				// updating session
				checkout.SetCurrentCheckout(context, currentCheckout)

				return "ok", nil
			}
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "bd07849e-8789-4316-924c-9c754efbc348", "payment method not allowed")
		}
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b8384a47-8806-4a54-90fc-cccb5e958b4e", "payment method not found")
}

// APISetShippingMethod assigns shipping method and shipping rate to current checkout
//   - "method" argument specifies requested shipping method (it should be available for a meaning time)
//   - "rate" argument specifies requested shipping rate (it should be available and belongs to shipping method)
func APISetShippingMethod(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// looking for shipping method
	for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
		if shippingMethod.GetCode() == context.GetRequestArgument("method") {
			if shippingMethod.IsAllowed(currentCheckout) {

				// looking for shipping rate
				for _, shippingRate := range shippingMethod.GetRates(currentCheckout) {
					if shippingRate.Code == context.GetRequestArgument("rate") {

						err := currentCheckout.SetShippingMethod(shippingMethod)
						if err != nil {
							return nil, env.ErrorDispatch(err)
						}

						err = currentCheckout.SetShippingRate(shippingRate)
						if err != nil {
							return nil, env.ErrorDispatch(err)
						}

						// updating session
						checkout.SetCurrentCheckout(context, currentCheckout)

						return "ok", nil
					}
				}

			} else {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d7fb6ff2-b914-467b-bf56-b8d2bea472ef", "shipping method not allowed")
			}
		}
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "279a645c-6a03-44de-95c0-2651a51440fa", "shipping method and/or rate were not found")
}

// APISubmitCheckout submits current checkout and creates a new order base on it
func APISubmitCheckout(context api.InterfaceApplicationContext) (interface{}, error) {

	// preparations
	//--------------
	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentVisitorID := utils.InterfaceToString(context.GetSession().Get(visitor.ConstSessionKeyVisitorID))

	addressInfoToAddress := func(addressInfo interface{}) (visitor.InterfaceVisitorAddress, error) {
		var addressData map[string]interface{}

		switch typedValue := addressInfo.(type) {
		case map[string]interface{}:
			typedValue["visitor_id"] = currentVisitorID
			addressData = typedValue
		case string:
			addressData = map[string]interface{}{"visitor_id": currentVisitorID, "id": typedValue}
		}
		return checkoutObtainAddress(addressData)
	}

	// checking for specified shipping address
	//-----------------------------------------
	if shippingAddressInfo := utils.GetFirstMapValue(requestData, "shipping_address", "shippingAddress", "ShippingAddress"); shippingAddressInfo != nil {
		address, err := addressInfoToAddress(shippingAddressInfo)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		currentCheckout.SetShippingAddress(address)
	}

	// checking for specified billing address
	//----------------------------------------
	if billingAddressInfo := utils.GetFirstMapValue(requestData, "billing_address", "billingAddress", "BillingAddress"); billingAddressInfo != nil {
		address, err := addressInfoToAddress(billingAddressInfo)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		currentCheckout.SetBillingAddress(address)
	}

	// checking for specified payment method
	//---------------------------------------
	if specifiedPaymentMethod := utils.GetFirstMapValue(requestData, "payment_method", "paymentMethod"); specifiedPaymentMethod != nil {
		var found bool
		for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
			if paymentMethod.GetCode() == specifiedPaymentMethod {
				if paymentMethod.IsAllowed(currentCheckout) {
					err := currentCheckout.SetPaymentMethod(paymentMethod)
					if err != nil {
						return nil, env.ErrorDispatch(err)
					}
					found = true
					break
				}
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "bd07849e-8789-4316-924c-9c754efbc348", "payment method not allowed")
			}
		}

		if !found {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b8384a47-8806-4a54-90fc-cccb5e958b4e", "payment method not found")
		}
	}

	// set ZeroPayment method for checkout without payment method
	if currentCheckout.GetPaymentMethod() == nil {
		for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
			if zeropay.ConstPaymentZeroPaymentCode == paymentMethod.GetCode() {
				if paymentMethod.IsAllowed(currentCheckout) {
					err := currentCheckout.SetPaymentMethod(paymentMethod)
					if err != nil {
						return nil, env.ErrorDispatch(err)
					}
				}
			}
		}
	}

	// checking for specified shipping method
	//----------------------------------------
	specifiedShippingMethod := utils.GetFirstMapValue(requestData, "shipping_method", "shipppingMethod")
	specifiedShippingMethodRate := utils.GetFirstMapValue(requestData, "shipppingRate", "shipping_rate")

	if specifiedShippingMethod != nil && specifiedShippingMethodRate != nil {
		var methodFound, rateFound bool

		for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
			if shippingMethod.GetCode() == context.GetRequestArgument("method") {
				if shippingMethod.IsAllowed(currentCheckout) {
					methodFound = true

					for _, shippingRate := range shippingMethod.GetRates(currentCheckout) {
						if shippingRate.Code == context.GetRequestArgument("rate") {
							err = currentCheckout.SetShippingMethod(shippingMethod)
							if err != nil {
								return nil, env.ErrorDispatch(err)
							}
							currentCheckout.SetShippingRate(shippingRate)
							if err != nil {
								return nil, env.ErrorDispatch(err)
							}

							rateFound = true
							break
						}
					}
					break
				}
			}
		}

		if !methodFound || !rateFound {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "279a645c-6a03-44de-95c0-2651a51440fa", "shipping method and/or rate were not found")
		}
	}

	return currentCheckout.Submit()
}
