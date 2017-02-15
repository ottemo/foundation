package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// List enumerates items of model type
func (it *DefaultCategoryCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	// loading data from DB
	//---------------------
	dbItems, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	// converting db record to StructListItem
	//-----------------------------------
	for _, dbItemData := range dbItems {
		categoryModel, err := category.GetCategoryModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		if err := categoryModel.FromHashMap(dbItemData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5702b4e7-dd12-4424-8618-da8ffd186baf", err.Error())
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		mediaPath, err := categoryModel.GetMediaPath("image")
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		resultItem.ID = categoryModel.GetID()
		resultItem.Name = categoryModel.GetName()
		resultItem.Image = ""
		resultItem.Desc = categoryModel.GetDescription()

		if categoryModel.GetImage() != "" {
			resultItem.Image = mediaPath + categoryModel.GetImage()
		}

		// serving extra attributes
		//-------------------------
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = categoryModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultCategoryCollection) ListAddExtraAttribute(attribute string) error {

	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var allowedAttributes []string
	for _, attributeInfo := range categoryModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}
	allowedAttributes = append(allowedAttributes, "parent")

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2509d847-ba1e-48bd-9b29-37edd0cac52b", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3282704a-a048-4de6-b910-b23c753083a9", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultCategoryCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	if err := it.listCollection.AddFilter(Attribute, Operator, Value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "29c448e4-dcc3-418d-9ab9-6e0af38c9bdc", err.Error())
	}
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultCategoryCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8dc3ed28-42d8-422c-9e2c-ba19486f49eb", err.Error())
	}
	return nil
}

// ListLimit specifies selection paging
func (it *DefaultCategoryCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
