package saleprice

import (
	"github.com/ottemo/foundation/app/models"
	salepriceModel "github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultSalePriceCollection) GetModelName() string {
	return salepriceModel.ConstSalePriceDbCollectionName
}

// GetImplementationName default model default implementation name
func (it *DefaultSalePriceCollection) GetImplementationName() string {
	return "Default" + salepriceModel.ConstSalePriceDbCollectionName
}

// New returns new instance of model implementation object
func (it *DefaultSalePriceCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(salepriceModel.ConstSalePriceDbCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultSalePriceCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceCollection implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// GetDBCollection returns database collection
func (it *DefaultSalePriceCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// List returns list of StructListItem items
func (it *DefaultSalePriceCollection) List() ([]models.StructListItem, error) {
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
		salePriceModel, err := salepriceModel.GetSalePriceModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		salePriceModel.FromHashMap(dbItemData)

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = salePriceModel.GetID()
		resultItem.Name = salePriceModel.GetProductID() + ", " +
			utils.InterfaceToString(salePriceModel.GetAmount()) + ", " +
			utils.InterfaceToString(salePriceModel.GetStartDatetime()) + ", " +
			utils.InterfaceToString(salePriceModel.GetEndDatetime())
		resultItem.Image = ""
		resultItem.Desc = "For product [" + salePriceModel.GetProductID() + "], " +
			" set sale price [" + utils.InterfaceToString(salePriceModel.GetAmount()) + "], " +
			" from [" + utils.InterfaceToString(salePriceModel.GetStartDatetime()) + "], " +
			" to [" + utils.InterfaceToString(salePriceModel.GetEndDatetime()) + "]"

		// serving extra attributes
		//-------------------------
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = salePriceModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute adds attribute to sale price collection
func (it *DefaultSalePriceCollection) ListAddExtraAttribute(attribute string) error {

	salePriceModel, err := salepriceModel.GetSalePriceModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var allowedAttributes []string
	for _, attributeInfo := range salePriceModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

	return nil
}

// ListFilterAdd adds filter to sale price collection
func (it *DefaultSalePriceCollection) ListFilterAdd(attribute string, operator string, value interface{}) error {
	it.listCollection.AddFilter(attribute, operator, value.(string))
	return nil
}

// ListFilterReset resets sale price collection filters
func (it *DefaultSalePriceCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// ListLimit limits sale price collection selected records
func (it *DefaultSalePriceCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
