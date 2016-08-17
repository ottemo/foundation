package paypal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
)

// GetInternalName returns the name of the payment method
func (it Express) GetInternalName() string {
	return ConstPaymentName
}

// GetName returns the user customized name of the payment method
func (it *Express) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTitle))
}

// GetCode returns payment method code
func (it *Express) GetCode() string {
	return ConstPaymentCode
}

// GetType returns the type of payment method
func (it *Express) GetType() string {
	return checkout.ConstPaymentTypeRemote
}

// IsTokenable checks for method applicability
func (it *Express) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return false
}

// IsAllowed checks for method applicability
func (it *Express) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEnabled))
}

// Authorize makes payment method authorize operation
func (it *Express) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	// getting order information
	//--------------------------
	grandTotal := orderInstance.GetGrandTotal()
	shippingPrice := orderInstance.GetShippingAmount()

	// getting request param values
	//-----------------------------
	user := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathUser))
	password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPass))
	signature := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathSignature))
	action := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAction))

	amount := fmt.Sprintf("%.2f", grandTotal)
	shippingAmount := fmt.Sprintf("%.2f", shippingPrice)
	itemAmount := fmt.Sprintf("%.2f", grandTotal-shippingPrice)

	description := "Purchase%20for%20%24" + fmt.Sprintf("%.2f", grandTotal)
	custom := orderInstance.GetID()

	cancelURL := app.GetFoundationURL("paypal/cancel")
	returnURL := app.GetFoundationURL("paypal/success")

	// making NVP request
	//-------------------
	requestParams := "USER=" + user +
		"&PWD=" + password +
		"&SIGNATURE=" + signature +
		"&METHOD=SetExpressCheckout" +
		"&VERSION=78" +
		"&PAYMENTREQUEST_0_PAYMENTACTION=" + action +
		"&PAYMENTREQUEST_0_AMT=" + amount +
		"&PAYMENTREQUEST_0_SHIPPINGAMT=" + shippingAmount +
		"&PAYMENTREQUEST_0_ITEMAMT=" + itemAmount +
		"&PAYMENTREQUEST_0_DESC=" + description +
		"&PAYMENTREQUEST_0_CUSTOM=" + custom +
		"&PAYMENTREQUEST_0_CURRENCYCODE=USD" +
		"&cancelURL=" + cancelURL +
		"&returnURL=" + returnURL

	//	println(requestParams)

	nvpGateway := paymentPayPalExpress[
		ConstPaymentPayPalNvp][
		utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalExpressGateway))]

		request, err := http.NewRequest("GET", nvpGateway+"?"+requestParams, nil)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// reading/decoding response from NVP
	//-----------------------------------
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	responseValues, err := url.ParseQuery(string(responseData))
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7b6a22eb-a7b1-40b5-a47a-5528753135b2", "payment unexpected response")
	}

	if responseValues.Get("ACK") != "Success" || responseValues.Get("TOKEN") == "" {
		if responseValues.Get("L_ERRORCODE0") != "" {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5ec8dc9b-f72c-4f35-9b1e-f4353753dd9e", "payment error "+responseValues.Get("L_ERRORCODE0")+": "+"L_LONGMESSAGE0")
		}
	}
	waitingTokensMutex.Lock()
	waitingTokens[responseValues.Get("TOKEN")] = utils.InterfaceToString(paymentInfo["sessionID"])
	waitingTokensMutex.Unlock()

	env.Log("paypal.log", env.ConstLogPrefixInfo, "NEW TRANSACTION: "+
		"Visitor ID - "+utils.InterfaceToString(orderInstance.Get("visitor_id"))+", "+
		"Order ID - "+utils.InterfaceToString(orderInstance.GetID())+", "+
		"TOKEN - "+utils.InterfaceToString(responseValues.Get("TOKEN")))

	// redirecting user to PayPal server for following checkout
	//---------------------------------------------------------
	payPalGateway := paymentPayPalExpress[
		ConstPaymentPayPalGateway][
		utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalExpressGateway))]
	redirectGateway := payPalGateway + "&token=" + responseValues.Get("TOKEN")
	return api.StructRestRedirect{
		Result:   "redirect",
		Location: redirectGateway,
	}, nil
}

// Capture payment method capture operation
func (it *Express) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8569a80f-ffc0-4c99-ab79-89d0cbb90ca7", "Not implemented")
}

// Refund makes payment method refund operation
func (it *Express) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bcd6d1c3-12ad-4e62-a2c6-e90bc61badd3", "Not implemented")
}

// Void makes payment method void operation
func (it *Express) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "78083247-0f7b-4cf1-84a2-09f6162ca350", "Not implemented")
}
