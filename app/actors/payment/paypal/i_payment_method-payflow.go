package paypal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

// GetName returns config value "Title" of payment method
func (it *PayFlowAPI) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowTitle))
}

// GetCode returns PayPal code value for this payment method
func (it *PayFlowAPI) GetCode() string {
	return ConstPaymentPayPalPayflowCode
}

// GetType returns the type of payment method
func (it *PayFlowAPI) GetType() string {
	return checkout.ConstPaymentTypeCreditCard
}

// IsAllowed checks for method applicability
func (it *PayFlowAPI) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathPayPalPayflowEnabled))
}

// Authorize makes payment method authorize operation (currently it's a Authorize zero amount + Sale operations)
func (it *PayFlowAPI) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	authorizeZeroResult, err := it.AuthorizeZeroAmount(orderInstance, paymentInfo)

	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	authorizeZeroResultMap := utils.InterfaceToMap(authorizeZeroResult)
	if value, present := authorizeZeroResultMap["transactionID"]; !present || utils.InterfaceToString(value) == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5e68f079-e8ce-4677-8fb9-89c6f7acbd7f", "Error: token was not created")
	}

	transactionID := utils.InterfaceToString(authorizeZeroResultMap["transactionID"])

	// getting order information
	//--------------------------
	grandTotal := orderInstance.GetGrandTotal()
	amount := fmt.Sprintf("%.2f", grandTotal)

	user := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowUser))
	password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowPass))
	vendor := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowVendor))

	// PayFlow Request Fields
	requestParams := "USER=" + user +
		"&PWD=" + password +
		"&VENDOR=" + vendor +
		"&PARTNER=PayPal" +
		"&VERSION=122" +
		"&TRXTYPE=S" + // Sale operation

		// Credit Card Details Fields
		"&TENDER=C" +
		"&ORIGID=" + utils.InterfaceToString(transactionID) +

		// Payment Details Fields
		"&AMT=" + amount +
		"&CURRENCY=USD" +
		"&VERBOSITY=HIGH" +
		"&INVNUM=" + orderInstance.GetID()

	// adding of access token info to request
	accessTokenInfo, err := it.GetAccessToken(requestParams)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	requestParams = requestParams + "&" + accessTokenInfo

	nvpGateway := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowURL))
	request, err := http.NewRequest("POST", nvpGateway, bytes.NewBufferString(requestParams))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	request.Header.Add("Content-Type", "text/name value")
	request.Header.Add("Host", utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowHost)))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	defer response.Body.Close()

	// reading/decoding response from PayPal
	//-----------------------------------
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	responseValues, err := url.ParseQuery(string(responseBody))
	if err != nil {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b18cdcad-8c21-4acf-a2e0-56e0541103de", "payment unexpected response")
	}

	// get info about transaction from response response
	orderTransactionID := utils.InterfaceToString(responseValues.Get("PNREF"))

	if responseValues.Get("RESPMSG") != "Approved" || orderTransactionID == "" {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e48403bb-c15d-4302-8894-da7146b93260", "payment error: "+responseValues.Get("RESPMSG")+", "+responseValues.Get("PREFPSMSG"))
	}

	env.Log("paypal.log", env.ConstLogPrefixInfo, "NEW TRANSACTION: "+
		"Visitor ID - "+utils.InterfaceToString(orderInstance.Get("visitor_id"))+", "+
		"Order ID - "+utils.InterfaceToString(orderInstance.GetID())+", "+
		"TRANSACTIONID - "+orderTransactionID)

	orderPaymentInfo := utils.InterfaceToMap(orderInstance.Get("payment_info"))

	if oldTransaction, present := orderPaymentInfo["transactionID"]; !present {
		orderPaymentInfo = map[string]interface{}{
			"transactionID":     orderTransactionID,
			"creditCardNumbers": responseValues.Get("ACCT"),
			"creditCardType":    getCreditCardName(utils.InterfaceToString(responseValues.Get("CARDTYPE"))),
		}
	} else {
		orderPaymentInfo["previosTransactionID"] = oldTransaction
		orderPaymentInfo["transactionID"] = orderTransactionID
	}

	orderInstance.Set("payment_info", orderPaymentInfo)
	orderInstance.SetStatus(order.ConstOrderStatusPending)
	orderInstance.Save()
	return nil, nil
}

