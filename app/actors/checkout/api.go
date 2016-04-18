package checkout

import (
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/actors/payment/paypal"
	"github.com/ottemo/foundation/app/actors/payment/zeropay"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("checkout", APIGetCheckout)

	// Addresses
	service.PUT("checkout/shipping/address", APISetShippingAddress)
	service.PUT("checkout/billing/address", APISetBillingAddress)

	// Shipping method
	service.GET("checkout/shipping/methods", APIGetShippingMethods)
	service.PUT("checkout/shipping/method/:method/:rate", APISetShippingMethod)

	// Payment method
	service.GET("checkout/payment/methods", APIGetPaymentMethods)
	service.PUT("checkout/payment/method/:method", APISetPaymentMethod)

	// Finalize
	service.PUT("checkout", APISetCheckoutInfo)
	service.POST("checkout/submit", APISubmitCheckout)

	// service.PUT("checkout/paymentdetails", APISetPaymentDetails)
	return nil
}

// APIGetCheckout returns information related to current checkkout
func APIGetCheckout(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(context, false)
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
	}

	result["grandtotal"] = currentCheckout.GetGrandTotal()
	result["subtotal"] = currentCheckout.GetSubtotal()

	result["shipping_amount"] = currentCheckout.GetShippingAmount()

	result["tax_amount"] = currentCheckout.GetTaxAmount()
	result["taxes"] = currentCheckout.GetTaxes()

	result["discount_amount"] = currentCheckout.GetDiscountAmount()
	result["discounts"] = currentCheckout.GetDiscounts()

	// The info map is only returned for logged out users
	infoMap := make(map[string]interface{})

	for key, value := range utils.InterfaceToMap(currentCheckout.GetInfo("*")) {
		// prevent from showing cc values in info
		if key != "cc" {
			infoMap[key] = value
		}
	}

	result["info"] = infoMap

	return result, nil
}

// APIGetPaymentMethods returns currently available payment methods
func APIGetPaymentMethods(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(context, false)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	type ResultValue struct {
		Name      string
		Code      string
		Type      string
		Tokenable bool
	}
	var result []ResultValue

	// for checkout that contain subscription items we will show only payment methods that allows to save token
	isSubscription := currentCheckout.IsSubscription()

	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.IsAllowed(currentCheckout) && (!isSubscription || paymentMethod.IsTokenable(currentCheckout)) {
			result = append(result, ResultValue{Name: paymentMethod.GetName(), Code: paymentMethod.GetCode(), Type: paymentMethod.GetType(), Tokenable: paymentMethod.IsTokenable(currentCheckout)})
		}
	}

	return result, nil
}

