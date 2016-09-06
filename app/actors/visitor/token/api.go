package token

import (
	"fmt"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.POST("visit/token", APICreateToken)
	service.GET("visit/tokens", APIListVisitorCards)
	service.DELETE("visit/token/:tokenID", APIDeleteToken)

	//service.PUT("token/:tokenID", APIGetToken)

	return nil
}


// APICreateToken creates a request body for posting credit card info to payment system with 0 amount payment
// for obtaining token on this card and saving it for visitor
func APICreateToken(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return "You are not logged in, please log in.", nil
	}

	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Context contains: %v\n\n\n", requestData)

	paymentMethodCode := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "payment", "payment_method"))
	if paymentMethodCode == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6d1691c8-2d26-44be-b90d-24d920e26301", "Please select a payment method.")
	}

	value, present := requestData["cc"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2e9f1bfc-ec9f-4017-83c6-4d04b95b9c08", "Missing field in credit card data.")
	}

	creditCardInfo := utils.InterfaceToMap(value)

	var paymentMethod checkout.InterfacePaymentMethod

	for _, payment := range checkout.GetRegisteredPaymentMethods() {
		if payment.GetCode() == paymentMethodCode {
			if !payment.IsTokenable(nil) {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "519ef43c-4d07-4b64-90f7-7fdc3657940a", "Cannot save selected Credit Card.")
			}
			paymentMethod = payment
		}
	}

	if paymentMethod == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c80c4106-1208-4d0b-8577-0889f608869b", "Provided payment method does not exist.")
	}

	// add email, first and lastname to paymentInfo map
	email := currentCheckout.GetVisitor().GetEmail()
	billingFirstName := currentCheckout.GetBillingAddress().GetFirstName()
	billingLastName := currentCheckout.GetBillingAddress().GetLastName()

	paymentInfo := map[string]interface{}{
		checkout.ConstPaymentActionTypeKey: checkout.ConstPaymentActionTypeCreateToken,
		"cc": creditCardInfo,
		"extra": map[string]interface{}{
			"email":        email,
			"billing_name": billingFirstName + " " + billingLastName,
		},
	}
	// contains creditCardLastFour, creditCardType, responseMessage, responseResult, transactionID, creditCardExp
	paymentResult, err := paymentMethod.Authorize(nil, paymentInfo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cardInfoMap := utils.InterfaceToMap(paymentResult)
	if !utils.KeysInMapAndNotBlank(cardInfoMap, "transactionID") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "22e17290-56f3-452a-8d54-18d5a9eb2833", "A transaction ID was not provided.")
	}

	// create visitor address operation
	//---------------------------------
	visitorCardModel, err := visitor.GetVisitorCardModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// create credit card map with info
	tokenRecord := map[string]interface{}{
		"visitor_id":      visitorID,
		"payment":         paymentMethodCode,
		"type":            cardInfoMap["creditCardType"],
		"number":          cardInfoMap["creditCardLastFour"],
		"expiration_date": cardInfoMap["creditCardExp"],
		"holder":          utils.InterfaceToString(requestData["holder"]),
		"token_id":        cardInfoMap["transactionID"],
		"token_updated":   time.Now(),
	}

	err = visitorCardModel.FromHashMap(tokenRecord)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorCardModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorCardModel.ToHashMap(), nil
}

// APIListVisitorCards return a list of existing tokens for visitor
func APIListVisitorCards(context api.InterfaceApplicationContext) (interface{}, error) {

	// if visitorID was specified - using this otherwise, taking current visitor
	visitorID := context.GetRequestArgument("visitorID")
	if visitorID == "" {

		sessionVisitorID := visitor.GetCurrentVisitorID(context)
		if sessionVisitorID == "" {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2ac4c16b-9241-406e-b35a-399813bb6ca5", "Please log in.")
		}
		visitorID = sessionVisitorID
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		if visitorID != visitor.GetCurrentVisitorID(context) {
			return nil, env.ErrorDispatch(err)
		}
	}

	// list operation
	//---------------
	visitorCardCollectionModel, err := visitor.GetVisitorCardCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := visitorCardCollectionModel.GetDBCollection()
	dbCollection.AddStaticFilter("visitor_id", "=", visitorID)

	// add allowed payment methods filter
	currentCheckout, err := checkout.GetCurrentCheckout(context, false)
	isSubscription := currentCheckout.IsSubscription()

	paymentMethods := make([]string, 0)
	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.IsAllowed(currentCheckout) && (!isSubscription || paymentMethod.IsTokenable(currentCheckout)) {
			paymentMethods = append(paymentMethods, paymentMethod.GetCode())
		}
	}
	dbCollection.AddStaticFilter("payment", "IN", paymentMethods)

	// filters handle
	models.ApplyFilters(context, dbCollection)

	// checking for a "count" request
	if context.GetRequestArgument("count") != "" {
		return visitorCardCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	visitorCardCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, visitorCardCollectionModel)


	return visitorCardCollectionModel.List()
}

// APIDeleteToken deletes credit card token by provided token_id
func APIDeleteToken (context api.InterfaceApplicationContext) (interface{}, error) {

	tokenID := context.GetRequestArgument("tokenID")
	if tokenID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "babd0a3a-5372-405f-9464-16184cd27ea0", "token_id was not specified")
	}

	// list operation
	//---------------
	visitorCardCollectionModel, err := visitor.GetVisitorCardModelAndSetID(tokenID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorCardCollectionModel.Delete()

	return "ok", nil
}
