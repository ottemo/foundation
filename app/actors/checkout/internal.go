package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/actors/discount/coupon"
	"github.com/ottemo/foundation/app/actors/discount/giftcard"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
)

// SendOrderConfirmationEmail sends an order confirmation email
func (it *DefaultCheckout) SendOrderConfirmationEmail() error {

	checkoutOrder := it.GetOrder()
	if checkoutOrder == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e7c69056-cc28-4632-9524-50d71b909d83", "given checkout order does not exists")
	}

	err := checkoutOrder.SendOrderConfirmationEmail()
	if err != nil {
		return env.ErrorDispatch(err)
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

	err = it.SendOrderConfirmationEmail()
	if err != nil {
		_ = env.ErrorDispatch(err)
	}

	return nil
}
