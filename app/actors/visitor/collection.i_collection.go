package visitor

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// enumerates items of Visitor model type
func (it *DefaultVisitorCollection) List() ([]models.T_ListItem, error) {
	result := make([]models.T_ListItem, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {
		visitorModel, err := visitor.GetVisitorModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		visitorModel.FromHashMap(dbRecordData)

		// retrieving minimal data needed for list
		resultItem := new(models.T_ListItem)

		resultItem.Id = visitorModel.GetId()
		resultItem.Name = visitorModel.GetFullName()
		resultItem.Image = ""
		resultItem.Desc = ""

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = visitorModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// allows to obtain additional attributes from  List() function
func (it *DefaultVisitorCollection) ListAddExtraAttribute(attribute string) error {

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	allowedAttributes := make([]string, 0)
	for _, attributeInfo := range visitorModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}
	allowedAttributes = append(allowedAttributes, "billing_address", "shipping_address")

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew("attribute already in list")
		}
	} else {
		return env.ErrorNew("not allowed attribute")
	}

	return nil
}

// adds selection filter to List() function
func (it *DefaultVisitorCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultVisitorCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// sets select pagination
func (it *DefaultVisitorCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
