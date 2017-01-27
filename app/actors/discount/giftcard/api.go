package giftcard

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
	"strings"
)

// setupAPI configures the API endpoints for the giftcard package
func setupAPI() error {

	service := api.GetRestService()

	// store
	service.GET("giftcards/:giftcode", GetSingleCode)
	service.GET("giftcards", GetList)

	// cart endpoints
	service.POST("cart/giftcards/:giftcode", Apply)
	service.DELETE("cart/giftcards/:giftcode", Remove)

	service.POST("giftcard/", createFromAdmin)

	return nil
}

// GetSingleCode returns the gift card and related info
//    - giftcode must be specified on the request
func GetSingleCode(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardID := context.GetRequestArgument("giftcode")
	if giftCardID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "06792fd7-c838-4acc-9c6f-cb8fcff833dd", "No giftcard code specified in the request.")
	}

	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.AddFilter("code", "=", giftCardID)
	rows, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dd7b2130-b5ed-4b26-b1fc-2d36c3bf147f", "No giftcard code matching the one supplied on the request found.")
	}

	return rows[0], nil
}

// GetList returns a list of gift cards for the visitor id in the context passed
//    - visitor must be logged in
func GetList(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "77d16dff-95bc-433d-9876-cc36e3645489", "Please log in to complete your request.")
	}

	if api.ValidateAdminRights(context) != nil {
		collection.AddFilter("visitor_id", "=", visitorID)
	}

	dbRecords, err := collection.Load()

	return dbRecords, env.ErrorDispatch(err)
}

// Apply applies the provided gift card to current checkout
//   - Gift Card code should be specified in "giftcode" argument
func Apply(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardCode := context.GetRequestArgument("giftcode")

	// getting applied gift codes array for current session
	appliedGiftCardCodes := utils.InterfaceToStringArray(context.GetSession().Get(ConstSessionKeyAppliedGiftCardCodes))

	// checking if codes have previously been applied
	if utils.IsInArray(giftCardCode, appliedGiftCardCodes) {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1c310f79-0f79-493a-b761-ad4f24542559", "This code, "+giftCardCode+" has already been applied.")
	}

	// loading gift codes for specified code
	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}
	err = collection.AddFilter("code", "=", giftCardCode)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	records, err := collection.Load()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// checking and applying provided gift card codes
	if len(records) == 1 && utils.InterfaceToString(records[0]["code"]) == giftCardCode {
		if utils.InterfaceToFloat64(records[0]["amount"]) <= 0 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "ce349f59-51c7-43ec-a64c-80f7d4af6d3c", "The provided giftcard value has been exhausted.")
		}

		// giftcard code is valid - applying it
		appliedGiftCardCodes = append(appliedGiftCardCodes, giftCardCode)
		context.GetSession().Set(ConstSessionKeyAppliedGiftCardCodes, appliedGiftCardCodes)

	} else {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2b55d714-2cba-49f8-ad7d-fdc542bfc2a3", "The provided giftcard code cannot be found, "+giftCardCode+".")
	}

	return "ok", nil
}

// Remove removes the application of the gift card value from the
// current checkout
//   - giftcard code should be specified in the "giftcode" argument
//   - use "*" as giftcard code to 'remove' all giftcard discounts
func Remove(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardID := context.GetRequestArgument("giftcode")
	if giftCardID == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e2bad33a-36e7-41d4-aea7-8fe1b97eb31c", "No giftcard code found on the request.")
	}

	if giftCardID == "*" {
		context.GetSession().Set(ConstSessionKeyAppliedGiftCardCodes, make([]string, 0))
		return "Remove successful", nil
	}

	appliedCoupons := utils.InterfaceToStringArray(context.GetSession().Get(ConstSessionKeyAppliedGiftCardCodes))
	if len(appliedCoupons) > 0 {
		var newAppliedCoupons []string
		for _, value := range appliedCoupons {
			if value != giftCardID {
				newAppliedCoupons = append(newAppliedCoupons, value)
			}
		}
		context.GetSession().Set(ConstSessionKeyAppliedGiftCardCodes, newAppliedCoupons)
	}

	return "Remove successful", nil
}


func createFromAdmin(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		env.ErrorDispatch(err)
		return false
	}

	if !utils.KeysInMapAndNotBlank(requestData, "amount", "message", "name", "recipient_mailbox", "sku") {
		context.SetResponseStatusBadRequest()
		// todo
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6372b9a3-29f3-4ea4-a19f-40051a8f330b", "email or group_ids have not been specified")
	}

	currentTime := time.Now()
	deliveryDate := utils.InterfaceToTime(requestData["delivery_date"])
	giftCardAmount := utils.InterfaceToFloat64(requestData["amount"])
	customMessage := utils.InterfaceToString(requestData["message"])
	recipientName := utils.InterfaceToString(requestData["name"])
	recipientEmail := utils.InterfaceToString(requestData["recipient_mailbox"])
	giftCardSku := utils.InterfaceToString(requestData["sku"])

	giftCardUniqueCode := utils.InterfaceToString(requestData["code"])
	if giftCardUniqueCode == "" {
		// generate unique code by unix nano time
		giftCardUniqueCode = utils.InterfaceToString(time.Now().UnixNano())
	}

	// todo check if code Unique

	// collect necessary info to variables
	// get a customer and his mail to set him as addressee
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		context.SetResponseStatusBadRequest()
		// todo
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "77d16dff-95bc-433d-9876-cc36e3645489", "Please log in to complete your request.")
	}

	giftCardSkuElement := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathGiftCardSKU))

	if strings.Contains(giftCardSku, giftCardSkuElement) == false {
		// todo
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "77d16dff-95bc-433d-9876-cc36e3645489", "Please log in to complete your request.")
	}


	giftCard := make(map[string]interface{})

	giftCard["code"] = giftCardUniqueCode
	giftCard["sku"] = giftCardSku

	giftCard["amount"] = giftCardAmount

	giftCard["visitor_id"] = visitorID

	giftCard["status"] = ConstGiftCardStatusNew
	giftCard["orders_used"] = make(map[string]float64)

	giftCard["name"] = recipientName
	giftCard["message"] = customMessage

	giftCard["recipient_mailbox"] = recipientEmail
	giftCard["delivery_date"] = deliveryDate

	giftCardID, err := giftCardCollection.Save(giftCard)
	if err != nil {
		env.ErrorDispatch(err)
		return false
	}

	//if deliveryDate.Truncate(time.Hour).Before(currentTime) {
	//	giftCardsToSendImmediately = append(giftCardsToSendImmediately, giftCardID)
	//}
	//
	//// run SendTask task to send immediately if delivery_date is today's date
	//if len(giftCardsToSendImmediately) > 0 {
	//	params := map[string]interface{}{
	//		"giftCards":          giftCardsToSendImmediately,
	//		"ignoreDeliveryDate": true,
	//	}
	//
	//	go SendTask(params)
	//}

	return true
}

