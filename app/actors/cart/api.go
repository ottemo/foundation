package cart

import (
	"fmt"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"
	"time"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("cart", APICartInfo)
	service.POST("cart/item", APICartItemAdd)
	service.PUT("cart/item/:itemIdx/:qty", APICartItemUpdate)
	service.DELETE("cart/item/:itemIdx", APICartItemDelete)

	service.GET("cart-abandon", testHandler)

	return nil
}

func sendAbandonEmail(context api.InterfaceApplicationContext) (interface{}, error) {
	// cartID := context.GetRequestArgument("cartID")
	// aCart, err := cart.LoadCartByID(cartID)
	// if err != nil {
	// 	return "issue loading that cart", env.ErrorDispatch(err)
	// }

	// err = sendAbandonEmail(emailData)

	return "ok", nil
}

func testHandler(context api.InterfaceApplicationContext) (interface{}, error) {
	fmt.Println("request")

	// Check frequency
	beforeDate, isEnabled := getConfigSendBefore()
	if !isEnabled {
		context.SetResponseStatusNotFound()
		return "endpoint not enabled", nil
	}

	resultCarts := getAbandonedCarts(beforeDate)
	actionableCarts := getActionableCarts(resultCarts)

	for _, aCart := range actionableCarts {
		err := sendAbandonEmail(aCart)
		if err != nil {
			continue
		}

		flagCartAsEmailed(aCart.Cart.ID)
	}

	return actionableCarts, nil
}

type AbandonCartEmailData struct {
	Visitor AbandonVisitor
	Cart    AbandonCart
}

type AbandonVisitor struct {
	Email     string
	FirstName string
	LastName  string
}

type AbandonCart struct {
	ID string
	// Items []AbandonCartItem
}

// type AbandonCartItem struct {
// 	Name  string
// 	SKU   string
// 	Price float64
// 	Image string
// }

func getConfigSendBefore() (time.Time, bool) {
	var isEnabled bool
	beforeConfig := utils.InterfaceToInt(env.ConfigGetValue(ConstConfigPathCartAbandonEmailSendTime))

	// Flag it as enabled
	if beforeConfig != 0 {
		isEnabled = true
	}

	// Build out the time to send before, we are expecting a config
	// that is a negative int
	beforeDuration := time.Duration(beforeConfig) * time.Hour
	beforeDate := time.Now().Add(beforeDuration)

	return beforeDate, isEnabled
}

// Get the abandoned carts
// - active
// - were updated in our time frame
// - have not been sent an abandon cart email
func getAbandonedCarts(beforeDate time.Time) []map[string]interface{} {
	dbEngine := db.GetDBEngine()
	cartCollection, _ := dbEngine.GetCollection(ConstCartCollectionName)
	cartCollection.AddFilter("active", "=", true)
	cartCollection.AddFilter("custom_info.is_abandon_email_sent", "!=", true)
	cartCollection.AddFilter("updated_at", "<", beforeDate)
	cartCollection.AddSort("updated_at", true)
	cartCollection.SetLimit(0, 3) //TODO: REMOVE
	resultCarts, _ := cartCollection.Load()

	fmt.Println("abandoned carts found", len(resultCarts)) //TODO: CLEANUP
	return resultCarts
}

func getActionableCarts(resultCarts []map[string]interface{}) []AbandonCartEmailData {
	allCartEmailData := []AbandonCartEmailData{}

	// Determine which carts have an email we can use
	for _, resultCart := range resultCarts {
		var email, firstName, lastName string
		cartID := utils.InterfaceToString(resultCart["_id"])
		sessionID := utils.InterfaceToString(resultCart["session_id"])
		visitorID := utils.InterfaceToString(resultCart["visitor_id"])

		//TODO: CLEANUP
		fmt.Println("cartid:", cartID)
		fmt.Println("visit :", visitorID)
		fmt.Println("sesh  :", sessionID)

		// try to get by visitor_id
		if visitorID != "" {
			vModel, _ := visitor.LoadVisitorByID(visitorID)
			email = vModel.GetEmail()
			firstName = vModel.GetFirstName()
			lastName = vModel.GetLastName()

		} else if sessionID != "" {
			create := false
			sessionWrapper, _ := api.GetSessionService().Get(sessionID, create)
			sCheckout := utils.InterfaceToMap(sessionWrapper.Get(checkout.ConstSessionKeyCurrentCheckout))

			scInfo := utils.InterfaceToMap(sCheckout["Info"])
			email = utils.InterfaceToString(scInfo["customer_email"])
			//NOTE: We have customer_name here as well, which we could split
			//      or we could look to see if the address is filled out yet

			// fmt.Println("info map:", scInfo) //TODO: CLEANUP
		}

		// TODO: if we don't have an email then flag this cart as don't update?

		// no email address for us to contact, move along
		if email == "" {
			continue
		}

		// Assemble the details needed for further actions
		cartEmailData := AbandonCartEmailData{
			Visitor: AbandonVisitor{
				Email:     email,
				FirstName: firstName,
				LastName:  lastName,
			},
			Cart: AbandonCart{
				ID: cartID,
			},
		}

		// NOTE: In v1 we aren't including cart item details
		// Get the cart items for the carts we are about to email
		// cartItemsCollection, err := dbEngine.GetCollection(ConstCartItemsCollectionName)
		// cartItemsCollection.AddFilter("cart_id", "=", it.GetID())
		// cartItems, err := cartItemsCollection.Load()

		allCartEmailData = append(allCartEmailData, cartEmailData)
	}

	return allCartEmailData
}

func sendAbandonEmail(emailData AbandonCartEmailData) error {
	subject := "Hi there" //TODO:
	template := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathCartAbandonEmailTemplate))
	if template == "" {
		return env.ErrorDispatch(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1756ec63-7cd7-4764-a8ff-64b142fc3f9f", "Abandon cart emails want to send but the template is empty"))
	}

	templateData := utils.InterfaceToMap(emailData)
	templateData["Site"] = map[string]interface{}{
		"Url": app.GetStorefrontURL(""),
	}

	body, err := utils.TextTemplate(template, templateData)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = app.SendMail(emailData.Visitor.Email, subject, body)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func flagCartAsEmailed(cartID string) {
	iCart, _ := cart.LoadCartByID(cartID)

	info := iCart.GetCustomInfo()
	info["is_abandon_email_sent"] = true
	info["abandon_email_sent_at"] = time.Now()
	iCart.SetCustomInfo(info)

	iCart.Save()
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
