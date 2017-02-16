package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// List enumerates items of Product model type
func (it *DefaultOrderItemCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {

		orderItemModel := new(DefaultOrderItem)
		err = orderItemModel.FromHashMap(dbRecordData)
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = orderItemModel.GetID()
		resultItem.Name = orderItemModel.GetName()
		resultItem.Image = ""
		resultItem.Desc = utils.InterfaceToString(orderItemModel.Get("description"))

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = orderItemModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultOrderItemCollection) ListAddExtraAttribute(attribute string) error {

	orderItemModel := new(DefaultOrderItem)

	var allowedAttributes []string
	for _, attributeInfo := range orderItemModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fa1d634a-c6c0-4ac6-98cf-e3c5a320c46c", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "205b9f78-3574-41a6-9384-b3f7dd482328", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultOrderItemCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	if err := it.listCollection.AddFilter(Attribute, Operator, Value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7d34fbe6-ac90-430b-8ede-9768982bac7d", err.Error())
	}
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultOrderItemCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a6f914e9-a6ae-4b0b-be9a-21a4fbf81a1d", err.Error())
	}
	return nil
}

// ListLimit specifies selection paging
func (it *DefaultOrderItemCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
