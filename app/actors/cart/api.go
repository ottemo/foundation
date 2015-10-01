package cart

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cart", api.ConstRESTOperationGet, APICartInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cart/item", api.ConstRESTOperationCreate, APICartItemAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cart/item/:itemIdx/:qty", api.ConstRESTOperationUpdate, APICartItemUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cart/item/:itemIdx", api.ConstRESTOperationDelete, APICartItemDelete)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APICartInfo returns get cart related information
func APICartInfo(context api.InterfaceApplicationContext) (interface{}, error) {

	currentCart, err := cart.GetCurrentCart(context, false)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var items []map[string]interface{}
	result := map[string]interface{}{
		"visitor_id": "",
		"cart_info":  nil,
		"items":      items,
	}

	if currentCart != nil {

		mediaStorage, err := media.GetMediaStorage()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		cartItems := currentCart.GetItems()
		for _, cartItem := range cartItems {

			item := make(map[string]interface{})

			item["_id"] = cartItem.GetID()
			item["idx"] = cartItem.GetIdx()
			item["qty"] = cartItem.GetQty()
			item["pid"] = cartItem.GetProductID()
			item["options"] = cartItem.GetOptions()

			if product := cartItem.GetProduct(); product != nil {

				product.ApplyOptions(cartItem.GetOptions())

				productData := make(map[string]interface{})

				productData["name"] = product.GetName()
				productData["sku"] = product.GetSku()
				productData["price"] = product.GetPrice()
				productData["weight"] = product.GetWeight()
				productData["options"] = product.GetOptions()

				productData["image"], err = mediaStorage.GetSizes(product.GetModelName(), product.GetID(), "image", product.GetDefaultImage())
				if err != nil {
					env.LogError(err)
				}

				item["product"] = productData
			}

			items = append(items, item)
		}

		result["visitor_id"] = currentCart.GetVisitorID()
		result["cart_info"] = currentCart.GetCartInfo()
		result["items"] = items
	}

	return result, nil
}

// APICartItemAdd adds specified product to cart
//   - "productID" and "qty" should be specified as arguments
func APICartItemAdd(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	pid := utils.InterfaceToString(api.GetArgumentOrContentValue(context, "pid"))
	if pid == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c21dac87-4f93-48dc-b997-bbbe558cfd29", "pid should be specified")
	}

	qty := 1
	requestedQty := api.GetArgumentOrContentValue(context, "qty")
	if requestedQty != "" {
		qty = utils.InterfaceToInt(requestedQty)
	}

	// we are considering json content as product options unless it have specified options key
	options, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestedOptions := api.GetArgumentOrContentValue(context, "options")
	if requestedOptions != nil {
		if reqestedOptionsAsMap, ok := requestedOptions.(map[string]interface{}); ok {
			options = reqestedOptionsAsMap
		} else {
			options = utils.InterfaceToMap(requestedOptions)
		}
	}

	// operation
	//----------
	currentCart, err := cart.GetCurrentCart(context, true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	_, err = currentCart.AddItem(pid, qty, options)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentCart.Save()

	eventData := map[string]interface{}{"session": context.GetSession(), "cart": currentCart, "pid": pid, "qty": qty, "options": options}
	env.Event(ConstEventAPIAdd, eventData)

	eventData = map[string]interface{}{"session": context.GetSession(), "cart": currentCart, "idx": nil, "pid": pid, "qty": qty, "options": options}
	env.Event(ConstEventAPIUpdate, eventData)

	return "ok", nil
}

// APICartItemUpdate changes qty and/or option for cart item
//   - "itemIdx" and "qty" should be specified as arguments
func APICartItemUpdate(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	if !utils.KeysInMapAndNotBlank(context.GetRequestArguments(), "itemIdx", "qty") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "16311f44-3f38-436d-82ca-8a9c08c47928", "itemIdx and qty should be specified")
	}

	itemIdx, err := utils.StringToInteger(context.GetRequestArgument("itemIdx"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	qty, err := utils.StringToInteger(context.GetRequestArgument("qty"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	if qty <= 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "701264ec-114b-4e18-971b-9965b70d534c", "qty should be greather then 0")
	}

	// operation
	//----------
	currentCart, err := cart.GetCurrentCart(context, true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	found := false
	cartItems := currentCart.GetItems()

	eventData := map[string]interface{}{"session": context.GetSession(), "cart": currentCart, "idx": itemIdx, "qty": qty}

	for _, cartItem := range cartItems {
		if cartItem.GetIdx() == itemIdx {
			cartItem.SetQty(qty)
			found = true

			eventData["pid"] = cartItem.GetProductID()
			eventData["options"] = cartItem.GetOptions()

			break
		}
	}

	if !found {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b1ae8e41-3aef-4f2e-b417-bd6975ff7bb1", "wrong itemIdx was specified")
	}

	currentCart.Save()

	env.Event(ConstEventAPIUpdate, eventData)

	return "ok", nil
}

// APICartItemDelete removes specified item from cart item from cart
//   - "itemIdx" should be specified as argument (item index can be obtained from APICartInfo)
func APICartItemDelete(context api.InterfaceApplicationContext) (interface{}, error) {

	reqItemIdx := context.GetRequestArgument("itemIdx")
	if reqItemIdx == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6afc9a4e-9fb4-4c31-b8c5-f46b514ef86e", "itemIdx should be specified")
	}

	itemIdx, err := utils.StringToInteger(reqItemIdx)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	currentCart, err := cart.GetCurrentCart(context, true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = currentCart.RemoveItem(itemIdx)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	currentCart.Save()

	eventData := map[string]interface{}{"session": context.GetSession(), "cart": currentCart, "idx": itemIdx, "qty": 0}
	env.Event(ConstEventAPIUpdate, eventData)

	return "ok", nil
}
