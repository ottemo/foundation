package token

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// List enumerates items of VisitorAddress model type
func (it *DefaultVisitorCardCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {
		visitorCardModel, err := visitor.GetVisitorCardModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		if err := visitorCardModel.FromHashMap(dbRecordData); err != nil {
			_ = env.ErrorDispatch(err)
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = visitorCardModel.GetID()
		resultItem.Name = visitorCardModel.GetType() + " (" + visitorCardModel.GetNumber() + ")"
		resultItem.Image = ""
		resultItem.Desc = visitorCardModel.GetPaymentMethodCode()

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = visitorCardModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultVisitorCardCollection) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "_id", "id", "visitor_id", "visitorID", "name", "holder", "number", "type", "payment", "expiration", "expiration_date", "created_at", "token_updated", "customer_id") {

		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e26dc487-2873-4fdd-8dff-b9a1005e08ab", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "393c0a07-9acf-4f66-9b59-8de31be541ce", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultVisitorCardCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	if err := it.listCollection.AddFilter(Attribute, Operator, Value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "023f1269-98e9-4d73-ab37-41368ed117f2", "Unable to add filter: " + err.Error())
	}
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultVisitorCardCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "51cb0061-1f1c-4319-88bb-f89b640d2347", "Unable to clear filters: " + err.Error())
	}
	return nil
}

// ListLimit sets select pagination
func (it *DefaultVisitorCardCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
