package braintree

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
)

// GetCode returns payment method code for use in business logic
func (it *BraintreePaymentMethod) GetCode() string {
	return ConstPaymentCode
}

// GetInternalName returns the human readable name of the payment method
func (it *BraintreePaymentMethod) GetInternalName() string {
	return ConstPaymentInternalName
}

// GetName returns the user customized name of the payment method
func (it *BraintreePaymentMethod) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathName))
}

// GetType returns type of payment method according to "github.com/ottemo/foundation/app/models/checkout"
func (it *BraintreePaymentMethod) GetType() string {
	return checkout.ConstPaymentTypePostCC // TODO decide this, other or new type
}

// IsAllowed checks for payment method applicability
func (it *BraintreePaymentMethod) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathBraintreeEnabled))
}

// IsTokenable returns possibility to save token for this payment method
func (it *BraintreePaymentMethod) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return (false)
}

// TODO Authorize makes payment method authorize operation
func (it *BraintreePaymentMethod) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	gateway := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "checkout/submit-fake"
	formValues := map[string]string{}
	token := "TODO generate token"

	htmlText := "<form id='braintree-form' method='post' action='" + gateway + "'>"
	for key, value := range formValues {
		htmlText += "<input type='hidden' name='" + key + "' value='" + value + "' />"
	}

	htmlText += "<div id='payment-form'></div>"
	htmlText += "<input type='submit' value='Submit' />"
	htmlText += "</form>"

	//htmlText += "<script src='https://js.braintreegateway.com/js/braintree-2.30.0.min.js'></script>"
	//htmlText += "<script>var script = document.createElement( 'script' );</script>"
	//htmlText += "<script>script.src = 'https://js.braintreegateway.com/js/braintree-2.30.0.min.js';</script>"
	//htmlText += "<script>$('#payment-form').append( script );</script>"
	htmlText += "<script> // " + token + "\n"
	//htmlText += "var brainformSubmitStage = 1; // " + token + "\n"
	//htmlText += "$( '#braintree-form' ).submit(function( event ) {"
	//htmlText += 	"	if (brainformSubmitStage == 1) {" +
	//		"		console.log( 'Handler for .submit() called. Preventing.' );event.preventDefault();" +
	//		"		brainformSubmitStage = 2;" +
	//		"	}"
	//htmlText += "		elseif (brainformSubmitStage == 2) {" +
	//		"		brainformSubmitStage = 3;" +
	//		"		console.log('request for nonse');" +
	//		"	}"
	//htmlText += "		elseif (brainformSubmitStage == 3) {" +
	//		"		console.log('NONSE should be here');" +
	//		"	}"
	//htmlText += "});"
	//htmlText += "$.getScript('https://js.braintreegateway.com/js/braintree-2.30.0.min.js', function(){"
	//htmlText += "	console.log('Token: "+token+"');"
	//htmlText += "	braintree.setup('"+token+"', 'dropin', {container: 'payment-form'});"
	//htmlText += "});"
	htmlText += "</script>"
	//htmlText += "<script>braintree.setup('"+token+"', 'dropin', {container: 'payment-form'});</script>"

	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "FORM: "+htmlText+"\n")

	//env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "NEW TRANSACTION: "+
	//	"Visitor ID - "+utils.InterfaceToString(orderInstance.Get("visitor_id"))+", "+
	//	"Order ID - "+utils.InterfaceToString(orderInstance.GetID()))

	return api.StructRestRedirect{Result: htmlText, Location: utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "checkout/submit"}, nil
}

// Capture makes payment method capture operation
// - at time of implementation this method is not used anywhere
func (it *BraintreePaymentMethod) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e40b3354-5522-4027-b56d-ea9e736f637f", "Not implemented")
}

// Refund will return funds on the given order
// - at time of implementation this method is not used anywhere
func (it *BraintreePaymentMethod) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f28dab69-77a1-438c-a14f-70c700a70c3b", "Not implemented")
}

// Void will mark the order and capture as void
// - at time of implementation this method is not used anywhere
func (it *BraintreePaymentMethod) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "194a4323-4cc4-41a0-ae90-0b578dc4b73a", "Not implemented")
}
