package cart

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cart", "GET", "info", restCartInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cart", "POST", "add/:productID/:qty", restCartAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cart", "PUT", "update/:itemIdx/:qty", restCartUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cart", "DELETE", "delete/:itemIdx", restCartDelete)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to get cart information
//   - parent categories and categorys will not be present in list
func restCartInfo(params *api.StructAPIHandlerParams) (interface{}, error) {

	currentCart, err := cart.GetCurrentCart(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var items []map[string]interface{}

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

			mediaPath, _ := product.GetMediaPath("image")

			productData["name"] = product.GetName()
			productData["sku"] = product.GetSku()
			productData["image"] = mediaPath + product.GetDefaultImage()
			productData["price"] = product.GetPrice()
			productData["weight"] = product.GetWeight()
			productData["options"] = product.GetOptions()

			item["product"] = productData
		}

		items = append(items, item)
	}

	result := map[string]interface{}{
		"visitor_id": currentCart.GetVisitorID(),
		"cart_info":  currentCart.GetCartInfo(),
		"items":      items,
	}

	return result, nil
}

// WEB REST API for adding new item into cart
//   - "pid" (product id) should be specified
//   - "qty" and "options" are optional params
func restCartAdd(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var pid string
	reqPid, present := params.RequestURLParams["productID"]
	pid = utils.InterfaceToString(reqPid)
	if !present || pid == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c21dac87-4f93-48dc-b997-bbbe558cfd29", "pid should be specified")
	}

	qty := 1
	reqQty, present := params.RequestURLParams["qty"]
	if present {
		qty = utils.InterfaceToInt(reqQty)
	}

	options := reqData
	reqOptions, present := reqData["options"]
	if present {
		if tmpOptions, ok := reqOptions.(map[string]interface{}); ok {
			options = tmpOptions
		}
	}

	// operation
	//----------
	currentCart, err := cart.GetCurrentCart(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	_, err = currentCart.AddItem(pid, qty, options)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentCart.Save()

	eventData := map[string]interface{}{"session": params.Session, "cart": currentCart, "pid": pid, "qty": qty, "options": options}
	env.Event("api.cart.addToCart", eventData)

	return "ok", nil
}

// WEB REST API used to update cart item qty
//   - "itemIdx" and "qty" should be specified in request URI
func restCartUpdate(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	if !utils.KeysInMapAndNotBlank(params.RequestURLParams, "itemIdx", "qty") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "16311f44-3f38-436d-82ca-8a9c08c47928", "itemIdx and qty should be specified")
	}

	itemIdx, err := utils.StringToInteger(params.RequestURLParams["itemIdx"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	qty, err := utils.StringToInteger(params.RequestURLParams["qty"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	if qty <= 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "701264ec-114b-4e18-971b-9965b70d534c", "qty should be greather then 0")
	}

	// operation
	//----------
	currentCart, err := cart.GetCurrentCart(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	found := false
	cartItems := currentCart.GetItems()

	for _, cartItem := range cartItems {
		if cartItem.GetIdx() == itemIdx {
			cartItem.SetQty(qty)
			found = true
			break
		}
	}

	if !found {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b1ae8e41-3aef-4f2e-b417-bd6975ff7bb1", "wrong itemIdx was specified")
	}

	currentCart.Save()

	return "ok", nil
}

// WEB REST API used to delete cart item from cart
//   - "itemIdx" should be specified in request URI
func restCartDelete(params *api.StructAPIHandlerParams) (interface{}, error) {

	_, present := params.RequestURLParams["itemIdx"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6afc9a4e-9fb4-4c31-b8c5-f46b514ef86e", "itemIdx should be specified")
	}

	itemIdx, err := utils.StringToInteger(params.RequestURLParams["itemIdx"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	currentCart, err := cart.GetCurrentCart(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = currentCart.RemoveItem(itemIdx)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	currentCart.Save()

	return "ok", nil
}
