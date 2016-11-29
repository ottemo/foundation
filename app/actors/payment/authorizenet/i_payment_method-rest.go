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

	profileId := ""
	paymentId := ""
	creditCard, creditCardOk := paymentInfo["cc"].(visitor.InterfaceVisitorCard);

	if  creditCardOk && creditCard != nil {
		profileId = creditCard.GetCustomerID()
		paymentId = creditCard.GetToken()
	}
	if utils.InterfaceToBool(ccInfo["save"]) != true && profileId == "" && paymentId == "" {
		return it.AuthorizeWithoutSave(orderInstance, paymentInfo)
	}

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
fmt.Println(profileId)
fmt.Println(paymentId)
	if profileId != "0"  && profileId != "" && paymentId == "" {
		// 3. Create a card

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

		paymentId = newPaymentID
	}

	if paymentId != "" && profileId != "" {
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
			"transactionID":      response["transId"].(string), // transactionID
			"creditCardLastFour": strings.Replace(response["accountNumber"].(string), "XXXX", "", -1), // number
			"creditCardType":     response["accountType"].(string), // type
			"creditCardExp":      utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]), // expiration_date
			"customerID":         profileId, // customer_id
			"tokenID":         paymentId, // token_id
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
		return nil, env.ErrorNew(ConstErrorModule, 1, "d43b4347-7560-4432-a9b3-b6941693f77f", "CVC field was left empty")
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
