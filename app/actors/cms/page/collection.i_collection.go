package page

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// List enumerates items of CMS page model
func (it *DefaultCMSPageCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {
		cmsPageModel, err := cms.GetCMSPageModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		cmsPageModel.FromHashMap(dbRecordData)

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = cmsPageModel.GetID()
		resultItem.Name = cmsPageModel.GetIdentifier()
		resultItem.Image = ""
		resultItem.Desc = cmsPageModel.GetTitle()

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = cmsPageModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultCMSPageCollection) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "_id", "id", "enabled", "identifier", "title", "pagetitle", "content", "created_at", "updated_at") {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8f5544bb-625e-4dec-bedd-d6f19a25f968", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "77bd942c-3fb4-492b-84b9-03de45c8eefa", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultCMSPageCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultCMSPageCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// ListLimit sets select pagination
func (it *DefaultCMSPageCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
