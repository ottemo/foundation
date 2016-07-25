package stripesubscription

import (
	"github.com/ottemo/foundation/db"
	"time"
)

// Package global constants
const (
	ConstConfigPathPlans = "payment.stripe.plans"
	ConstErrorModule     = "stripesubscription"
	ConstCollectionNameStripeSubscription = "stripe_subscription"

	ConstModelNameStripeSubscription = "StripeSubscription"
	ConstModelNameStripeSubscriptionCollection = "StripeSubscriptionCollection"
)

// DefaultSubscription struct to hold subscription information and represent
// default implementer of InterfaceSubscription
type DefaultStripeSubscription struct {
	id string

	VisitorID       string
	CustomerName    string
	CustomerEmail   string
	BillingAddress  map[string]interface{}
	ShippingAddress map[string]interface{}

	Total float64

	CreatedAt time.Time
	UpdatedAt time.Time

	Description string
	Info        map[string]interface{}
	Status      string

	StripeSubscriptionID string
	StripeCustomerID     string
}

// DefaultSubscriptionCollection is a default implementer of InterfaceSubscriptionCollection
type DefaultStripeSubscriptionCollection struct {
	listCollection      db.InterfaceDBCollection
	listExtraAttributes []string
}
