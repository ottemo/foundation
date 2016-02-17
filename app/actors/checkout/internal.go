package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/actors/discount/coupon"
	"github.com/ottemo/foundation/app/actors/discount/giftcard"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"

	"strings"
)

// SendOrderConfirmationMail sends an order confirmation email
func (it *DefaultCheckout) SendOrderConfirmationMail() error {

	checkoutOrder := it.GetOrder()
	if checkoutOrder == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e7c69056-cc28-4632-9524-50d71b909d83", "given checkout order does not exists")
	}

	confirmationEmail := utils.InterfaceToString(env.ConfigGetValue(checkout.ConstConfigPathConfirmationEmail))
	if confirmationEmail != "" {
		timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
		giftCardSku := utils.InterfaceToString(env.ConfigGetValue(giftcard.ConstConfigPathGiftCardSKU))

		email := utils.InterfaceToString(checkoutOrder.Get("customer_email"))
		if email == "" {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1202fcfb-da3f-4a0f-9a2e-92f288fd3881", "customer email for order is not set")
		}

		visitorMap := make(map[string]interface{})
		if visitorModel := it.GetVisitor(); visitorModel != nil {
			visitorMap = visitorModel.ToHashMap()
		} else {
			visitorMap["first_name"] = checkoutOrder.Get("customer_name")
			visitorMap["email"] = checkoutOrder.Get("customer_email")
		}

		orderMap := checkoutOrder.ToHashMap()
		var orderItems []map[string]interface{}

		for _, item := range checkoutOrder.GetItems() {
			options := make(map[string]interface{})

			for optionName, optionKeys := range item.GetOptions() {
				optionMap := utils.InterfaceToMap(optionKeys)
				options[optionName] = optionMap["value"]

				// Giftcard's delivery date
				if strings.Contains(item.GetSku(), giftCardSku) {
					if utils.IsAmongStr(optionName, "Date", "Delivery Date", "send_date", "Send Date", "date") {
						// Localize and format the date
						giftcardDeliveryDate, _ := utils.MakeTZTime(utils.InterfaceToTime(optionMap["value"]), timeZone)
						if !utils.IsZeroTime(giftcardDeliveryDate) {
							//TODO: Should be "Monday Jan 2 15:04 (MST)" but we have a bug
							options[optionName] = giftcardDeliveryDate.Format("Monday Jan 2 15:04")
						}
					}
				}
			}

			orderItems = append(orderItems, map[string]interface{}{
				"name":    item.GetName(),
				"options": options,
				"sku":     item.GetSku(),
				"qty":     item.GetQty(),
				"price":   item.GetPrice()})
		}

		// convert date of order creation to store time zone
		if date, present := orderMap["created_at"]; present {
			convertedDate, _ := utils.MakeTZTime(utils.InterfaceToTime(date), timeZone)
			if !utils.IsZeroTime(convertedDate) {
				orderMap["created_at"] = convertedDate
			}
		}

		orderMap["items"] = orderItems
		orderMap["payment_method_title"] = it.GetPaymentMethod().GetName()
		orderMap["shipping_method_title"] = it.GetShippingMethod().GetName()

		customInfo := make(map[string]interface{})
		customInfo["base_storefront_url"] = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStorefrontURL))

		confirmationEmail, err := utils.TextTemplate(confirmationEmail,
			map[string]interface{}{
				"Order":   orderMap,
				"Visitor": visitorMap,
				"Info":    customInfo,
			})
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = app.SendMail(email, "Order confirmation", confirmationEmail)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// CheckoutSuccess will save the order and clear the shopping in the session.
func (it *DefaultCheckout) CheckoutSuccess(checkoutOrder order.InterfaceOrder, session api.InterfaceSession) error {
	var err error

	// making sure order was specified
	if checkoutOrder == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "17d45365-7808-4a1b-ad36-1741a83e820f", "Must specify an Order.")
	}

	// check order status for funds collected before  proceeding to checkout success
	if orderStatus := checkoutOrder.GetStatus(); orderStatus != order.ConstOrderStatusProcessed && orderStatus != order.ConstOrderStatusCompleted {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7dec976e-084b-4d29-9301-bb3b328be95f", "There was an error collecting funds on the order.")
	}

	// checkout information cleanup
	//-----------------------------
	currentCart := it.GetCart()

	if currentCart != nil {
		err = currentCart.Deactivate()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = currentCart.Save()
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	if session != nil {
		session.Set(cart.ConstSessionKeyCurrentCart, nil)
		session.Set(checkout.ConstSessionKeyCurrentCheckout, nil)
		session.Set(coupon.ConstSessionKeyCurrentRedemptions, make([]string, 0))
		session.Set(giftcard.ConstSessionKeyAppliedGiftCardCodes, make([]string, 0))
	}

	// sending notifications
	//----------------------
	eventData := map[string]interface{}{"checkout": it, "order": checkoutOrder, "session": session, "cart": currentCart}
	env.Event("checkout.success", eventData)

	err = it.SendOrderConfirmationMail()
	if err != nil {
		env.LogError(err)
	}

	return nil
}
