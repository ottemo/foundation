package paypal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
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

// IsTokenable checks for method applicability
func (it *PayFlowAPI) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathPayPalPayflowEnabled)) && utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathPayPalPayflowTokenable))
}

// IsAllowed checks for method applicability
func (it *PayFlowAPI) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathPayPalPayflowEnabled))
}

// Authorize makes payment method authorize operation (currently it's a Authorize zero amount + Sale operations)
func (it *PayFlowAPI) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	// make only authorize zero amount procedure to obtain and return it's token with additional payment info
	if value, present := paymentInfo[checkout.ConstPaymentActionTypeKey]; present && utils.InterfaceToString(value) == checkout.ConstPaymentActionTypeCreateToken {
		return it.AuthorizeZeroAmount(orderInstance, paymentInfo)

		// currently we will do request directly from foundation side but our next step is to do it just create request content
		// and make post from storefront
		//		return it.CreateAuthorizeZeroAmountRequest(orderInstance, paymentInfo)
	}

	var transactionID string
	var visitorCreditCard visitor.InterfaceVisitorCard

	// try to obtain visitor token info
	if cc, present := paymentInfo["cc"]; present {
		if creditCard, ok := cc.(visitor.InterfaceVisitorCard); ok && creditCard != nil {
			transactionID = creditCard.GetToken()
			visitorCreditCard = creditCard
		}
	}

	// otherwise we will authorize in default way
	if transactionID == "" {
		authorizeZeroResult, err := it.AuthorizeZeroAmount(orderInstance, paymentInfo)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		authorizeZeroResultMap := utils.InterfaceToMap(authorizeZeroResult)
		if value, present := authorizeZeroResultMap["transactionID"]; !present || utils.InterfaceToString(value) == "" {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5e68f079-e8ce-4677-8fb9-89c6f7acbd7f", "Error: token was not created")
		}

		transactionID = utils.InterfaceToString(authorizeZeroResultMap["transactionID"])
	}

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
		return nil, env.ErrorDispatch(err)
	}

	responseValues, err := url.ParseQuery(string(responseBody))
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b18cdcad-8c21-4acf-a2e0-56e0541103de", "payment unexpected response")
	}

	if responseValues.Get("RESPMSG") != "Approved" {
		env.Log("paypal.log", env.ConstLogPrefixInfo, "Redjected payment: "+fmt.Sprint(responseValues))
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e48403bb-c15d-4302-8894-da7146b93260", checkout.ConstPaymentErrorDeclined+" Reason: "+responseValues.Get("RESPMSG")+", "+responseValues.Get("PREFPSMSG"))
	}

	// get info about transaction from payment response
	orderTransactionID := utils.InterfaceToString(responseValues.Get("PNREF"))
	if orderTransactionID == "" {
		env.Log("paypal.log", env.ConstLogPrefixInfo, "Redjected payment: "+fmt.Sprint(responseValues))
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d1d0a2d6-786a-4a29-abb1-3eb7667fbc3e", checkout.ConstPaymentErrorTechnical+" Reason: "+responseValues.Get("RESPMSG")+". "+responseValues.Get("PREFPSMSG"))
	}

	env.Log("paypal.log", env.ConstLogPrefixInfo, "NEW TRANSACTION: "+
		"Visitor ID - "+utils.InterfaceToString(orderInstance.Get("visitor_id"))+", "+
		"Order ID - "+utils.InterfaceToString(orderInstance.GetID())+", "+
		"TRANSACTIONID - "+orderTransactionID)

	orderPaymentInfo := map[string]interface{}{
		"transactionID":     orderTransactionID,
		"creditCardNumbers": responseValues.Get("ACCT"),
		"creditCardExp":     responseValues.Get("EXPDATE"),
		"creditCardType":    getCreditCardName(utils.InterfaceToString(responseValues.Get("CARDTYPE"))),
	}

	// id presense in credit card means it was saved so we can update token for it
	if visitorCreditCard != nil && visitorCreditCard.GetID() != "" {
		orderPaymentInfo["creditCardID"] = visitorCreditCard.GetID()

		visitorCreditCard.Set("token_id", orderTransactionID)
		visitorCreditCard.Set("token_updated", time.Now())
		visitorCreditCard.Save()
	}

	return orderPaymentInfo, nil
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
		env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "Can't obtain secure token: "+fmt.Sprint(responseValues))
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f3608dfb-3c7a-4549-82c1-83d6e9d8b7cb", checkout.ConstPaymentErrorTechnical+" "+responseValues.Get("RESPMSG"))
	}

	token := responseValues.Get("SECURETOKEN")
	if responseValues.Get("SECURETOKENID") != secureTokenID {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9b095f62-b371-4eaf-965f-98eb24206e53", checkout.ConstPaymentErrorTechnical+" Unexpected response, SECURETOKENID value changed")
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

		// Payment Details Fields
		"&AMT=0" +
		"&VERBOSITY=HIGH"

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
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "550c824b-86cf-4c8d-a13e-73f92da15bde", checkout.ConstPaymentErrorTechnical+" Payment unexpected response")
	}
	responseResult := utils.InterfaceToString(responseValues.Get("RESULT"))
	responseMessage := utils.InterfaceToString(responseValues.Get("RESPMSG"))
	transactionID := utils.InterfaceToString(responseValues.Get("PNREF"))

	if responseResult == "" && responseMessage == "" || len(responseValues) == 0 {
		env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "TRANSACTION NO RESPONSE: "+
			"RESPONSE - "+fmt.Sprint(responseValues))

		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4d941690-d981-4d20-9b4e-ab903d1ea526", checkout.ConstPaymentErrorTechnical+" The payment server is not responding.")
	}

	result := map[string]interface{}{
		"transactionID":      transactionID,
		"responseMessage":    responseMessage,
		"responseResult":     responseResult,
		"creditCardLastFour": responseValues.Get("ACCT"),
		"creditCardType":     getCreditCardName(responseValues.Get("CARDTYPE")),
		"creditCardExp":      responseValues.Get("EXPDATE"),
	}

	// Check all values in response for valid credit card data
	if _, ccSecureCodePresent := ccInfo["cvv"]; ccSecureCodePresent {
		if utils.InterfaceToString(responseValues.Get("CVV2MATCH")) == "N" {
			// invalid CVV2
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "51d1a2c9-2f0a-4eee-9aa2-527ca6d83f28", checkout.ConstPaymentErrorDeclined+" Reason: CVV code is not valid.")
		}
	}

	// utils.InterfaceToString(result["transactionID"])
	// if status is ok and card is Verified - return result with valid values
	if transactionID != "" {
		if responseResult == "0" && responseMessage[0:8] == "Verified" {
			return result, nil
		}

		// On review of by Fraud Service -- possible to continue
		if responseResult == "126" {
			env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "ZERO AMOUNT ATHORIZE TRANSACTION WITH COMMENT: "+
				"MESSAGE - "+responseMessage+
				"TRANSACTIONID - "+transactionID)

			return result, nil
		}
	}

	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "ZERO AMOUNT ATHORIZE FAIL: "+
		"MESSAGE - "+responseMessage+" "+
		"RESULT - "+responseResult)

	return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a050604a-b9e9-44cc-a4d1-e5c0bfab5c69", checkout.ConstPaymentErrorDeclined+"Reason: "+responseMessage)
}

// CreateAuthorizeZeroAmountRequest will do Account Verification and return transaction ID for refer transaction if all info is valid
func (it *PayFlowAPI) CreateAuthorizeZeroAmountRequest(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

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
		"&ACCT=" + "$CC_NUM" +
		"&EXPDATE=" + "$CC_MONTH$CC_YEAR" +

		// Payment Details Fields
		"&AMT=0" +
		"&VERBOSITY=HIGH"

	nvpGateway := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowURL))

	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "NEW  obtain token transaction request created")
	env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "Params: "+requestParams)

	return api.StructRestRedirect{Result: requestParams, Location: nvpGateway}, nil
}
