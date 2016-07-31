package stripesubscription

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GetStripeSubscriptionCollectionModel retrieves current InterfaceStripeSubscriptionCollection model implementation
func GetStripeSubscriptionCollectionModel() (InterfaceStripeSubscriptionCollection, error) {
	model, err := models.GetModel(ConstModelNameStripeSubscriptionCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stripeSubscriptionCollectionModel, ok := model.(InterfaceStripeSubscriptionCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "68f75930-7147-4bbd-8214-9f70e4a6bb95", "model "+model.GetImplementationName()+" is not 'InterfaceStripeSubscriptionCollection' capable")
	}

	return stripeSubscriptionCollectionModel, nil
}

// GetStripeSubscriptionModel retrieves current InterfaceStripeSubscription model implementation
func GetStripeSubscriptionModel() (InterfaceStripeSubscription, error) {
	model, err := models.GetModel(ConstModelNameStripeSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stripeSubscriptionModel, ok := model.(InterfaceStripeSubscription)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0af5ae4e-0a4e-4840-a270-402c7548980e", "model "+model.GetImplementationName()+" is not 'InterfaceStripeSubscription' capable")
	}

	return stripeSubscriptionModel, nil
}

// LoadStripeSubscriptionByID loads subscription data into current InterfaceStripeSubscription model implementation
func LoadStripeSubscriptionByID(id string) (InterfaceStripeSubscription, error) {

	stripeSubscriptionModel, err := GetStripeSubscriptionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = stripeSubscriptionModel.Load(id)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return stripeSubscriptionModel, nil
}
