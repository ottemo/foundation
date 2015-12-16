package subscription

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GetSubscriptionCollectionModel retrieves current InterfaceSubscriptionCollection model implementation
func GetSubscriptionCollectionModel() (InterfaceSubscriptionCollection, error) {
	model, err := models.GetModel(ConstModelNameSubscriptionCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionCollectionModel, ok := model.(InterfaceSubscriptionCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "954efbe7-9d4c-4072-8ef2-850ecf5f17a8", "model "+model.GetImplementationName()+" is not 'InterfaceSubscriptionCollection' capable")
	}

	return subscriptionCollectionModel, nil
}

// GetSubscriptionModel retrieves current InterfaceSubscription model implementation
func GetSubscriptionModel() (InterfaceSubscription, error) {
	model, err := models.GetModel(ConstModelNameSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionModel, ok := model.(InterfaceSubscription)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a4e0418a-c508-42de-b994-e6ad08fd796a", "model "+model.GetImplementationName()+" is not 'InterfaceSubscription' capable")
	}

	return subscriptionModel, nil
}

// LoadSubscriptionByID loads subscription data into current InterfaceSubscription model implementation
func LoadSubscriptionByID(subscriptionID string) (InterfaceSubscription, error) {

	subscriptionModel, err := GetSubscriptionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = subscriptionModel.Load(subscriptionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return subscriptionModel, nil
}
