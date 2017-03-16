package giftcard

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"math"
	"time"
)

// setupAPI configures the API endpoints for the giftcard package
func setupAPI() error {

	service := api.GetRestService()

	// store
	service.GET("giftcards/:giftcode", GetSingleCode)
	service.GET("giftcards", GetList)
	service.GET("check/giftcards/:giftcode", IfGiftCardCodeUnique)
	service.GET("generate/giftcards/code", GetUniqueGiftCode)

	// cart endpoints
	service.POST("cart/giftcards/:giftcode", Apply)
	service.DELETE("cart/giftcards/:giftcode", Remove)

	// Admin Only
	service.GET("giftcard/:id", GetSingleID)
	service.POST("giftcard", api.IsAdminHandler(Create))
	service.PUT("giftcard/:id", api.IsAdminHandler(Edit))
	service.GET("giftcard/:id/history", api.IsAdminHandler(GetHistory))

	return nil
}

// GetSingleCode returns the gift card and related info
//    - giftcode must be specified on the request
func GetSingleCode(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardID := context.GetRequestArgument("giftcode")
	if giftCardID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "06792fd7-c838-4acc-9c6f-cb8fcff833dd", "No giftcard code specified in the request.")
	}

	giftCard, err := getGiftCardByCode(giftCardID)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(context) != nil && giftCard["status"] == ConstGiftCardStatusCancelled {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "34deea74-1378-4ec8-b7c1-53d73d8a8987", "Giftcard has cancelled.")
	}

	return giftCard, nil
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

	if api.ValidateAdminRights(context) != nil {
		visitorID := visitor.GetCurrentVisitorID(context)
		if visitorID == "" {
			context.SetResponseStatusBadRequest()
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "77d16dff-95bc-433d-9876-cc36e3645489", "Please log in to complete your request.")
		}

		collection.AddFilter("visitor_id", "=", visitorID)
	}

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return collection.Count()
	}

	dbRecords, err := collection.Load()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	for _, value := range dbRecords {

		initialAmount := utils.InterfaceToFloat64(value["amount"])
		for _, amount := range utils.InterfaceToMap(value["orders_used"]) {
			initialAmount = initialAmount + math.Abs(utils.InterfaceToFloat64(amount))
		}

		value["initial_amount"] = initialAmount
	}

	return dbRecords, nil
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
	record, err := getGiftCardByCode(giftCardCode)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	if record["status"] == ConstGiftCardStatusCancelled {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "132a21e7-67bf-42f9-a21d-bdb5b0ca1cf2", "Giftcard has cancelled.")
	}

	// checking and applying provided gift card codes
	if utils.InterfaceToString(record["code"]) == giftCardCode {
		if utils.InterfaceToFloat64(record["amount"]) <= 0 {
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

// GetHistory returns a history of gift cards for the admin in the context passed
//    - giftcard id should be specified in the "giftid" argument
func GetHistory(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardID := context.GetRequestArgument("id")
	if giftCardID == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "10ab8fd5-05ca-43e2-9da9-8acac0ea13f9", "No giftcard code specified in the request.")
	}

	giftCard, err := GetGiftCardByID(giftCardID)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	var historyData []map[string]interface{}

	for orderId, amount := range utils.InterfaceToMap(giftCard["orders_used"]) {
		orderData, err := order.LoadOrderByID(orderId)
		if err != nil {
			context.SetResponseStatusBadRequest()
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cb86e6de-b94d-4480-bc87-90301676f4fe", "system error loading id from db: "+utils.InterfaceToString(orderId))
		}
		historyData = append(historyData, map[string]interface{}{
			"order_id":         utils.InterfaceToString(orderId),
			"amount":           math.Abs(utils.InterfaceToFloat64(amount)),
			"transaction_date": orderData.Get("created_at"),
		})
	}

	return historyData, nil
}

// Create gift card from admin panel
func Create(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "amount", "message", "name", "recipient_mailbox", "sku", "code") {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e4a6ad26-fd34-428b-8cca-9baed590a67e", "amount or message or name or recipient_mailbox or sku or code have not been specified")
	}

	currentTime := time.Now()
	deliveryDate := utils.InterfaceToTime(requestData["delivery_date"])
	giftCardAmount := utils.InterfaceToInt(requestData["amount"])
	customMessage := utils.InterfaceToString(requestData["message"])
	recipientName := utils.InterfaceToString(requestData["name"])
	recipientEmail := utils.InterfaceToString(requestData["recipient_mailbox"])
	giftCardSku := utils.InterfaceToString(requestData["sku"])

	giftCardUniqueCode := utils.InterfaceToString(requestData["code"])

	_, err = getGiftCardByCode(giftCardUniqueCode)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	// collect necessary info to variables
	// get a customer and his mail to set him as addressee
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" && api.ValidateAdminRights(context) != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "39a37b12-93fb-4660-836e-ef5e07c2af52", "Please log in to complete your request.")
	}

	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return false, env.ErrorDispatch(err)
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
	giftCard["created_at"] = currentTime

	giftCardID, err := giftCardCollection.Save(giftCard)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return false, env.ErrorDispatch(err)
	}
	giftCard["_id"] = giftCardID

	var giftCardsToSendImmediately []string

	// run SendTask task to send immediately if delivery_date is today's date
	if deliveryDate.Truncate(time.Hour).Before(currentTime) {
		giftCardsToSendImmediately = append(giftCardsToSendImmediately, giftCardID)

		params := map[string]interface{}{
			"giftCards":          giftCardsToSendImmediately,
			"ignoreDeliveryDate": true,
		}

		go SendTask(params)
	}

	return giftCard, nil
}

