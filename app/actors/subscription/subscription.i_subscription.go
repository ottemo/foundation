package subscription

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

// GetEmail returns subscriber e-mail
func (it *DefaultSubscription) GetEmail() string {
	return it.Email
}

// GetName returns name of subscriber
func (it *DefaultSubscription) GetName() string {
	return it.Name
}

// GetVisitorID returns the Subscription's Visitor ID
func (it *DefaultSubscription) GetVisitorID() string {
	return it.VisitorID
}

// GetCartID returns the Subscription's Cart ID
func (it *DefaultSubscription) GetCartID() string {
	return it.CartID
}

// GetOrderID returns the Subscription's Order ID
func (it *DefaultSubscription) GetOrderID() string {
	return it.OrderID
}

// GetStatus returns the Subscription status
func (it *DefaultSubscription) GetStatus() string {
	return it.Status
}

// GetState returns the Subscription state
func (it *DefaultSubscription) GetState() string {
	return it.State
}

// GetAction returns the Subscription action
func (it *DefaultSubscription) GetActionDate() time.Time {
	return it.ActionDate
}

// GetPeriod returns the Subscription action
func (it *DefaultSubscription) GetPeriod() int {
	return it.Period
}

// SetStatus set Subscription status
func (it *DefaultSubscription) SetStatus(status string) error {
	if status != ConstSubscriptionStatusSuspended && status != ConstSubscriptionStatusConfirmed && status != ConstSubscriptionStatusCanceled {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3b7d17c3-c5fa-4369-a039-49bafec2fb9d", "new subscription status should be one of allowed")
	}

	if it.Status == status {
		return nil
	}

	it.Status = status
	return nil
}

// SetState set Subscription state
func (it *DefaultSubscription) SetState(state string) error {
	it.State = state
	return nil
}

// SetActionDate set Subscription action date
func (it *DefaultSubscription) SetActionDate(actionDate time.Time) error {
	if err := validateSubscriptionDate(actionDate); err != nil {
		return env.ErrorDispatch(err)
	}

	it.ActionDate = actionDate
	return nil
}

// SetPeriod set Subscription period
func (it *DefaultSubscription) SetPeriod(days int) error {
	if err := validateSubscriptionPeriod(days); err != nil {
		return env.ErrorDispatch(err)
	}

	it.Period = days
	return nil
}

// SetShippingAddress sets shipping address for subscription
func (it *DefaultSubscription) SetShippingAddress(address visitor.InterfaceVisitorAddress) error {
	if address == nil {
		it.ShippingAddress = nil
		return nil
	}

	it.ShippingAddress = address.ToHashMap()
	return nil
}

// GetShippingAddress returns subscription shipping address
func (it *DefaultSubscription) GetShippingAddress() visitor.InterfaceVisitorAddress {
	if it.ShippingAddress == nil {
		return nil
	}

	shippingAddress, err := visitor.GetVisitorAddressModel()
	if err != nil {
		env.LogError(err)
		return nil
	}

	err = shippingAddress.FromHashMap(it.ShippingAddress)
	if err != nil {
		env.LogError(err)
		return nil
	}

	return shippingAddress
}

// SetBillingAddress sets billing address for subscription
func (it *DefaultSubscription) SetBillingAddress(address visitor.InterfaceVisitorAddress) error {
	if address == nil {
		it.BillingAddress = nil
		return nil
	}

	it.BillingAddress = address.ToHashMap()
	return nil
}

// GetBillingAddress returns subscription billing address
func (it *DefaultSubscription) GetBillingAddress() visitor.InterfaceVisitorAddress {
	if it.BillingAddress == nil {
		return nil
	}

	billingAddress, err := visitor.GetVisitorAddressModel()
	if err != nil {
		env.LogError(err)
		return nil
	}

	err = billingAddress.FromHashMap(it.BillingAddress)
	if err != nil {
		env.LogError(err)
		return nil
	}

	return billingAddress
}

// SetPaymentMethod sets payment method for subscription
func (it *DefaultSubscription) SetCreditCard(creditCard visitor.InterfaceVisitorCard) error {
	if creditCard == nil {
		it.PaymentInstrument = nil
		return nil
	}

	it.PaymentInstrument = creditCard.ToHashMap()
	return nil
}

// GetCreditCard sets payment method for subscription
func (it *DefaultSubscription) GetCreditCard() visitor.InterfaceVisitorCard {

	visitorCardModel, err := visitor.GetVisitorCardModel()
	if err != nil {
		env.LogError(err)
		return nil
	}

	if it.PaymentInstrument == nil {
		return visitorCardModel
	}

	err = visitorCardModel.FromHashMap(it.PaymentInstrument)
	if err != nil {
		env.LogError(err)
		return visitorCardModel
	}

	return visitorCardModel
}

