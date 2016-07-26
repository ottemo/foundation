package stripesubscription

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/stripesubscription"
)

// List enumerates items of Stripe Subscription model type in a Stripe Subscription collection
func (it *DefaultStripeSubscriptionCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {

		stripeSubscriptionModel, err := stripesubscription.GetStripeSubscriptionModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		stripeSubscriptionModel.FromHashMap(dbRecordData)

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = stripeSubscriptionModel.GetID()
		resultItem.Name = stripeSubscriptionModel.Get("customer_email")
		resultItem.Image = ""
		resultItem.Desc = utils.InterfaceToString(stripeSubscriptionModel.Get("description"))

		// if extra attributes were required
		if len(it.listExtraAttributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAttributes {
				resultItem.Extra[attributeName] = stripeSubscriptionModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute provides the ability to add additional attributes if the attribute does not already exist
func (it *DefaultStripeSubscriptionCollection) ListAddExtraAttribute(attribute string) error {

	stripeSubscriptionModel := new(DefaultStripeSubscription)

	var allowedAttributes []string
	for _, attributeInfo := range stripeSubscriptionModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAttributes) {
			it.listExtraAttributes = append(it.listExtraAttributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "24125249-ad36-4e0c-b5f1-384829d4a66c", "Attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c8bb9e09-8a5c-49c7-942a-c748aac45663", "Not allowed attribute")
	}

	return nil
}

// ListFilterAdd provides the ability to add a selection filter to List() function
func (it *DefaultStripeSubscriptionCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// ListFilterReset clears the presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultStripeSubscriptionCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// ListLimit sets the pagination when provided offset and limit values
func (it *DefaultStripeSubscriptionCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
