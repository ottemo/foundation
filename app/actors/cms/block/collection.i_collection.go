package block

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// List enumerates items of CMS block model
func (it *DefaultCMSBlockCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {
		cmsBlockModel, err := cms.GetCMSBlockModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		if err := cmsBlockModel.FromHashMap(dbRecordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3d84fee6-a9a2-4ea1-bb10-ba4d5fc0d3f6", err.Error())
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = cmsBlockModel.GetID()
		resultItem.Name = cmsBlockModel.GetIdentifier()
		resultItem.Image = ""
		resultItem.Desc = ""

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = cmsBlockModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultCMSBlockCollection) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "_id", "id", "identifier", "content", "created_at", "updated_at") {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "48018ded-44c2-4501-a0bc-bd55a9eaa570", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "018884c4-b2cd-452a-a45e-e88c11e69988", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultCMSBlockCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	if err := it.listCollection.AddFilter(Attribute, Operator, Value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9d1685d0-0a77-4399-86c7-c5956c85157e", err.Error())
	}
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultCMSBlockCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9197a345-bb0d-4903-8160-bb93488f5867", err.Error())
	}
	return nil
}

// ListLimit sets select pagination
func (it *DefaultCMSBlockCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