// getGiftCardByCode returns a list of gift cards for the giftCardCode
func getGiftCardByCode(giftCardCode string) (map[string]interface{}, error) {

	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.AddFilter("code", "=", giftCardCode)
	rows, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return make(map[string]interface{}), nil
	}

	return rows[0], nil
}

// IfGiftCardCodeUnique returns a history of gift cards for the admin in the context passed
//    - giftcard code should be specified in the "gift card code" argument
func IfGiftCardCodeUnique(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardUniqueCode := context.GetRequestArgument("giftcode")
	if giftCardUniqueCode == "" {
		context.SetResponseStatusBadRequest()
		return false, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e2940eda-4023-4a27-80d3-c39bab1c28fe", "giftcode have not been specified")
	}

	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.AddFilter("code", "=", giftCardUniqueCode)
	rows, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if len(rows) != 0 {
		return false, nil
	}

	return true, nil
}

// GetUniqueGiftCode returns unique gift card code
func GetUniqueGiftCode(context api.InterfaceApplicationContext) (interface{}, error) {
	return utils.InterfaceToString(time.Now().UnixNano()), nil
}

// Edit gift card
func Edit(context api.InterfaceApplicationContext) (interface{}, error) {

	giftID := context.GetRequestArgument("id")
	if giftID == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "ac41eca7-4590-439d-92a2-4eaf2d9a45f2", "ID have not been specified")
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "amount", "message", "recipient_mailbox", "code", "sku", "delivery_date", "recipient_mailbox", "name") {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "237a4b72-2373-4e65-a546-4194a35e3d82", "amount or message or recipient_mailbox or code or sku or delivery_date or recipient_mailbox or recipient_name have not been specified")
	}

	giftCard, err := GetGiftCardByID(giftID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if giftCard["status"] == ConstGiftCardStatusCancelled {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "39195589-ca3b-43ac-b8b2-d879fa80057d", "Giftcard has cancelled.")
	}

	giftCardAmount := utils.InterfaceToInt(requestData["amount"])
	customMessage := utils.InterfaceToString(requestData["message"])
	recipientEmail := utils.InterfaceToString(requestData["recipient_mailbox"])
	giftCardUniqueCode := utils.InterfaceToString(requestData["code"])
	sku := utils.InterfaceToString(requestData["sku"])
	delivery_date := utils.InterfaceToTime(requestData["delivery_date"])
	recipient_mailbox := utils.InterfaceToString(requestData["recipient_mailbox"])
	name := utils.InterfaceToString(requestData["name"])

	row, err := getGiftCardByCode(giftCardUniqueCode)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	if row["_id"] != "" && row["_id"] != giftID {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0f26ec81-e555-4856-89ed-e2b4e050c808", "Gift code must be unique")
	}

	giftCardCollection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if utils.InterfaceToString(requestData["status"]) == ConstGiftCardStatusCancelled {
		giftCard["status"] = ConstGiftCardStatusCancelled
	}

	giftCard["code"] = giftCardUniqueCode
	giftCard["amount"] = giftCardAmount
	giftCard["message"] = customMessage
	giftCard["recipient_mailbox"] = recipientEmail
	giftCard["sku"] = sku
	giftCard["delivery_date"] = delivery_date
	giftCard["recipient_mailbox"] = recipient_mailbox
	giftCard["name"] = name

	_, err = giftCardCollection.Save(giftCard)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	return giftCard, nil

}

// GetSingleCode returns the gift card and related info
//    - id must be specified on the request
func GetSingleID(context api.InterfaceApplicationContext) (interface{}, error) {

	giftCardID := context.GetRequestArgument("id")
	if giftCardID == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cc227376-4654-4036-b01f-6706a3ed55c1", "No giftcard id specified in the request.")
	}

	//giftCard, err := GetGiftCardByID(giftCardID)
	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(context) != nil {
		visitorID := visitor.GetCurrentVisitorID(context)
		if visitorID == "" {
			context.SetResponseStatusBadRequest()
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "77d16dff-95bc-433d-9876-cc36e3645489", "Please log in to complete your request.")
		}

		collection.AddFilter("visitor_id", "=", visitorID)
	}

	collection.AddFilter("_id", "=", giftCardID)
	rows, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "96afa35f-24d9-4de0-9521-22b5aece5f57", "No giftcard code matching the one supplied on the request found.")
	}

	giftCard := rows[0]

	if api.ValidateAdminRights(context) != nil && giftCard["status"] == ConstGiftCardStatusCancelled {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "58d83871-bc9e-4e98-a0f3-d57caa59f039", "Giftcard has cancelled.")
	}

	return giftCard, nil
}

func GetGiftCardByID(giftCardID string) (map[string]interface{}, error) {

	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.AddFilter("_id", "=", giftCardID)
	rows, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if len(rows) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8bd579a3-bacc-4bcf-acd6-2f8f1393745a", "No giftcard code matching the one supplied on the request found.")
	}

	return rows[0], nil
}
