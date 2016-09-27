package seo

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/seo"
)

// ---------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models")
// ---------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultSEOCollection) GetModelName() string {
	return ConstCollectionNameURLRewrites
}

// GetImplementationName returns model implementation name
func (it *DefaultSEOCollection) GetImplementationName() string {
	return "Default" + ConstCollectionNameURLRewrites
}

// New returns new instance of model implementation object
func (it *DefaultSEOCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultSEOCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}

//----------------------------------------------------------------------------------------------------------------------
// InterfaceCollection implementation (package "github.com/ottemo/foundation/app/models/interfaces")
//----------------------------------------------------------------------------------------------------------------------

// GetDBCollection returns database collection
func (it *DefaultSEOCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// List enumerates items of model type
func (it *DefaultSEOCollection) List() ([]models.StructListItem, error) {
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
		seoItemModel, err := seo.GetSEOItemModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		err = seoItemModel.FromHashMap(dbItemData)
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = seoItemModel.GetID()

		// serving extra attributes
		//-------------------------
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = seoItemModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultSEOCollection) ListAddExtraAttribute(attribute string) error {

	seoItemModel, err := seo.GetSEOItemModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var allowedAttributes []string
	for _, attributeInfo := range seoItemModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7d0fc10d-cd7c-47af-9e0e-68bf3be9bb8e", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "cf94ab6c-43d0-428d-8a4e-69769d3c9428", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultSEOCollection) ListFilterAdd(attribute string, operator string, value interface{}) error {
	it.listCollection.AddFilter(attribute, operator, value.(string))
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultSEOCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// ListLimit specifies selection paging
func (it *DefaultSEOCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}

// -----------------------------------------------------------------------------------------------------
// InterfaceSEOCollection implementation (package "github.com/ottemo/foundation/app/models/seo")
// -----------------------------------------------------------------------------------------------------

// ListSEOItems returns array of SEO items in model instance form
func (it *DefaultSEOCollection) ListSEOItems() []seo.InterfaceSEOItem {
	var result []seo.InterfaceSEOItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, dbRecordData := range dbRecords {
		seoItemModel, err := seo.GetSEOItemModel()
		if err != nil {
			return result
		}
		seoItemModel.FromHashMap(dbRecordData)

		result = append(result, seoItemModel)
	}

	return result
}