// APIGetShippingMethods returns currently available shipping methods
func APIGetShippingMethods(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCheckout, err := checkout.GetCurrentCheckout(context, false)
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

	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
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

	visitorAddressModel, err := checkout.ValidateAddress(addressData)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// setting address owner to current visitor (for sure)
	if currentVisitorID != "" {
		visitorAddressModel.Set("visitor_id", currentVisitorID)
	}

	// if address id was specified it means that address was changed, so saving it
	// new address we are not saving as if could be temporary address
	if (visitorAddressModel.GetID() != "" || currentVisitorID != "") && utils.InterfaceToBool(addressData["save"]) {
		err = visitorAddressModel.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	return visitorAddressModel, nil
}

// APISetShippingAddress specifies shipping address for a current checkout
func APISetShippingAddress(context api.InterfaceApplicationContext) (interface{}, error) {
	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
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
	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
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

	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
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

	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
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

// checkoutObtainToken is an internal usage function used to create or load credit card for visitor
func checkoutObtainToken(currentCheckout checkout.InterfaceCheckout, creditCardInfo map[string]interface{}) (visitor.InterfaceVisitorCard, error) {

	currentVisitor := currentCheckout.GetVisitor()
	currentVisitorID := ""
	if currentVisitor != nil {
		currentVisitorID = currentCheckout.GetVisitor().GetID()
	}

	// checking for address id was specified, if it was - making sure it correct
	if creditCardID := utils.GetFirstMapValue(creditCardInfo, "id", "_id"); currentVisitorID != "" && creditCardID != nil {

		// loading specified credit card by id
		visitorCard, err := visitor.LoadVisitorCardByID(utils.InterfaceToString(creditCardID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// checking address owner is current visitor
		if visitorCard.GetVisitorID() != currentVisitorID {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3b5446ef-dd70-4bdc-8d30-817cb7b48d05", "credit card id is not related to current visitor")
		}

		return visitorCard, nil
	}

	paymentMethod := currentCheckout.GetPaymentMethod()
	if !paymentMethod.IsTokenable(currentCheckout) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5b05cc24-2184-47cc-b2dc-77cb41035698", "for selected payment method credit card can't be saved")
	}

	// put required key to create token from payment method using only zero amount authorize
	paymentInfo := map[string]interface{}{
		checkout.ConstPaymentActionTypeKey: checkout.ConstPaymentActionTypeCreateToken,
		"cc": creditCardInfo,
	}

	// contains creditCardLastFour, creditCardType, responseMessage, responseResult, transactionID, creditCardExp
	paymentResult, err := paymentMethod.Authorize(nil, paymentInfo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	authorizeCardResult := utils.InterfaceToMap(paymentResult)
	if !utils.KeysInMapAndNotBlank(authorizeCardResult, "transactionID", "creditCardLastFour") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "22e17290-56f3-452a-8d54-18d5a9eb2833", "transaction can't be obtained")
	}

	// create visitor address operation
	//---------------------------------
	visitorCardModel, err := visitor.GetVisitorCardModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// override credit card info with provided from payment info
	// TODO: payment should have interface method that return predefined struct for 0 authorize
	creditCardInfo["token_id"] = authorizeCardResult["transactionID"]
	creditCardInfo["payment"] = paymentMethod.GetCode()
	creditCardInfo["type"] = authorizeCardResult["creditCardType"]
	creditCardInfo["number"] = authorizeCardResult["creditCardLastFour"]
	creditCardInfo["expiration_date"] = authorizeCardResult["creditCardExp"]
	creditCardInfo["token_updated"] = time.Now()
	creditCardInfo["created_at"] = time.Now()

	// filling new instance with request provided data
	// TODO: check other places with such code:
	// it's possible to put here id's or visitorID's in different way as they used in Set method
	// and override value that should be
	for attribute, value := range creditCardInfo {
		err := visitorCardModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// setting credit card owner to current visitor (for sure)
	if currentVisitorID != "" {
		visitorCardModel.Set("visitor_id", currentVisitorID)
	}

	// save cc token if using appropriate payment adapter
	if (visitorCardModel.GetID() != "" || currentVisitorID != "") &&
		paymentMethod.GetCode() == paypal.ConstPaymentPayPalPayflowCode {

		err = visitorCardModel.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	return visitorCardModel, nil
}

// APISetPaymentDetails specifies payment details for a current checkout
func APISetPaymentDetails(context api.InterfaceApplicationContext) (interface{}, error) {
	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestContents, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	creditCard, err := checkoutObtainToken(currentCheckout, requestContents)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentCheckout.SetInfo("cc", creditCard)

	// updating session
	checkout.SetCurrentCheckout(context, currentCheckout)

	var result map[string]interface{}

	// hide token ID
	for key, value := range creditCard.ToHashMap() {
		if key != "token" {
			result[key] = value
		}
	}

	return result, nil
}

// APISubmitCheckout submits current checkout and creates a new order base on it
func APISubmitCheckout(context api.InterfaceApplicationContext) (interface{}, error) {

	// preparations
	//--------------
	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Handle custom information set in case of one request submit
	if customInfo := utils.GetFirstMapValue(requestData, "custom_info"); customInfo != nil {
		for key, value := range utils.InterfaceToMap(customInfo) {
			currentCheckout.SetInfo(key, value)
		}
	}

	currentCheckout.SetInfo("session_id", context.GetSession().GetID())
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

	currentPaymentMethod := currentCheckout.GetPaymentMethod()
	// set ZeroPayment method for checkout without payment method
	if currentPaymentMethod == nil || (currentCheckout.GetGrandTotal() == 0 && currentPaymentMethod.GetCode() != zeropay.ConstPaymentZeroPaymentCode) {
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

	// Now that checkout is about to submit we want to see if we can turn our cc info into a token
	// if this errors out, it just means that the criteria wasn't met to create a token. Which is ok
	specifiedCreditCard := currentCheckout.GetInfo("cc")
	creditCard, err := checkoutObtainToken(currentCheckout, utils.InterfaceToMap(specifiedCreditCard))
	if err == nil {
		currentCheckout.SetInfo("cc", creditCard)
	}

	return currentCheckout.Submit()
}
