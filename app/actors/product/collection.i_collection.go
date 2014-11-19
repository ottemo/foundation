package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/product"
)

// enumerates items of Product model type
func (it *DefaultProductCollection) List() ([]models.StructListItem, error) {
	result := make([]models.StructListItem, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {

		productModel, err := product.GetProductModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		err = productModel.FromHashMap(dbRecordData)
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		mediaPath, err := productModel.GetMediaPath("image")
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		resultItem.ID = productModel.GetID()
		resultItem.Name = "[" + productModel.GetSku() + "] " + productModel.GetName()
		resultItem.Image = ""
		resultItem.Desc = productModel.GetShortDescription()

		if productModel.GetDefaultImage() != "" {
			resultItem.Image = mediaPath + productModel.GetDefaultImage()
		}

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = productModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// allows to obtain additional attributes from  List() function
func (it *DefaultProductCollection) ListAddExtraAttribute(attribute string) error {

	productModel, err := product.GetProductModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	allowedAttributes := make([]string, 0)
	for _, attributeInfo := range productModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

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
func (it *DefaultProductCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultProductCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// specifies selection paging
func (it *DefaultProductCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
