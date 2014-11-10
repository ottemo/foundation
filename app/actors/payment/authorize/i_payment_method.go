package authorize

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

// GetName returns the payment method name
func (it *AuthorizeNetDPM) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_DPM_TITLE))
}

// GetCode returns the payment method code
func (it *AuthorizeNetDPM) GetCode() string {
	return PAYMENT_CODE_DPM
}

// GetType returns the type of payment method
func (it *AuthorizeNetDPM) GetType() string {
	return checkout.PAYMENT_TYPE_POST_CC
}

// IsAllowed checks for method applicability
func (it *AuthorizeNetDPM) IsAllowed(checkoutInstance checkout.I_Checkout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(CONFIG_PATH_DPM_ENABLED))
}

// Authorize executes the payment method authorization
func (it *AuthorizeNetDPM) Authorize(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {

	// crypting fingerprint
	//---------------------
	loginID := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_DPM_LOGIN))
	sequence := fmt.Sprintf("%d", rand.Intn(999)+1)
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	amount := fmt.Sprintf("%0.0f", orderInstance.GetGrandTotal())
	transactionKey := []byte(utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_DPM_KEY)))

	hmacEncoder := hmac.New(md5.New, transactionKey)
	hmacEncoder.Write([]byte(loginID + "^" + sequence + "^" + timeStamp + "^" + amount + "^USD"))
	fingerprint := hex.EncodeToString(hmacEncoder.Sum(nil))

	billingAddress := orderInstance.GetBillingAddress()
	shippingAddress := orderInstance.GetShippingAddress()

	// preparing post form values
	//---------------------------
	formValues := map[string]string{
		"x_relay_response": utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_DPM_CHECKOUT)),
		"x_relay_url":      utils.InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_FOUNDATION_URL)) + "authorizenet/relay",

		"x_test_request": utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_DPM_TEST)),

		"x_fp_sequence":  sequence,
		"x_fp_timestamp": timeStamp,
		"x_fp_hash":      string(fingerprint),

		"x_amount": amount,

		"x_login":         loginID,
		"x_type":          utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_DPM_ACTION)),
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
		"x_session":          utils.InterfaceToString(paymentInfo["sessionId"]),
	}

	// generating post form
	//---------------------
	gateway := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_DPM_GATEWAY))

	htmlText := "<form method='post' action='" + gateway + "'>"
	for key, value := range formValues {
		htmlText += "<input type='hidden' name='" + key + "' value='" + value + "' />"
	}
	htmlText += "<input type='submit' value='Submit' />"
	htmlText += "</form>"

	env.Log("authorizenet.log", env.LOG_PREFIX_INFO, "NEW TRANSACTION: "+
		"Visitor ID - "+utils.InterfaceToString(orderInstance.Get("visitorID"))+", "+
		"Order ID - "+utils.InterfaceToString(orderInstance.GetId()))

	return api.T_RestRedirect{Result: htmlText, Location: utils.InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_FOUNDATION_URL)) + "authorizenet/relay"}, nil
}

// Capture will secure the funds once an order has been fulfilled :: Not Implemented Yet
func (it *AuthorizeNetDPM) Capture(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew("Not implemented")
}

// Refund will return funds on the given order :: Not Implemented Yet
func (it *AuthorizeNetDPM) Refund(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew("Not implemented")
}

// Void will mark the order and capture as void
func (it *AuthorizeNetDPM) Void(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew("Not implemented")
}
