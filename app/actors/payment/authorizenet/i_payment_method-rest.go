package authorizenet

import (
	"fmt"
	"time"
	"strings"
	"regexp"

	"github.com/hunterlong/authorizecim"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)


// GetInternalName returns the name of the payment method
func (it RestAPI) GetInternalName() string {
	return ConstPaymentAuthorizeNetRestApiName
}

// GetName returns the user customized name of the payment method
func (it *RestAPI) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAuthorizeNetRestApiTitle))
}

// GetCode returns payment method code
func (it *RestAPI) GetCode() string {
	return ConstPaymentAuthorizeNetRestApiCode
}

// IsTokenable returns possibility to save token for this payment method
func (it *RestAPI) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return true
}

// GetType returns type of payment method
func (it *RestAPI) GetType() string {
	return checkout.ConstPaymentTypeCreditCard
}

// IsAllowed checks for method applicability
func (it *RestAPI) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathAuthorizeNetRestApiEnabled))
}

// Authorize makes payment method authorize operation
func (it *RestAPI) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	_, err := it.ConnectToAuthorize()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var profileId = ""
	var paymentId = ""

	action := paymentInfo[checkout.ConstPaymentActionTypeKey]
	isCreateToken := utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken
	if isCreateToken {
		ccInfo := utils.InterfaceToMap(paymentInfo["cc"])

		if profileId == "" {

			newProfileId, err := it.CreateProfile(orderInstance, paymentInfo)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			profileId = newProfileId
		}

		if profileId != "" {
			// 3. Create a card
			newPaymentID, _, err := it.CreatePaymentProfile(orderInstance, paymentInfo, profileId)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			paymentId = newPaymentID
			numberString := utils.InterfaceToString(ccInfo["number"])
			cardType, err := getCardTypeByNumber(utils.InterfaceToString(numberString))
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			// This response looks like our normal authorize response
			// but this map is translated into other keys to store a token
			result := map[string]interface{}{
				"transactionID":      paymentId, // transactionID
				"creditCardLastFour": numberString[len(numberString)-4:], // number
				"creditCardType":     cardType, // type
				"creditCardExp":      utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]), // expiration_date
				"customerID":         profileId, // customer_id
			}
			fmt.Println(result)

			return result, nil
		}
	}

	creditCard, creditCardOk := paymentInfo["cc"].(visitor.InterfaceVisitorCard);
	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	if  creditCardOk && creditCard != nil {
		profileId = creditCard.GetCustomerID()
		paymentId = creditCard.GetToken()
	}

	if utils.InterfaceToBool(ccInfo["save"]) != true && profileId == "" && paymentId == "" {
		return it.AuthorizeWithoutSave(orderInstance, paymentInfo)
	}
	if paymentId != "" && profileId != "" {

		// Waiting for 5 seconds to allow Authorize.net to keep up
		time.Sleep(5000 * time.Millisecond)
		grandTotal := orderInstance.GetGrandTotal()
		amount := fmt.Sprintf("%.2f", grandTotal)

		item := AuthorizeCIM.LineItem{
			ItemID: orderInstance.GetID(),
			Name: "Order #" + orderInstance.GetID(),
			Description: "",
			Quantity: "1",
			UnitPrice: amount,
		}

		response, approved, success := AuthorizeCIM.CreateTransaction(profileId, paymentId, item, amount)
		// outputs transaction response, approved status (true/false), and success status (true/false)
		var orderTransactionID string
		if !success {
			env.Log("authorizenet.log", env.ConstLogPrefixInfo, "Transaction has failed: "+fmt.Sprint(response))
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "da966f67-666f-412c-a381-a080edd915d0", checkout.ConstPaymentErrorTechnical)
		}

		orderTransactionID = response["transId"].(string)
		status := "denied"
		if approved {
			status = "approved"
		}

		env.Log("authorizenet.log", env.ConstLogPrefixInfo, "NEW TRANSACTION ("+status+"): "+
			"Visitor ID - "+utils.InterfaceToString(orderInstance.Get("visitor_id"))+", "+
			"LASTNAME - "+orderInstance.GetBillingAddress().GetLastName()+", "+
			"Order ID - "+utils.InterfaceToString(orderInstance.GetID())+", "+
			"TRANSACTIONID - "+orderTransactionID)


		// This response looks like our normal authorize response
		// but this map is translated into other keys to store a token
		result := map[string]interface{}{
			"transactionID":      response["transId"].(string), // transactionID
			"creditCardLastFour": strings.Replace(response["accountNumber"].(string), "XXXX", "", -1), // number
			"creditCardType":     response["accountType"].(string), // type
			"creditCardExp":      utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]), // expiration_date
			"customerID":         profileId, // customer_id
			"tokenID":            paymentId, // token_id
		}

		if !creditCardOk {
			_, err := it.SaveToken(orderInstance, result)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}

		return result, nil
	}

	return nil, nil
}

