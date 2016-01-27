// Package subscription implements base subscription functionality
package subscription

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// Package global constants
const (
	ConstErrorModule = "subscription"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstCollectionNameSubscription = "subscription"

	ConstTimeDay = time.Hour * 24
)

var subscriptionEnabled = false

// DefaultSubscription struct to hold subscription information
type DefaultSubscription struct {
	id string

	VisitorID string
	CartID    string
	OrderID   string

	CustomerEmail string
	CustomerName  string

	Status     string
	State      string
	ActionDate time.Time
	Period     int

	ShippingAddress map[string]interface{}
	BillingAddress  map[string]interface{}

	ShippingMethodCode string

	ShippingRate checkout.StructShippingRate

	// should be stored credit card info with payment method in it
	PaymentInstrument map[string]interface{}

	LastSubmit time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DefaultSubscriptionCollection is a default implementer of InterfaceSubscriptionCollection
type DefaultSubscriptionCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
