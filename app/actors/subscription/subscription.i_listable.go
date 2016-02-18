package subscription

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/subscription"
)

// GetCollection returns collection of current instance type
func (it *DefaultSubscription) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(subscription.ConstModelNameSubscriptionCollection)
	if result, ok := model.(subscription.InterfaceSubscriptionCollection); ok {
		return result
	}

	return nil
}
