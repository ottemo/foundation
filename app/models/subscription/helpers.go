package subscription

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
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

// IsSubscriptionEnabled return status of subscription
func IsSubscriptionEnabled() bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathSubscriptionEnabled))
}

// ContainsSubscriptionItems used to check checkout for subscription items
func ContainsSubscriptionItems(checkoutInstance checkout.InterfaceCheckout) bool {
	currentCart := checkoutInstance.GetCart()
	if currentCart == nil {
		return false
	}

	for _, cartItem := range currentCart.GetItems() {
		itemOptions := cartItem.GetOptions()
		if optionValue, present := itemOptions[ConstSubscriptionOptionName]; present {
			if GetSubscriptionPeriodValue(utils.InterfaceToString(optionValue)) != 0 {
				return true
			}
		}
	}

	return false
}

// GetSubscriptionPeriodValue used to obtain valid period value from option value
func GetSubscriptionPeriodValue(option string) int {

	if value, present := optionValues[option]; present {
		return value
	}

	if value, present := optionValues[strings.ToLower(option)]; present {
		return value
	}

	return 0
}

// GetSubscriptionOptionValues return map of known options for subscription
func GetSubscriptionOptionValues() map[string]int {
	return optionValues
}

// GetSubscriptionCronExpr return cron expression by option value
func GetSubscriptionCronExpr(key int) string {
	if cronExprValue, present := cronExpr[key]; present {
		return cronExprValue
	}
	return cronExpr[optionValues[ConstConfigPathSubscriptionExecutionOptionHour]]
}
