package stripesubscription

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// Package global constants
const (
	ConstConfigPathGroup   = "general.stripesubscription"
	ConstConfigPathAPIKey  = "general.stripesubscription.apiKey"
	ConstConfigPathEnabled = "general.stripesubscription.enabled"
	ConstConfigPathPlans   = "general.stripesubscription.plans"

	ConstConfigPathEmailCancelSubject     = "general.stripesubscription.emailCancelSubject"
	ConstConfigPathEmailCancelTemplate    = "general.stripesubscription.emailCancelTemplate"
	ConstConfigPathEmailSubscribeSubject  = "general.stripesubscription.emailSubscribeSubject"
	ConstConfigPathEmailSubscribeTemplate = "general.stripesubscription.emailSubscribeTemplate"

	ConstErrorModule                      = "stripesubscription"
	ConstErrorLevel                       = env.ConstErrorLevelActor
	ConstCollectionNameStripeSubscription = "stripe_subscription"

	ConstSubscriptionStatusSuspended = "suspended"
	ConstSubscriptionStatusConfirmed = "confirmed"
	ConstSubscriptionStatusCanceled  = "canceled"
)

// DefaultStripeSubscription struct to hold subscription information and represent
// default implementer of InterfaceStripeSubscription
type DefaultStripeSubscription struct {
	id string

	VisitorID       string
	CustomerName    string
	CustomerEmail   string
	BillingAddress  map[string]interface{}
	ShippingAddress map[string]interface{}

	Description     string
	Status          string
	LastPaymentInfo map[string]interface{}

	StripeCustomerID     string
	StripeSubscriptionID string
	StripeCoupon         string
	Price                float64

	CreatedAt time.Time
	UpdatedAt time.Time
	Info      map[string]interface{}

	PeriodEnd   time.Time
	NotifyRenew bool
}

// DefaultStripeSubscriptionCollection is a default implementer of InterfaceStripeSubscriptionCollection
type DefaultStripeSubscriptionCollection struct {
	listCollection      db.InterfaceDBCollection
	listExtraAttributes []string
}
