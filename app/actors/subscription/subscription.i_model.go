package subscription

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/subscription"
)

// GetModelName returns model name for the Subscription
func (it *DefaultSubscription) GetModelName() string {
	return subscription.ConstModelNameSubscription
}

// GetImplementationName returns model implementation name for the Subscription
func (it *DefaultSubscription) GetImplementationName() string {
	return "Default" + subscription.ConstModelNameSubscription
}

// New returns new instance of model implementation object for the Subscription
func (it *DefaultSubscription) New() (models.InterfaceModel, error) {
	return &DefaultSubscription{Info: make(map[string]interface{})}, nil
}
