package stripesubscription

import (
	"github.com/ottemo/foundation/app/models"
)

// GetModelName returns model name for the Stripe Subscription
func (it *DefaultStripeSubscription) GetModelName() string {
	return ConstModelNameStripeSubscription
}

// GetImplementationName returns model implementation name for the Stripe Subscription
func (it *DefaultStripeSubscription) GetImplementationName() string {
	return "Default" + ConstModelNameStripeSubscription
}

// New returns new instance of model implementation object for the Stripe Subscription
func (it *DefaultStripeSubscription) New() (models.InterfaceModel, error) {
	return &DefaultStripeSubscription{}, nil
}
