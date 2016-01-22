package subscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	// Administrative
	err = api.GetRestService().RegisterAPI("subscriptions", api.ConstRESTOperationGet, AdminList)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("subscriptions/:id", api.ConstRESTOperationGet, AdminOne)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("subscriptions/:id", api.ConstRESTOperationUpdate, Update)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Public
	err = api.GetRestService().RegisterAPI("visit/subscriptions", api.ConstRESTOperationGet, List)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("visit/subscriptions/:id", api.ConstRESTOperationUpdate, Update)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Other thing
	err = api.GetRestService().RegisterAPI("subscriptional/checkout", api.ConstRESTOperationGet, APICheckCheckoutSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// AdminList returns a list of subscriptions for visitor
//   - if "action" parameter is set to "count" result value will be just a number of list items
func AdminList(context api.InterfaceApplicationContext) (interface{}, error) {

	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// list operation
	subscriptionCollectionModel, err := subscription.GetSubscriptionCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbCollection := subscriptionCollectionModel.GetDBCollection()

	// filters handle
	models.ApplyFilters(context, dbCollection)

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return dbCollection.Count()
	}

	// limit parameter handle
	subscriptionCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, subscriptionCollectionModel)

	return subscriptionCollectionModel.List()
}

// List returns a list of subscriptions for visitor
//   - if "action" parameter is set to "count" result value will be just a number of list items
func List(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c73e39c9-dc23-463b-9792-a5d3f7e4d9dd", "You should log in first")
	}

	//subscriptionCollection =
	// list operation
	subscriptionCollectionModel, err := subscription.GetSubscriptionCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbCollection := subscriptionCollectionModel.GetDBCollection()
	dbCollection.AddStaticFilter("visitor_id", "=", visitorID)
	dbCollection.AddStaticFilter("status", "=", ConstSubscriptionStatusConfirmed)
	models.ApplyFilters(context, dbCollection)

	// checking for a "count" request
	if context.GetRequestArgument("count") != "" {
		return dbCollection.Count()
	}

	// limit parameter handle
	dbCollection.SetLimit(models.GetListLimit(context))

	dbCollection.SetResultColumns("_id", "period", "action_date", "cart_id")
	// extra parameter handle
	extra := context.GetRequestArgument("extra")
	extraAttributes := utils.Explode(extra, ",")
	for _, attributeName := range extraAttributes {
		dbCollection.SetResultColumns(attributeName)
	}

	records, err := dbCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, record := range records {
		storedCart, err := cart.LoadCartByID(utils.InterfaceToString(record["cart_id"]))
		if err != nil {
			continue
		}

		var items []interface{}

		for _, cartItem := range storedCart.GetItems() {

			item := make(map[string]interface{})

			item["_id"] = cartItem.GetID()
			item["idx"] = cartItem.GetIdx()
			item["qty"] = cartItem.GetQty()
			item["pid"] = cartItem.GetProductID()
			item["options"] = cartItem.GetOptions()
			items = append(items, item)
		}

		record["items"] = items
	}

	return records, nil
}

// AdminOne return specified subscription information
//   - subscription id should be specified in "id" argument
func AdminOne(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	subscriptionID := context.GetRequestArgument("id")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b626ec0a-a317-4b63-bd05-cc23932bdfe0", "subscription id should be specified")
	}

	subscriptionModel, err := subscription.LoadSubscriptionByID(subscriptionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := subscriptionModel.ToHashMap()

	// Attach the order items
	orderModel, err := order.LoadOrderByID(subscriptionModel.GetOrderID())
	result["order_items"] = orderModel.GetItems()

	result["payment_method_name"] = subscriptionModel.GetPaymentMethod().GetName()
	result["shipping_method_name"] = subscriptionModel.GetShippingMethod().GetName()

	return result, nil
}

// APICheckCheckoutSubscription provide check is current checkout allows to create new subscription
func APICheckCheckoutSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check visitor to be registered
	visitorID := visitor.GetCurrentVisitorID(context)
	if api.ValidateAdminRights(context) != nil && visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e6109c04-e35a-4a90-9593-4cc1f141a358", "you are not logged in")
	}

	currentCheckout, err := checkout.GetCurrentCheckout(context, false)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if err := validateCheckoutToSubscribe(currentCheckout); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

func Update(context api.InterfaceApplicationContext) (interface{}, error) {

	// validate params
	subscriptionID := context.GetRequestArgument("id")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4e8f9873-9144-42ae-b119-d1e95bb1bbfd", "subscription id should be specified")
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestedStatus := utils.InterfaceToString(requestData["status"])
	if requestedStatus == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "71fc926c-d2a0-4c8a-9462-b5274346ed23", "status should be specified")
	}

	subscriptionInstance, err := subscription.LoadSubscriptionByID(subscriptionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// validate ownership
	isAdmin := api.ValidateAdminRights(context) == nil
	isOwner := subscriptionInstance.GetVisitorID() == visitor.GetCurrentVisitorID(context)

	if !isAdmin && !isOwner {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "bae87bfa-0fa2-4256-ab11-2fffa20bfa00", "Subscription ownership could not be verified")
	}

	err = subscriptionInstance.SetStatus(requestedStatus)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", subscriptionInstance.Save()
}
