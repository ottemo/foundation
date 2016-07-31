package stripesubscription

import (
	"time"
)

// GetVisitorID returns the Subscription's Visitor ID
func (it *DefaultStripeSubscription) GetVisitorID() string {
	return it.VisitorID
}

// GetCustomerEmail returns subscriber e-mail
func (it *DefaultStripeSubscription) GetCustomerEmail() string {
	return it.CustomerEmail
}

// GetPeriodEnd returns subscription current period end time
func (it *DefaultStripeSubscription) GetPeriodEnd() time.Time {
	return it.PeriodEnd
}

// GetStripeSubscriptionID returns subscription ID in Stripe
func (it *DefaultStripeSubscription) GetStripeSubscriptionID() string {
	return it.StripeSubscriptionID
}
