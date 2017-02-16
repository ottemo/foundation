package product

// DefaultProductCollection type implements:
//	- InterfaceModel
//	- InterfaceCollection
//	- InterfaceProductCollection

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// --------------------------------------------------------------------------------------
// InterfaceCollection implementation (package "github.com/ottemo/foundation/app/models")
// --------------------------------------------------------------------------------------

// List enumerates items of Product model type
func (it *DefaultProductCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

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
			// load external attributes
			err = productModel.LoadExternalAttributes()
			if err != nil {
				return result, env.ErrorDispatch(err)
			}

			// populate required attribute values
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = productModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultProductCollection) ListAddExtraAttribute(attribute string) error {

	productModel, err := product.GetProductModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var allowedAttributes []string
	for _, attributeInfo := range productModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7bb09474-0f5c-4eae-8b19-9573e17806f8", "attribute '"+attribute+"' already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fdb287f6-200f-4867-9822-d3939d098f19", "not allowed attribute '"+attribute+"'")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultProductCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	if err := it.listCollection.AddFilter(Attribute, Operator, Value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "84234105-3560-447c-bcc8-6ccf238dba9c", err.Error())
	}
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultProductCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ffaccc12-e45a-4cd9-a79e-168e89b0d2c1", err.Error())
	}
	return nil
}

// ListLimit specifies selection paging
func (it *DefaultProductCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}

// ---------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models")
// ---------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultProductCollection) GetModelName() string {
	return product.ConstModelNameProductCollection
}

// GetImplementationName returns model implementation name
func (it *DefaultProductCollection) GetImplementationName() string {
	return "Default" + product.ConstModelNameProductCollection
}

// New returns new instance of model implementation object
func (it *DefaultProductCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultProductCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}

// -----------------------------------------------------------------------------------------------------
// InterfaceProductCollection implementation (package "github.com/ottemo/foundation/app/models/product")
// -----------------------------------------------------------------------------------------------------

// GetDBCollection returns database collection
func (it *DefaultProductCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListProducts returns array of products in model instance form
func (it *DefaultProductCollection) ListProducts() []product.InterfaceProduct {
	var result []product.InterfaceProduct

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, dbRecordData := range dbRecords {
		productModel, err := product.GetProductModel()
		if err != nil {
			return result
		}
		if err := productModel.FromHashMap(dbRecordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c67e86dd-d87e-49eb-a421-ba53a0d7157e", err.Error())
		}

		// load external attributes
		err = productModel.LoadExternalAttributes()
		if err != nil {
			return result
		}

		result = append(result, productModel)
	}

	return result
}
