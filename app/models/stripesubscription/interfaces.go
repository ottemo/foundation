package stripesubscription

import (
	"github.com/ottemo/foundation/app/models"
)

// InterfaceStripeSubscription represents interface to access business layer implementation of purchase subscription object
type InterfaceStripeSubscription interface {
	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
	models.InterfaceListable
}

// InterfaceSubscriptionCollection represents interface to access business layer implementation of purchase subscription collection
type InterfaceStripeSubscriptionCollection interface {
	ListSubscriptions() []InterfaceStripeSubscription

	models.InterfaceCollection
}
