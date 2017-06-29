package stock

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/stock"
	"github.com/ottemo/foundation/utils"
)

// ---------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models")
// ---------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultStockCollection) GetModelName() string {
	return ConstCollectionNameStock
}

// GetImplementationName returns model implementation name
func (it *DefaultStockCollection) GetImplementationName() string {
	return "Default" + ConstCollectionNameStock
}

// New returns new instance of model implementation object
func (it *DefaultStockCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultStockCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}

//----------------------------------------------------------------------------------------------------------------------
// InterfaceCollection implementation (package "github.com/ottemo/foundation/app/models/interfaces")
//----------------------------------------------------------------------------------------------------------------------

// GetDBCollection returns database collection
func (it *DefaultStockCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// List enumerates items of model type
func (it *DefaultStockCollection) List() ([]models.StructListItem, error) {
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
		stockModel, err := stock.GetStockModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		err = stockModel.FromHashMap(dbItemData)
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = stockModel.GetID()

		// serving extra attributes
		//-------------------------
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = stockModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultStockCollection) ListAddExtraAttribute(attribute string) error {

	stockModel, err := stock.GetStockModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var allowedAttributes []string
	for _, attributeInfo := range stockModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7d0fc10d-cd7c-47af-9e0e-68bf3be9bb8e", "attribute `" + attribute + "` already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "cf94ab6c-43d0-428d-8a4e-69769d3c9428", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultStockCollection) ListFilterAdd(attribute string, operator string, value interface{}) error {
	if err := it.listCollection.AddFilter(attribute, operator, value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe0b923d-9571-4c64-9a6a-99d2cda3d587", err.Error())
	}
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultStockCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dd115bb7-6c5e-4fc0-86c5-a13a5d12cf92", err.Error())
	}
	return nil
}

// ListLimit specifies selection paging
func (it *DefaultStockCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}


// -----------------------------------------------------------------------------------------------------
// InterfaceSEOCollection implementation (package "github.com/ottemo/foundation/app/models/seo")
// -----------------------------------------------------------------------------------------------------

// ListStocks returns array of stock items in model instance form
func (it *DefaultStockCollection) ListStocks() []stock.InterfaceStock {
	var result []stock.InterfaceStock

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, dbRecordData := range dbRecords {
		stockModel, err := stock.GetStockModel()
		if err != nil {
			return result
		}
		if err := stockModel.FromHashMap(dbRecordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d6811b7c-1d4f-4435-9365-8b10b6bbe9a0", err.Error())
		}

		result = append(result, stockModel)
	}

	return result
}
