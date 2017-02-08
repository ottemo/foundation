package subscription

import (
	"time"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetCustomerEmail returns subscriber e-mail
func (it *DefaultSubscription) GetCustomerEmail() string {
	return it.CustomerEmail
}

// GetCustomerName returns name of subscriber
func (it *DefaultSubscription) GetCustomerName() string {
	return it.CustomerName
}

// GetVisitorID returns the Subscription's Visitor ID
func (it *DefaultSubscription) GetVisitorID() string {
	return it.VisitorID
}

// GetOrderID returns the Subscription's Order ID
func (it *DefaultSubscription) GetOrderID() string {
	return it.OrderID
}

// GetStatus returns the Subscription status
func (it *DefaultSubscription) GetStatus() string {
	return it.Status
}

// GetActionDate returns the date of next action
func (it *DefaultSubscription) GetActionDate() time.Time {
	return it.ActionDate
}

// GetPeriod returns the Subscription action
func (it *DefaultSubscription) GetPeriod() int {
	return it.Period
}

// SetStatus set Subscription status
func (it *DefaultSubscription) SetStatus(status string) error {
	if status != subscription.ConstSubscriptionStatusSuspended && status != subscription.ConstSubscriptionStatusConfirmed && status != subscription.ConstSubscriptionStatusCanceled {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3b7d17c3-c5fa-4369-a039-49bafec2fb9d", "new subscription status should be one of allowed")
	}

	if it.Status == status {
		return nil
	}

	it.Status = status
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

// UpdateActionDate set Subscription action date
func (it *DefaultSubscription) UpdateActionDate() error {

	actionDate := it.GetActionDate()
	periodValue := it.GetPeriod()
	if periodValue > 0 {
		actionDate = actionDate.Add(ConstTimeDay * time.Duration(periodValue))
	} else {
		actionDate = actionDate.Add(time.Hour * time.Duration(periodValue*-1))
	}

	return it.SetActionDate(actionDate)
}

// SetPeriod set Subscription period
func (it *DefaultSubscription) SetPeriod(days int) error {
	if err := validateSubscriptionPeriod(days); err != nil {
		return env.ErrorDispatch(err)
	}

	it.Period = days
	return nil
}

// GetItems return items of subscription
func (it *DefaultSubscription) GetItems() []subscription.StructSubscriptionItem {
	return it.items
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
		_ = env.ErrorDispatch(err)
		return nil
	}

	err = shippingAddress.FromHashMap(it.ShippingAddress)
	if err != nil {
		_ = env.ErrorDispatch(err)
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
		_ = env.ErrorDispatch(err)
		return nil
	}

	err = billingAddress.FromHashMap(it.BillingAddress)
	if err != nil {
		_ = env.ErrorDispatch(err)
		return nil
	}

	return billingAddress
}

// SetCreditCard sets credit card for subscription
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
		_ = env.ErrorDispatch(err)
		return nil
	}

	if it.PaymentInstrument == nil {
		return visitorCardModel
	}

	err = visitorCardModel.FromHashMap(it.PaymentInstrument)
	if err != nil {
		_ = env.ErrorDispatch(err)
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

// GetCheckout return checkout object created from subscription
func (it *DefaultSubscription) GetCheckout() (checkout.InterfaceCheckout, error) {

	checkoutInstance, err := checkout.GetCheckoutModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// set visitor basic info
	visitorID := it.GetVisitorID()
	if visitorID != "" {
		if err := checkoutInstance.Set("VisitorID", visitorID); err != nil {
			return checkoutInstance, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "205c5663-e5d7-404e-b731-7de78c055fac", err.Error())
		}
	}

	if err := checkoutInstance.SetInfo("customer_email", it.GetCustomerEmail()); err != nil {
		return checkoutInstance, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "12ef8ee8-0f77-4f1f-9a65-cc94648cf981", err.Error())
	}
	if err := checkoutInstance.SetInfo("customer_name", it.GetCustomerName()); err != nil {
		return checkoutInstance, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "7054a0f4-294c-47a3-8921-96a4406b7ed4", err.Error())
	}

	// set billing and shipping address
	shippingAddress := it.GetShippingAddress()
	if shippingAddress != nil {
		if err := checkoutInstance.Set("ShippingAddress", shippingAddress.ToHashMap()); err != nil {
			return checkoutInstance, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "37424b34-0406-4361-a882-43514292c5bd", err.Error())
		}
	}

	billingAddress := it.GetBillingAddress()
	if billingAddress != nil {
		if err := checkoutInstance.Set("BillingAddress", billingAddress.ToHashMap()); err != nil {
			return checkoutInstance, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "f07174df-4eb7-4447-9b3b-338076516369", err.Error())
		}
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

	if err := checkoutInstance.SetInfo("cc", it.GetCreditCard()); err != nil {
		return checkoutInstance, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "4b47d1a1-577d-4ef0-a57a-8e7f00679b72", err.Error())
	}

	// handle cart
	currentCart, err := cart.GetCartModel()
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}

	for _, item := range it.GetItems() {
		if _, err = currentCart.AddItem(item.ProductID, item.Qty, item.Options); err != nil {
			return checkoutInstance, env.ErrorDispatch(err)
		}
	}

	if err = currentCart.ValidateCart(); err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}

	if err = currentCart.Save(); err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}

	if err = checkoutInstance.SetCart(currentCart); err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}

	checkoutInstance.GetGrandTotal()

	return checkoutInstance, nil
}

// Validate allows to validate subscription object for data presence
// TODO: validate ALL values and their existence
func (it *DefaultSubscription) Validate() error {
	if err := validateSubscriptionPeriod(it.Period); err != nil {
		return env.ErrorDispatch(err)
	}

	if err := validateSubscriptionDate(it.ActionDate); err != nil {
		return env.ErrorDispatch(err)
	}

	if !utils.ValidEmailAddress(it.CustomerEmail) {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "1c033c36-d63b-4659-95e8-9f348f5e2880", "Subscription invalid: email")
	}

	if len(it.items) == 0 {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "1c033c36-d63b-4659-95e8-9f348f5e2880", "no items in subscription")
	}

	return nil
}

// SetInfo sets additional info for subscription
func (it *DefaultSubscription) SetInfo(key string, value interface{}) {
	if value == nil {
		if _, present := it.Info[key]; present {
			delete(it.Info, key)
		}
	} else {
		it.Info[key] = value
	}

	return
}

// GetInfo returns additional subscription info value or nil,
//   - use "*" as a key to get all keys and values currently set
func (it *DefaultSubscription) GetInfo(key string) interface{} {
	if key == "*" {
		return it.Info
	}

	if value, present := it.Info[key]; present {
		return value
	}

	return nil
}
