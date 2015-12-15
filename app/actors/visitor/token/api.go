package token

import (
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
	var err error

	err = api.GetRestService().RegisterAPI("visit/token", api.ConstRESTOperationCreate, APICreateToken)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("visit/tokens", api.ConstRESTOperationGet, APIListVisitorCards)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	//	err = api.GetRestService().RegisterAPI("token/:tokenID", api.ConstRESTOperationUpdate, APIGetToken)
	//	if err != nil {
	//		return env.ErrorDispatch(err)
	//	}

	return nil
}

// APICreateToken creates a request body for posting credit card info to payment system with 0 amount payment
// for obtaining token on this card and saving it for visitor
func APICreateToken(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return "you are not logined in", nil
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, err
	}

	paymentMethodCode := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "payment", "payment_method"))
	if paymentMethodCode == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6d1691c8-2d26-44be-b90d-24d920e26301", "payment method not selected")
	}

	value, present := requestData["cc"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2e9f1bfc-ec9f-4017-83c6-4d04b95b9c08", "payment info not specified")
	}

	creditCardInfo := utils.InterfaceToMap(value)

	var paymentMethod checkout.InterfacePaymentMethod

	for _, payment := range checkout.GetRegisteredPaymentMethods() {
		if payment.GetCode() == paymentMethodCode {
			if !payment.IsTokenable(nil) {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "519ef43c-4d07-4b64-90f7-7fdc3657940a", "for selected payment method credit card can't be saved")
			}
			paymentMethod = payment
		}
	}

	if paymentMethod == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c80c4106-1208-4d0b-8577-0889f608869b", "such payment method not existing")
	}

	paymentInfo := map[string]interface{}{
		checkout.ConstPaymentActionTypeKey: checkout.ConstPaymentActionTypeCreateToken,
		"cc": creditCardInfo,
	}

	// contains creditCardLastFour, creditCardType, responseMessage, responseResult, transactionID, creditCardExp
	paymentResult, err := paymentMethod.Authorize(nil, paymentInfo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cardInfoMap := utils.InterfaceToMap(paymentResult)
	if !utils.KeysInMapAndNotBlank(cardInfoMap, "transactionID", "creditCardLastFour") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "22e17290-56f3-452a-8d54-18d5a9eb2833", "transaction can't be obtained")
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

	return tokenRecord, nil
}

// APIListVisitorCards return a list of existing tokens for visitor
func APIListVisitorCards(context api.InterfaceApplicationContext) (interface{}, error) {

	// if visitorID was specified - using this otherwise, taking current visitor
	visitorID := context.GetRequestArgument("visitorID")
	if visitorID == "" {

		sessionVisitorID := visitor.GetCurrentVisitorID(context)
		if sessionVisitorID == "" {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2ac4c16b-9241-406e-b35a-399813bb6ca5", "you are not logined in")
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
