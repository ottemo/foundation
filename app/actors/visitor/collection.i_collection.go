package visitor

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// List enumerates items of Visitor model type in a Visitor collection
func (it *DefaultVisitorCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {
		visitorModel, err := visitor.GetVisitorModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		if err := visitorModel.FromHashMap(dbRecordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e46fc881-1820-4eaf-9f06-6452be240cfd", err.Error())
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = visitorModel.GetID()
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

// ListAddExtraAttribute provides the ability to add additional attributes if the attribute does not already exist
func (it *DefaultVisitorCollection) ListAddExtraAttribute(attribute string) error {

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var allowedAttributes []string
	for _, attributeInfo := range visitorModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}
	allowedAttributes = append(allowedAttributes, "billing_address", "shipping_address", "token")

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b3700f6c-65a1-4fd0-b8f5-37af1b6922a7", "The provided attribute, "+attribute+", already exists, unable to insert it.")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a4bee81b-f975-40a2-94cf-fdbdfbf80627", "The following attribute is not allowed to be added to the attribute list: "+attribute+".")
	}

	return nil
}

// ListFilterAdd provides the ability to add a selection filter to List() function
func (it *DefaultVisitorCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	if err := it.listCollection.AddFilter(Attribute, Operator, Value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "46acf901-4b7b-49b4-a97c-b1d1e6257ebe", err.Error())
	}
	return nil
}

// ListFilterReset clears the presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultVisitorCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a0d1dbbb-e482-469c-8685-12712128b371", err.Error())
	}
	return nil
}

// ListLimit sets the pagination when provided offset and limit values
func (it *DefaultVisitorCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
