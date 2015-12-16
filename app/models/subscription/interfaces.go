// Package subscription represents abstraction of business layer purchase subscription object
package subscription

import (
	"github.com/ottemo/foundation/app/models"
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

	GetPeriod() int
	SetPeriod(days int) error

	GetStatus() string
	GetState() string
	GetAction() string

	GetAddress() map[string]interface{}

	GetLastSubmit() time.Time
	GetActionDate() time.Time

	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time

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
