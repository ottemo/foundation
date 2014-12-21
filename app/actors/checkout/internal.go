package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// SendOrderConfirmationMail sends an order confirmation email
func (it *DefaultCheckout) SendOrderConfirmationMail() error {

	checkoutOrder := it.GetOrder()
	if checkoutOrder == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e7c69056cc284632952450d71b909d83", "given checkout order does not exists")
	}

	confirmationEmail := utils.InterfaceToString(env.ConfigGetValue(checkout.ConstConfigPathConfirmationEmail))
	if confirmationEmail != "" {
		email := utils.InterfaceToString(checkoutOrder.Get("customer_email"))
		if email == "" {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1202fcfbda3f4a0f9a2e92f288fd3881", "customer email for order is not set")
		}

		visitorMap := make(map[string]interface{})
		if visitorModel := it.GetVisitor(); visitorModel != nil {
			visitorMap = visitorModel.ToHashMap()
		} else {
			visitorMap["first_name"] = checkoutOrder.Get("customer_name")
			visitorMap["email"] = checkoutOrder.Get("customer_email")
		}

		confirmationEmail, err := utils.TextTemplate(confirmationEmail,
			map[string]interface{}{
				"Order":   checkoutOrder.ToHashMap(),
				"Visitor": visitorMap,
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

	if checkoutOrder == nil || session == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "17d4536578084a1bad361741a83e820f", "Order or session is null")
	}

	err := checkoutOrder.NewIncrementID()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = checkoutOrder.SetStatus(order.ConstOrderStatusPending)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = checkoutOrder.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// cleanup checkout information
	//-----------------------------
	currentCart := it.GetCart()

	err = currentCart.Deactivate()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = currentCart.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	session.Set(cart.ConstSessionKeyCurrentCart, nil)
	session.Set(checkout.ConstSessionKeyCurrentCheckout, nil)

	// sending notifications
	//-----------------------------
	eventData := map[string]interface{}{"checkout": it, "order": checkoutOrder, "session": session, "cart": currentCart}
	env.Event("checkout.success", eventData)

	err = it.SendOrderConfirmationMail()
	if err != nil {
		env.ErrorDispatch(err)
	}

	return nil
}
