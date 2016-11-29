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

	fmt.Println("Authorize \n\n\n\n")

	var apiLoginId = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAuthorizeNetRestApiApiLoginId))
	if apiLoginId == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "88111f54-e8a1-4c43-bc38-0e660c4caa16", "account id was not specified")
	}

	var transactionKey = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAuthorizeNetRestApiTransactionKey))
	if transactionKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "88111f54-e8a1-4c43-bc38-0e660c4caa16", "account id was not specified")
	}

	// @todo check if test mode
	AuthorizeCIM.SetAPIInfo(apiLoginId, transactionKey, "test")

	connected := AuthorizeCIM.TestConnection()
	if !connected {
		fmt.Println("There was an issue connecting to Authorize.net")
	}

	action := paymentInfo[checkout.ConstPaymentActionTypeKey]
	isCreateToken := utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken
	if isCreateToken {
		return nil, nil
	}

	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	if utils.InterfaceToBool(ccInfo["save"]) != true {
		return it.AuthorizeWithoutSave(orderInstance, paymentInfo)
	}

	// 1. Get our customer token
	visitorID := utils.InterfaceToString(orderInstance.Get("visitor_id"))

	profileId := getAuthorizenetCustomerToken(visitorID, it.GetCode())
	if profileId == "" {

		extra := utils.InterfaceToMap(paymentInfo["extra"])
		userEmail := utils.InterfaceToString(extra["email"])
		billingName := utils.InterfaceToString(extra["billing_name"])

		customerInfo := AuthorizeCIM.AuthUser{
			"0",
			userEmail,
			billingName,
		}

		profileNewId, response, success := AuthorizeCIM.CreateCustomerProfile(customerInfo)

		if success {
			profileId = profileNewId

			fmt.Println("New Customer Profile ID: ", profileId + "\n")
		} else {
			fmt.Println("There was an issue creating the Customer Profile \n")
			fmt.Println(response)
			fmt.Println(utils.InterfaceToString(response))
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


	}

	if profileId != "0"  && profileId != "" {
		// 3. Create a card
		ccCVC := utils.InterfaceToString(ccInfo["cvc"])
		if ccCVC == "" {
			return nil, env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
		}

		address := AuthorizeCIM.Address{
			FirstName: orderInstance.GetBillingAddress().GetFirstName(),
			LastName: orderInstance.GetBillingAddress().GetLastName(),
			Address: orderInstance.GetBillingAddress().GetAddress(),
			City: orderInstance.GetBillingAddress().GetCity(),
			State: orderInstance.GetBillingAddress().GetState(),
			Zip: orderInstance.GetBillingAddress().GetZipCode(),
			Country: orderInstance.GetBillingAddress().GetCountry(),
			PhoneNumber:  orderInstance.GetBillingAddress().GetPhone(),
		}
		credit_card := AuthorizeCIM.CreditCard{
			CardNumber: utils.InterfaceToString(ccInfo["number"]),
			ExpirationDate: utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]),
		}

		newPaymentID, success := AuthorizeCIM.CreateCustomerBillingProfile(profileId, credit_card, address)
		if success {
			fmt.Println("New Credit Card was added, Billing ID: ", newPaymentID)
		} else {
			fmt.Println("There was an issue inserting a credit card into the user account")
		}

		fmt.Println("Waiting for 10 seconds to allow Authorize.net to keep up")
		time.Sleep(10000 * time.Millisecond)

		paymentId := newPaymentID

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
		var tranxID string
		if success {
			tranxID = response["transId"].(string)
			if approved {
				fmt.Println("Transaction was approved! " + tranxID + "\n")
			} else {
				fmt.Println("Transaction was denied! " + tranxID + "\n")
			}
		} else {
			fmt.Println("Transaction has failed! \n")
		}

		// This response looks like our normal authorize response
		// but this map is translated into other keys to store a token
		result := map[string]interface{}{
			"transactionID":      response["transId"].(string), // token_id
			"creditCardLastFour": strings.Replace(response["accountNumber"].(string), "XXXX", "", -1), // number
			"creditCardType":     response["accountType"].(string), // type
			"creditCardExp":      utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]), // expiration_date
			"customerID":         profileId, // customer_id
		}

		_, err := it.SaveToken(orderInstance, result)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		return result, nil
	}

	return nil, nil

}

func (it *RestAPI) AuthorizeWithoutSave(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	ccCVC := utils.InterfaceToString(ccInfo["cvc"])
	if ccCVC == "" {
		err := env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
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

	var tranxID string
	if success {
		tranxID = response["transId"].(string)
		if approved {
			fmt.Println("Transaction was approved! " + tranxID + "\n")
		} else {
			fmt.Println("Transaction was denied! " + tranxID + "\n")
		}
	} else {
		fmt.Println("Transaction has failed! \n")
	}

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

func (it *RestAPI) SaveToken(orderInstance order.InterfaceOrder, creditCardInfo map[string]interface{}) (visitor.InterfaceVisitorCard, error) {

	fmt.Println("\n\n\n\n\n\nSaveToken \n\n\n\n\n\n\n")
	visitorID := utils.InterfaceToString(orderInstance.Get("visitor_id"))
	fmt.Println(visitorID)
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
	}

	authorizeCardResult := utils.InterfaceToMap(creditCardInfo)
	if !utils.KeysInMapAndNotBlank(authorizeCardResult, "transactionID", "creditCardLastFour") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "22e17290-56f3-452a-8d54-18d5a9eb2833", "transaction can't be obtained")
	}
	fmt.Println(authorizeCardResult)
	// create visitor card and fill required fields
	//---------------------------------
	visitorCardModel, err := visitor.GetVisitorCardModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// override credit card info with provided from payment info
	creditCardInfo["token_id"] = authorizeCardResult["transactionID"]
	creditCardInfo["payment"] = it.GetCode()
	creditCardInfo["customer_id"] = authorizeCardResult["customerID"]
	creditCardInfo["type"] = authorizeCardResult["creditCardType"]
	creditCardInfo["number"] = authorizeCardResult["creditCardLastFour"]
	creditCardInfo["expiration_date"] = authorizeCardResult["creditCardExp"] // mmyy
	creditCardInfo["token_updated"] = time.Now()
	creditCardInfo["created_at"] = time.Now()
	fmt.Println(creditCardInfo)
	// filling new instance with request provided data
	for attribute, value := range creditCardInfo {
		err := visitorCardModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// setting credit card owner to current visitor (for sure)
	visitorCardModel.Set("visitor_id", visitorID)

	fmt.Println(visitorCardModel.ToHashMap())

	// save card info if checkbox is checked on frontend
	if utils.InterfaceToBool(creditCardInfo["save"]) {
		err = visitorCardModel.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}
	fmt.Println(visitorCardModel.ToHashMap())
	return visitorCardModel, nil
}

// Capture makes payment method capture operation
func (it *RestAPI) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ed753163-d708-4884-aae8-3aa1dc9bf9f4", "Not implemented")
}

// Refund will return funds on the given order :: Not Implemented Yet
func (it *RestAPI) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2dc9b523-4e53-4ff5-9d49-1dfdadf3fb44", "Not implemented")
}

// Void will mark the order and capture as void
func (it *RestAPI) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d682f87e-4d51-473b-a00d-191d28e807f5", "Not implemented")
}