func (it *RestAPI) AuthorizeWithoutSave(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	ccCVC := utils.InterfaceToString(ccInfo["cvc"])
	if ccCVC == "" {
		err := env.ErrorNew(ConstErrorModule, 1, "fdcb2ecd-a31d-4fa7-a4e8-df51e10a5332", "CVC field was left empty")
		return nil, err
	}

	grandTotal := orderInstance.GetGrandTotal()
	amount := fmt.Sprintf("%.2f", grandTotal)

	credit_card := AuthorizeCIM.CreditCardCVV{
		CardNumber: utils.InterfaceToString(ccInfo["number"]),
		ExpirationDate: utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]),
		CardCode: ccCVC,
	}

	response, approved, success := AuthorizeCIM.AuthorizeCard(credit_card, amount)
	// outputs transaction response, approved status (true/false), and success status (true/false)

	var orderTransactionID string
	if !success {
		env.Log("authorizenet.log", env.ConstLogPrefixInfo, "Transaction has failed: "+fmt.Sprint(response))
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "48352140-873c-4deb-8cf0-b3140225d8fb", checkout.ConstPaymentErrorTechnical)
	}

	status := "denied"
	if approved {
		status = "approved"
	}

	env.Log("authorizenet.log", env.ConstLogPrefixInfo, "NEW TRANSACTION ("+status+"): "+
		"Visitor ID - "+utils.InterfaceToString(orderInstance.Get("visitor_id"))+", "+
		"LASTNAME - "+orderInstance.GetBillingAddress().GetLastName()+", "+
		"Order ID - "+utils.InterfaceToString(orderInstance.GetID())+", "+
		"TRANSACTIONID - "+orderTransactionID)

	// This response looks like our normal authorize response
	// but this map is translated into other keys to store a token
	result := map[string]interface{}{
		"transactionID":      response["transId"].(string), // token_id
		"creditCardLastFour": strings.Replace(response["accountNumber"].(string), "XXXX", "", -1), // number
		"creditCardType":     response["accountType"].(string), // type
		"creditCardExp":      utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]), // expiration_date
		"customerID":         0, // customer_id
	}

	return result, nil

}

func (it *RestAPI) CreateProfile(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (string, error) {
	profileId := ""
	extra := utils.InterfaceToMap(paymentInfo["extra"])
	userEmail := utils.InterfaceToString(extra["email"])
	billingName := utils.InterfaceToString(extra["billing_name"])

	customerInfo := AuthorizeCIM.AuthUser{
		"0",
		userEmail,
		billingName,
	}

	newProfileId, response, success := AuthorizeCIM.CreateCustomerProfile(customerInfo)
	response = utils.InterfaceToMap(response)
	if success {
		profileId = newProfileId

		if orderInstance != nil {
			env.Log("authorizenet.log", env.ConstLogPrefixInfo, "New Customer Profile: " +
				"Visitor ID - " + utils.InterfaceToString(orderInstance.Get("visitor_id")) + ", " +
				"BILLNAME - " + billingName + ", " +
				"Order ID - " + utils.InterfaceToString(orderInstance.GetID()) + ", " +
				"Profile ID - " + profileId)
		}
	} else {
		messages, _ := response["messages"].(map[string]interface{})
		if messages != nil {
			// Array
			messageArray, _ := messages["message"].([]interface{})
			// Hash
			text := (messageArray[0].(map[string]interface{}))["text"]

			re := regexp.MustCompile("[0-9]+")
			profileId = re.FindString(text.(string))
		}

	}

	if profileId == "" || profileId == "0" {
		return "", env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "221aaa5a-a87e-4dc3-a1a9-a8cfee975f48", "profileId can't be empty")
	}

	return profileId, nil
}


