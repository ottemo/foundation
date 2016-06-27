package address

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
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
		visitorAddressModel.FromHashMap(dbRecordData)

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
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e26dc487-2873-4fdd-8dff-b9a1005e08ab", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "393c0a07-9acf-4f66-9b59-8de31be541ce", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultVisitorAddressCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultVisitorAddressCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// ListLimit sets select pagination
func (it *DefaultVisitorAddressCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
