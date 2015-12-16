package subscription

import (
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
func (it *DefaultSubscription) GetAction() string {
	return it.Action
}

// GetPeriod returns the Subscription action
func (it *DefaultSubscription) GetPeriod() int {
	return it.Period
}

// GetAddress returns the Subscription address
func (it *DefaultSubscription) GetAddress() map[string]interface{} {
	return it.Address
}

// GetLastSubmit returns the Subscription last submit date
func (it *DefaultSubscription) GetLastSubmit() time.Time {
	return it.CreatedAt
}

// GetActionDate returns the Subscription action date
func (it *DefaultSubscription) GetActionDate() time.Time {
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

// SetPeriod allows to set new period for subscription
func (it *DefaultSubscription) SetPeriod(days int) error {
	it.Period = days
	return nil
}
