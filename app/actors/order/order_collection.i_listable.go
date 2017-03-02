package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/order"
)

// List enumerates items of Product model type
func (it *DefaultOrderCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {

		orderModel, err := order.GetOrderModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		err = orderModel.FromHashMap(dbRecordData)
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = orderModel.GetID()
		resultItem.Name = orderModel.GetIncrementID()
		resultItem.Image = ""
		resultItem.Desc = utils.InterfaceToString(orderModel.Get("description"))

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = orderModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultOrderCollection) ListAddExtraAttribute(attribute string) error {

	orderModel, err := order.GetOrderModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var allowedAttributes []string
	for _, attributeInfo := range orderModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "43aa5749-55b8-41b6-98b4-851abd2962bd", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "23ab202d-2af5-4bf9-a888-6dc767ef53fe", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultOrderCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	if err := it.listCollection.AddFilter(Attribute, Operator, Value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "57b29f89-94ab-434a-9001-a9fea0b71bd0", err.Error())
	}
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultOrderCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f0900b86-68ee-4e2d-a234-f16e235f1c4f", err.Error())
	}
	return nil
}

// ListLimit specifies selection paging
func (it *DefaultOrderCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