// Capture payment method capture operation
func (it *PayFlowAPI) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2dc38587-de12-4bdf-9468-a4cef846afe5", "Not implemented")
}

// Refund makes payment method refund operation
func (it *PayFlowAPI) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2dc38587-de12-4bdf-9468-a4cef846afe5", "Not implemented")
}

// Void makes payment method void operation
func (it *PayFlowAPI) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2dc38587-de12-4bdf-9468-a4cef846afe5", "Not implemented")
}

// GetAccessToken returns application access token
func (it *PayFlowAPI) GetAccessToken(originRequestParams string) (string, error) {

	// TRXTYPE=S&AMT=0&USER=xom94ok&PWD=KOL78963&PARTNER=PayPal&SILENTTRAN=TRUE&SECURETOKENID=123448308de1413abc3d60c86cb1f4c5&CREATESECURETOKEN=Y

	secureTokenID := utils.InterfaceToString(time.Now().UnixNano())
	// making NVP request
	//-------------------.
	// PayFlow Request Fields
	requestParams := originRequestParams +
		"&CREATESECURETOKEN=Y" +
		"&SILENTTRAN=TRUE" +
		"&SECURETOKENID=" + secureTokenID

	nvpGateway := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowURL))

	request, err := http.NewRequest("POST", nvpGateway, bytes.NewBufferString(requestParams))
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	request.Header.Add("Content-Type", "text/name value")
	request.Header.Add("Host", utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowHost)))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	defer response.Body.Close()

	// reading/decoding response from PayPal
	//-----------------------------------
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	responseValues, err := url.ParseQuery(string(responseBody))
	if err != nil {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b18cdcad-8c21-4acf-a2e0-56e0541103de", "payment unexpected response")
	}

	if responseValues.Get("RESPMSG") != "Approved" || responseValues.Get("SECURETOKEN") == "" {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e48403bb-c15d-4302-8894-da7146b93260", "payment error: "+responseValues.Get("RESPMSG"))
	}

	token := responseValues.Get("SECURETOKEN")
	if responseValues.Get("SECURETOKENID") != secureTokenID {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9b095f62-b371-4eaf-965f-98eb24206e53", "unexpected response, SECURETOKENID value changed")
	}

	return "SECURETOKEN=" + utils.InterfaceToString(token) + "&SECURETOKENID=" + utils.InterfaceToString(secureTokenID), nil
}

