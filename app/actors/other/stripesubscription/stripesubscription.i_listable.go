package stripesubscription

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/stripesubscription"
)

// GetCollection returns collection of current instance type
func (it *DefaultStripeSubscription) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(stripesubscription.ConstModelNameStripeSubscriptionCollection)
	if result, ok := model.(stripesubscription.InterfaceStripeSubscriptionCollection); ok {
		return result
	}

	return nil
}
