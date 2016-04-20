package authorizenet

import (
	"fmt"
	"time"

	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"math/rand"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
)

// GetInternalName returns the name of the payment method
func (it DirectPostMethod) GetInternalName() string {
	return ConstPaymentNameDPM
}

// GetName returns the user customized name of the payment method
func (it *DirectPostMethod) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDPMTitle))
}

// GetCode returns payment method code
func (it *DirectPostMethod) GetCode() string {
	return ConstPaymentCodeDPM
}

// IsTokenable returns possibility to save token for this payment method
func (it *DirectPostMethod) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return false
}

// GetType returns type of payment method
func (it *DirectPostMethod) GetType() string {
	return checkout.ConstPaymentTypePostCC
}

// IsAllowed checks for method applicability
func (it *DirectPostMethod) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathDPMEnabled))
}

// Authorize makes payment method authorize operation
func (it *DirectPostMethod) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	// crypting fingerprint
	//---------------------
	loginID := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDPMLogin))
	sequence := fmt.Sprintf("%d", rand.Intn(999)+1)
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	amount := fmt.Sprintf("%.2f", orderInstance.GetGrandTotal())
	transactionKey := []byte(utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDPMKey)))

	hmacEncoder := hmac.New(md5.New, transactionKey)
	hmacEncoder.Write([]byte(loginID + "^" + sequence + "^" + timeStamp + "^" + amount + "^USD"))
	fingerprint := hex.EncodeToString(hmacEncoder.Sum(nil))

	billingAddress := orderInstance.GetBillingAddress()
	shippingAddress := orderInstance.GetShippingAddress()

	// preparing post form values
	//---------------------------
	formValues := map[string]string{
		"x_relay_response": utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDPMCheckout)),
		"x_relay_url":      app.GetFoundationURL("") + "authorizenet/relay",

		"x_test_request": utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDPMTest)),

		"x_fp_sequence":  sequence,
		"x_fp_timestamp": timeStamp,
		"x_fp_hash":      string(fingerprint),

		"x_amount": amount,

		"x_login":         loginID,
		"x_type":          utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDPMAction)),
		"x_method":        "CC",
		"x_currency_code": "USD",

		"x_first_name": billingAddress.GetFirstName(),
		"x_last_name":  billingAddress.GetLastName(),
		"x_company":    billingAddress.GetCompany(),
		"x_address":    billingAddress.GetAddress(),
		"x_city":       billingAddress.GetCity(),
		"x_state":      billingAddress.GetState(),
		"x_zip":        billingAddress.GetZipCode(),
		"x_country":    billingAddress.GetCountry(),
		"x_phone":      billingAddress.GetPhone(),
		"x_cust_id":    billingAddress.GetVisitorID(),
		"x_email":      utils.InterfaceToString(orderInstance.Get("customer_email")),

		"x_ship_to_first_name": shippingAddress.GetFirstName(),
		"x_ship_to_last_name":  shippingAddress.GetLastName(),
		"x_ship_to_company":    shippingAddress.GetCompany(),
		"x_ship_to_address":    shippingAddress.GetAddress(),
		"x_ship_to_city":       shippingAddress.GetCity(),
		"x_ship_to_state":      shippingAddress.GetState(),
		"x_ship_to_zip":        shippingAddress.GetZipCode(),
		"x_ship_to_country":    shippingAddress.GetCountry(),

		"x_exp_date":         "$CC_MONTH/$CC_YEAR",
		"x_card_num":         "$CC_NUM",
		"x_duplicate_window": "30",
		"x_session":          utils.InterfaceToString(paymentInfo["sessionID"]),
	}

	// generating post form
	//---------------------
	gateway := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDPMGateway))

	htmlText := "<form method='post' action='" + gateway + "'>"
	for key, value := range formValues {
		htmlText += "<input type='hidden' name='" + key + "' value='" + value + "' />"
	}
	htmlText += "<input type='submit' value='Submit' />"
	htmlText += "</form>"

	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "FORM: "+htmlText)

	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "NEW TRANSACTION: "+
		"Visitor ID - "+utils.InterfaceToString(orderInstance.Get("visitor_id"))+", "+
		"Order ID - "+utils.InterfaceToString(orderInstance.GetID()))

	return api.StructRestRedirect{Result: htmlText, Location: utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathFoundationURL)) + "authorizenet/relay"}, nil
}

// Capture makes payment method capture operation
func (it *DirectPostMethod) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ed753163-d708-4884-aae8-3aa1dc9bf9f4", "Not implemented")
}

// Refund will return funds on the given order :: Not Implemented Yet
func (it *DirectPostMethod) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2dc9b523-4e53-4ff5-9d49-1dfdadf3fb44", "Not implemented")
}

// Void will mark the order and capture as void
func (it *DirectPostMethod) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d682f87e-4d51-473b-a00d-191d28e807f5", "Not implemented")
}
