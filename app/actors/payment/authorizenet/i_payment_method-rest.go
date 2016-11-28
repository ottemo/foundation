package authorizenet

import (
	"github.com/hunterlong/authorizecim"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"fmt"
	"time"
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

	AuthorizeCIM.SetAPIInfo(apiLoginId, transactionKey, "test")

	connected := AuthorizeCIM.TestConnection()
	if !connected {
		fmt.Println("There was an issue connecting to Authorize.net")
	}

	// Check if we are just supposed to create a Customer (aka a token)
	//action := paymentInfo[checkout.ConstPaymentActionTypeKey]
	//isCreateToken := utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken
	//if isCreateToken {
	// NOTE: `orderInstance = nil` when creating a token

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
	extra := utils.InterfaceToMap(paymentInfo["extra"])
	visitorID := utils.InterfaceToString(extra["visitor_id"])
	profileId := getAuthorizenetCustomerToken(visitorID)
	if profileId == "" {

		userEmail := utils.InterfaceToString(extra["email"])
		billingName := utils.InterfaceToString(extra["billing_name"])

		customerInfo := AuthorizeCIM.AuthUser{
			"0",
			userEmail,
			billingName,
		}

		profileNewId, success := AuthorizeCIM.CreateCustomerProfile(customerInfo)

		if success {
			fmt.Println("New Customer Profile ID: ", profileId + "\n")
		} else {
			fmt.Println("There was an issue creating the Customer Profile \n")
		}

		profileId = profileNewId

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

		fmt.Println(response)

		// This response looks like our normal authorize response
		// but this map is translated into other keys to store a token
		result := map[string]interface{}{
			"transactionID":      response["transId"].(string), // token_id
			"creditCardLastFour": response["accountNumber"].(string), // number
			"creditCardType":     response["accountType"].(string), // type
			"creditCardExp":      utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]), // expiration_date
			"customerID":         profileId, // customer_id
		}

		fmt.Println(result)

		return result, nil
	}
	//
	////billingFirstName := orderInstance.GetBillingAddress().GetFirstName()
	//
	return nil, nil
	//return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ed753163-d708-4884-aae8-3aa1dc9bf9f4", "Not implemented")

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

	fmt.Println(response)

	// This response looks like our normal authorize response
	// but this map is translated into other keys to store a token
	result := map[string]interface{}{
		"transactionID":      response["transId"].(string), // token_id
		"creditCardLastFour": response["accountNumber"].(string), // number
		"creditCardType":     response["accountType"].(string), // type
		"creditCardExp":      utils.InterfaceToString(ccInfo["expire_year"]) + "-" + utils.InterfaceToString(ccInfo["expire_month"]), // expiration_date
		"customerID":         0, // customer_id
	}

	fmt.Println(result)

	return result, nil

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