func (it *RestAPI) CreatePaymentProfile(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}, profileId string) (string, map[string]interface{}, error) {
	paymentID := ""
	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	var address AuthorizeCIM.Address
	if orderInstance != nil {
		address = AuthorizeCIM.Address{
			FirstName: orderInstance.GetBillingAddress().GetFirstName(),
			LastName: orderInstance.GetBillingAddress().GetLastName(),
			Address: orderInstance.GetBillingAddress().GetAddress(),
			City: orderInstance.GetBillingAddress().GetCity(),
			State: orderInstance.GetBillingAddress().GetState(),
			Zip: orderInstance.GetBillingAddress().GetZipCode(),
			Country: orderInstance.GetBillingAddress().GetCountry(),
			PhoneNumber:  orderInstance.GetBillingAddress().GetPhone(),
		}
	} else {
		address = AuthorizeCIM.Address{
			FirstName: "",
			LastName: "",
			Address: "",
			City: "",
			State: "",
			Zip: "",
			Country: "",
			PhoneNumber:  "",
		}
	}
	credit_card := AuthorizeCIM.CreditCard{
		CardNumber: utils.InterfaceToString(ccInfo["number"]),
		ExpirationDate: utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]),
	}

	newPaymentID, response, success := AuthorizeCIM.CreateCustomerBillingProfile(profileId, credit_card, address)
	response = utils.InterfaceToMap(response)
	if success {
		paymentID = newPaymentID

		if orderInstance != nil {
			env.Log("authorizenet.log", env.ConstLogPrefixInfo, "New Credit Card was added: "+
				"Visitor ID - "+utils.InterfaceToString(orderInstance.Get("visitor_id"))+", "+
				"LASTNAME - "+orderInstance.GetBillingAddress().GetLastName()+", "+
				"Order ID - "+utils.InterfaceToString(orderInstance.GetID())+", "+
				"Billing ID - "+paymentID)
		}

	} else {
		messages, _ := response["messages"].(map[string]interface{})
		if messages != nil {
			// Array
			messageArray, _ := messages["message"].([]interface{})
			// Hash
			text := (messageArray[0].(map[string]interface{}))["text"]
			return "", response, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5609f3bf-bad6-4e93-8d1e-bf525ddf17f9", text.(string))
		}
		env.Log("authorizenet.log", env.ConstLogPrefixInfo, "There was an issue inserting a credit card into the user account")
	}

	if paymentID == "" || paymentID == "0" {
		return "", response, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5609f3bf-bad6-4e93-8d1e-bf525ddf17f9", "paymentID can't be empty")
	}

	return paymentID, response, nil
}

func (it *RestAPI) SaveToken(orderInstance order.InterfaceOrder, creditCardInfo map[string]interface{}) (visitor.InterfaceVisitorCard, error) {

	visitorID := utils.InterfaceToString(orderInstance.Get("visitor_id"))

	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, 1, "d43b4347-7560-4432-a9b3-b6941693f77f", "CVC field was left empty")
	}

	authorizeCardResult := utils.InterfaceToMap(creditCardInfo)
	if !utils.KeysInMapAndNotBlank(authorizeCardResult, "transactionID", "creditCardLastFour") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "22e17290-56f3-452a-8d54-18d5a9eb2833", "transaction can't be obtained")
	}

	// create visitor card and fill required fields
	//---------------------------------
	visitorCardModel, err := visitor.GetVisitorCardModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// create credit card map with info
	tokenRecord := map[string]interface{}{
		"visitor_id":      visitorID,
		"payment":         it.GetCode(),
		"type":            authorizeCardResult["creditCardType"],
		"number":          authorizeCardResult["creditCardLastFour"],
		"expiration_date": authorizeCardResult["creditCardExp"],
		"holder":          utils.InterfaceToString(authorizeCardResult["holder"]),
		"token_id":        authorizeCardResult["tokenID"],
		"customer_id":     authorizeCardResult["customerID"],
		"token_updated":   time.Now(),
		"created_at":      time.Now(),
	}

	err = visitorCardModel.FromHashMap(tokenRecord)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorCardModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorCardModel, nil
}


func (it *RestAPI) ConnectToAuthorize() (bool, error) {
	var apiLoginId = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAuthorizeNetRestApiApiLoginId))
	if apiLoginId == "" {
		return false, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "88111f54-e8a1-4c43-bc38-0e660c4caa16", "api login id was not specified")
	}

	var transactionKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAuthorizeNetRestApiTransactionKey))
	if transactionKey == "" {
		return false, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "35de21dd-3f07-4ec2-9630-a15fa07d00a5", "transaction key was not specified")
	}

	var mode = ""
	var isTestMode = utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathAuthorizeNetRestApiTest))
	if isTestMode {
		mode = "test"
	}

	AuthorizeCIM.SetAPIInfo(apiLoginId, transactionKey, mode)

	connected := AuthorizeCIM.TestConnection()
	if !connected {
		return false, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4faf7f78-cda7-464f-9a9e-459806907069", "cannot connect to Authorize.Net")
	}

	return true, nil
}

// Capture makes payment method capture operation
func (it *RestAPI) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ebbac9ac-94e3-48f7-ae8a-8a562ee09907", "Not implemented")
}

// Refund will return funds on the given order :: Not Implemented Yet
func (it *RestAPI) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "baaf0cac-2924-4340-a9a1-cc3e407326d3", "Not implemented")
}

// Void will mark the order and capture as void
func (it *RestAPI) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "eb391185-161d-4e0f-8d08-470dda867fed", "Not implemented")
}
