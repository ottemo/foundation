package address

import (
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"

	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/visitor"
)

// List enumerates items of VisitorAddress model type
func (it *DefaultVisitorAddressCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {
		visitorAddressModel, err := visitor.GetVisitorAddressModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		if err := visitorAddressModel.FromHashMap(dbRecordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2891345a-2f9f-4cec-b0f6-9bd9e6533455", err.Error())
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = visitorAddressModel.GetID()
		resultItem.Name =
			visitorAddressModel.GetAddress() + ", " +
				visitorAddressModel.GetCity() + ", " +
				visitorAddressModel.GetState() + ", " +
				visitorAddressModel.GetZipCode() + ", " +
				visitorAddressModel.GetCountry()

		resultItem.Image = ""
		resultItem.Desc = "Zip: " + visitorAddressModel.GetZipCode() + ", State: " + visitorAddressModel.GetState() +
			", City: " + visitorAddressModel.GetCity() + ", Address: " + visitorAddressModel.GetAddress() +
			", Phone: " + visitorAddressModel.GetPhone()

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = visitorAddressModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultVisitorAddressCollection) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "_id", "id", "visitor_id", "visitorID", "fname", "first_name", "lname", "last_name",
		"address_line1", "address_line2", "company", "country", "city", "state", "phone", "zip", "zip_code") {

		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e58839af-a0a2-405e-aab9-ad9eea5768c5", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "97e1a1a1-5ce1-4137-9223-ab1881917dbc", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultVisitorAddressCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	if err := it.listCollection.AddFilter(Attribute, Operator, Value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5267288e-5aa0-4656-94b8-ce49013fd038", err.Error())
	}
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultVisitorAddressCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3dbf4efc-2b7b-4a92-821b-866923312008", err.Error())
	}
	return nil
}

// ListLimit sets select pagination
func (it *DefaultVisitorAddressCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
