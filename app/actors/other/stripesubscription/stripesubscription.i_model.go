package stripesubscription

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/stripesubscription"
)

// GetModelName returns model name for the Stripe Subscription
func (it *DefaultStripeSubscription) GetModelName() string {
	return stripesubscription.ConstModelNameStripeSubscription
}

// GetImplementationName returns model implementation name for the Stripe Subscription
func (it *DefaultStripeSubscription) GetImplementationName() string {
	return "Default" + stripesubscription.ConstModelNameStripeSubscription
}

// New returns new instance of model implementation object for the Stripe Subscription
func (it *DefaultStripeSubscription) New() (models.InterfaceModel, error) {
	return &DefaultStripeSubscription{}, nil
}