// GetPaymentMethod returns subscription payment method
func (it *DefaultSubscription) GetPaymentMethod() checkout.InterfacePaymentMethod {
	creditCard := it.GetCreditCard()
	if creditCard == nil {
		return nil
	}

	return checkout.GetPaymentMethodByCode(creditCard.GetPaymentMethodCode())
}

// SetShippingMethod sets payment method for subscription
func (it *DefaultSubscription) SetShippingMethod(shippingMethod checkout.InterfaceShippingMethod) error {
	it.ShippingMethodCode = shippingMethod.GetCode()
	return nil
}

// GetShippingMethod returns a subscription shipping method
func (it *DefaultSubscription) GetShippingMethod() checkout.InterfaceShippingMethod {
	return checkout.GetShippingMethodByCode(it.ShippingMethodCode)
}

// SetShippingRate sets shipping rate for subscription
func (it *DefaultSubscription) SetShippingRate(shippingRate checkout.StructShippingRate) error {
	it.ShippingRate = shippingRate
	return nil
}

// GetShippingRate returns a subscription shipping rate
func (it *DefaultSubscription) GetShippingRate() checkout.StructShippingRate {
	return it.ShippingRate
}

// GetLastSubmit returns the Subscription last submit date
func (it *DefaultSubscription) GetLastSubmit() time.Time {
	return it.CreatedAt
}

// GetCreatedAt returns the Subscription creation date
func (it *DefaultSubscription) GetCreatedAt() time.Time {
	return it.CreatedAt
}

// GetUpdatedAt returns the Subscription update date
func (it *DefaultSubscription) GetUpdatedAt() time.Time {
	return it.CreatedAt
}

// Validate allows to validate subscription object for data presence
func (it *DefaultSubscription) GetCheckout() (checkout.InterfaceCheckout, error) {

	checkoutInstance, err := checkout.GetCheckoutModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// set visitor basic info
	visitorID := it.GetVisitorID()
	if visitorID != "" {
		checkoutInstance.Set("VisitorID", visitorID)
	}

	checkoutInstance.SetInfo("customer_email", it.GetEmail())
	checkoutInstance.SetInfo("customer_name", it.GetName())

	// set billing and shipping address
	shippingAddress := it.GetShippingAddress()
	if shippingAddress != nil {
		checkoutInstance.Set("ShippingAddress", shippingAddress.ToHashMap())
	}

	billingAddress := it.GetBillingAddress()
	if billingAddress != nil {
		checkoutInstance.Set("BillingAddress", billingAddress.ToHashMap())
	}

	// check payment and shipping methods for availability
	shippingMethod := it.GetShippingMethod()
	if shippingMethod != nil {
		if !shippingMethod.IsAllowed(checkoutInstance) {
			return checkoutInstance, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "db2e8933-d0eb-4a16-a28b-78c169fe20c0", "shipping method not allowed")
		}

		err = checkoutInstance.SetShippingMethod(shippingMethod)
		if err != nil {
			return checkoutInstance, env.ErrorDispatch(err)
		}

		err = checkoutInstance.SetShippingRate(it.GetShippingRate())
		if err != nil {
			return checkoutInstance, env.ErrorDispatch(err)
		}
	}

	paymentMethod := it.GetPaymentMethod()
	if paymentMethod != nil {
		if !paymentMethod.IsAllowed(checkoutInstance) {
			return checkoutInstance, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "e7cfc56b-97d9-43f5-862e-fb370004c8cf", "payment method not allowed")
		}

		err = checkoutInstance.SetPaymentMethod(paymentMethod)
		if err != nil {
			return checkoutInstance, env.ErrorDispatch(err)
		}
	}

	checkoutInstance.SetInfo("cc", it.GetCreditCard())

	// handle cart
	currentCart, err := cart.LoadCartByID(it.GetCartID())
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}

	err = currentCart.ValidateCart()
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}

	err = checkoutInstance.SetCart(currentCart)
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}

	return checkoutInstance, nil
}

// Validate allows to validate subscription object for data presence
// TODO: validate ALL values and thre exisitng
func (it *DefaultSubscription) Validate() error {
	if err := validateSubscriptionPeriod(it.Period); err != nil {
		return env.ErrorDispatch(err)
	}

	if err := validateSubscriptionDate(it.ActionDate); err != nil {
		return env.ErrorDispatch(err)
	}

	if !utils.ValidEmailAddress(it.Email) {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "1c033c36-d63b-4659-95e8-9f348f5e2880", "Subscription invalid: email")
	}

	if it.CartID == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "1c033c36-d63b-4659-95e8-9f348f5e2880", "Subscription missing: cart_id")
	}

	return nil
}
