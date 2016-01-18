// Package subscription represents abstraction of business layer purchase subscription object
package subscription

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"time"
)

// Package global constants
const (
	ConstModelNameSubscription           = "Subscription"
	ConstModelNameSubscriptionCollection = "SubscriptionCollection"

	ConstErrorModule = "subscription"
	ConstErrorLevel  = env.ConstErrorLevelModel
)

// InterfaceSubscription represents interface to access business layer implementation of purchase subscription object
type InterfaceSubscription interface {
	GetEmail() string
	GetName() string

	GetOrderID() string
	GetCartID() string
	GetVisitorID() string

	SetShippingAddress(address visitor.InterfaceVisitorAddress) error
	GetShippingAddress() visitor.InterfaceVisitorAddress

	SetBillingAddress(address visitor.InterfaceVisitorAddress) error
	GetBillingAddress() visitor.InterfaceVisitorAddress

	SetCreditCard(creditCard visitor.InterfaceVisitorCard) error
	GetCreditCard() visitor.InterfaceVisitorCard

	GetPaymentMethod() checkout.InterfacePaymentMethod

	SetShippingMethod(shippingMethod checkout.InterfaceShippingMethod) error
	GetShippingMethod() checkout.InterfaceShippingMethod

	SetShippingRate(shippingRate checkout.StructShippingRate) error
	GetShippingRate() checkout.StructShippingRate

	GetStatus() string
	SetStatus(status string) error

	GetActionDate() time.Time
	SetActionDate(actionDate time.Time) error

	GetPeriod() int
	SetPeriod(days int) error

	UpdateActionDate() error

	GetLastSubmit() time.Time
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time

	Validate() error
	GetCheckout() (checkout.InterfaceCheckout, error)

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
	models.InterfaceListable
}

// InterfaceSubscriptionCollection represents interface to access business layer implementation of purchase subscription collection
type InterfaceSubscriptionCollection interface {
	ListSubscriptions() []InterfaceSubscription

	models.InterfaceCollection
}