// AuthorizeZeroAmount will do Account Verification and return transaction ID for refer transaction if all info is valid
func (it *PayFlowAPI) AuthorizeZeroAmount(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	if ccInfo, present := paymentInfo["cc"]; !present || !utils.StrKeysInMap(utils.InterfaceToMap(ccInfo), "number", "expire_month", "expire_year") {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "39a27c94-7d39-453d-b464-fd24f7beebcc", "credit card info was not specified")
	}

	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])

	ccExpirationDate := utils.InterfaceToString(ccInfo["expire_year"])
	ccExpirationDate = utils.InterfaceToString(ccInfo["expire_month"]) + ccExpirationDate[len(ccExpirationDate)-2:]
	if len(utils.InterfaceToString(ccInfo["expire_month"])) == 1 {
		ccExpirationDate = "0" + ccExpirationDate
	}

	// getting order information
	//--------------------------
	billingAddress := orderInstance.GetBillingAddress()
	if billingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "17a1b377-7915-4cd0-a4e8-40379e34851f", "no billing address information")
	}

	if orderInstance == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ae8dfe99-4895-43fa-b2af-a575d94fd80a", "no created order")
	}

	user := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowUser))
	password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowPass))
	vendor := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowVendor))

	// PayFlow Request Fields
	requestParams := "USER=" + user +
		"&PWD=" + password +
		"&VENDOR=" + vendor +
		"&PARTNER=PayPal" +
		"&VERSION=122" +
		"&TRXTYPE=A" + // Authorize

		// Credit Card Details Fields
		"&TENDER=C" +
		"&ACCT=" + utils.InterfaceToString(ccInfo["number"]) +
		"&EXPDATE=" + ccExpirationDate +

		// Payer Information Fields
		"&EMAIL=" + utils.InterfaceToString(orderInstance.Get("customer_email")) +
		"&BILLTOFIRSTNAME=" + billingAddress.GetFirstName() +
		"&BILLTOLASTNAME=" + billingAddress.GetLastName() +
		"&BILLTOZIP=" + billingAddress.GetZipCode() +

		// Payment Details Fields
		"&AMT=0" +
		"&VERBOSITY=HIGH" +
		"&INVNUM=" + orderInstance.GetID()

	// add additional params to request
	if ccSecureCode, ccSecureCodePresent := ccInfo["cvv"]; ccSecureCodePresent {
		requestParams = requestParams + "&CVV2=" + utils.InterfaceToString(ccSecureCode)
	}

	nvpGateway := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowURL))
	request, err := http.NewRequest("POST", nvpGateway, bytes.NewBufferString(requestParams))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	request.Header.Add("Content-Type", "text/name value")
	request.Header.Add("Host", utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowHost)))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	defer response.Body.Close()

	// reading/decoding response from PayPal
	//-----------------------------------
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	responseValues, err := url.ParseQuery(string(responseBody))
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "550c824b-86cf-4c8d-a13e-73f92da15bde", "payment unexpected response")
	}

	// Check all values in response for valid credit card data
	if _, ccSecureCodePresent := ccInfo["cvv"]; ccSecureCodePresent {
		if utils.InterfaceToString(responseValues.Get("CVV2MATCH")) != "Y" {
			// invalid CVV2
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "51d1a2c9-2f0a-4eee-9aa2-527ca6d83f28", "Payment validation error: CVV2 code incorect")
		}
	}

	result := map[string]interface{}{
		"transactionID":      responseValues.Get("PNREF"),
		"responseMessage":    responseValues.Get("RESPMSG"),
		"responseResult":     responseValues.Get("RESULT"),
		"creditCardLastFour": responseValues.Get("ACCT"),
		"creditCardType":     responseValues.Get("CARDTYPE"),
	}

	// utils.InterfaceToString(result["transactionID"])
	// if status is ok return result with valid values
	if utils.InterfaceToString(result["transactionID"]) != "" {
		if utils.InterfaceToString(result["responseMessage"]) == "Verified" {
			return result, nil
		}

		// On review of by Fraud Service -- possible to continue
		if utils.InterfaceToString(result["responseResult"]) == "126" {
			env.Log("paypal.log", env.ConstLogPrefixInfo, "ZERO AMOUNT ATHORIZE TRANSACTION WITH COMMENT: "+
				"MESSAGE - "+utils.InterfaceToString(result["responseMessage"])+
				"TRANSACTIONID - "+utils.InterfaceToString(result["transactionID"]))

			return result, nil
		}

		env.Log("paypal.log", env.ConstLogPrefixInfo, "ZERO AMOUNT ATHORIZE FAIL: "+
			"MESSAGE - "+utils.InterfaceToString(result["responseMessage"])+
			"RESULT - "+utils.InterfaceToString(result["responseResult"]))

		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a050604a-b9e9-44cc-a4d1-e5c0bfab5c69", "Payment error: "+utils.InterfaceToString(result["responseMessage"]))
	}

	env.Log("paypal.log", env.ConstLogPrefixInfo, "ZERO AMOUNT ATHORIZE FAIL: "+
		"MESSAGE - "+utils.InterfaceToString(result["responseMessage"])+
		"RESULT - "+utils.InterfaceToString(result["responseResult"]))

	return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a050604a-b9e9-44cc-a4d1-e5c0bfab5c69", "Payment error: "+utils.InterfaceToString(result["responseMessage"]))
}
