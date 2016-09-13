package saleprice

// SalePriceDelegate type implements:
//	- InterfaceAttributesDelegate

import (
	"time"

	"github.com/ottemo/foundation/api/context"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceAttributesDelegate implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// New creates new SalePriceDelegate with associated product
func (it *SalePriceDelegate) New(productInstance interface{}) (models.InterfaceAttributesDelegate, error) {
	if productModel, ok := productInstance.(product.InterfaceProduct); ok {
		return &SalePriceDelegate{productInstance: productModel}, nil
	}
	return nil, newErrorHelper("unexpected instance for sale price delegate", "6ac6965a-1f1e-44ae-b854-ad430d5b85a6")
}

// Get returns product external attributes managed by sale price package
func (it *SalePriceDelegate) Get(attribute string) interface{} {
	currentContext := context.GetContext()
	isAdmin := false
	if currentContext != nil {
		if contextIsAdmin, ok := currentContext["is_admin"]; ok {
			isAdmin = utils.InterfaceToBool(contextIsAdmin)
		}
	}

	switch attribute {
	case "sale_prices":
		if it.SalePrices == nil {
			salePriceCollectionModel, err := saleprice.GetSalePriceCollectionModel()
			if err != nil {
				newErrorHelper("can not get sale price collection", "1ee355c3-20f4-4706-9723-9fe6c7e1bda4")
				return nil
			}

			salePriceCollectionModel.GetDBCollection().AddFilter("product_id", "=", it.productInstance.GetID())

			salePriceStructListItems, err := salePriceCollectionModel.List()
			if err != nil {
				newErrorHelper("can not get sale prices list", "9d77e24b-e45f-426a-8b7d-dd859271b0d2")
				return nil
			}

			today := time.Now()
			for _, salePriceStructListItem := range salePriceStructListItems {
				salePriceModel, err := saleprice.GetSalePriceModel()
				if err != nil {
					newErrorHelper("can not get sale price model", "d5f43503-d73c-4d60-a349-2668ae37c6b0")
					continue
				}

				err = salePriceModel.Load(salePriceStructListItem.ID)
				if err != nil {
					newErrorHelper("can not load sale price model", "dd08dffe-6147-4d96-8306-c6b60dcb704f")
					continue
				}

				if isAdmin || (salePriceModel.GetStartDatetime().Before(today) &&
					today.Before(salePriceModel.GetEndDatetime())) {
					it.SalePrices = append(it.SalePrices, salePriceModel.ToHashMap())
				}
			}
		}

		return it.SalePrices
	}
	return nil
}

// Set saves product external attributes managed by sale price package
func (it *SalePriceDelegate) Set(attribute string, value interface{}) error {
	var saveError error

	switch attribute {
	case "sale_prices":
		if value != nil {
			// store old records to compare
			salePriceCollectionModel, err := saleprice.GetSalePriceCollectionModel()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			salePriceCollectionModel.GetDBCollection().AddFilter("product_id", "=", it.productInstance.GetID())

			salePriceStructListItems, err := salePriceCollectionModel.List()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			// try to save updated/new records
			newSalePrices := utils.InterfaceToArray(value)
			var newSalePriceIDs []string
			for _, salePrice := range newSalePrices {
				salePriceHashMap := utils.InterfaceToMap(salePrice)

				if salePriceHashMap["_id"] != nil {
					newSalePriceID := utils.InterfaceToString(salePriceHashMap["_id"])
					if newSalePriceID != "" {
						newSalePriceIDs = append(newSalePriceIDs, newSalePriceID)
					}
				}

				salePriceModel, err := saleprice.GetSalePriceModel()
				if err != nil {
					return env.ErrorDispatch(err)
				}

				err = salePriceModel.FromHashMap(salePriceHashMap)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				salePriceModel.SetProductID(it.productInstance.GetID())
				if err != nil {
					return env.ErrorDispatch(err)
				}

				err = salePriceModel.Save()
				// do not exit on save error
				if err != nil {
					// store first error
					if saveError == nil {
						saveError = err
					}
					err = nil
				}

				it.SalePrices = append(it.SalePrices, salePriceModel.ToHashMap())
			}

			// remove old records which are not present
			for _, salePriceStructListItem := range salePriceStructListItems {
				// check new set of records do not include old record
				foundOld := false

				for _, newSalePriceID := range newSalePriceIDs {
					if utils.InterfaceToString(newSalePriceID) ==
						salePriceStructListItem.ID {
						foundOld = true
					}
				}

				// is there are no corresponding new item - remove old one
				if !foundOld {
					salePriceModel, err := saleprice.GetSalePriceModel()
					if err != nil {
						return env.ErrorDispatch(err)
					}

					salePriceModel.SetID(salePriceStructListItem.ID)
					if err != nil {
						return env.ErrorDispatch(err)
					}

					err = salePriceModel.Delete()
					if err != nil {
						return env.ErrorDispatch(err)
					}
				}
			}
		}
	}

	// dispatch first save error
	if saveError != nil {
		return env.ErrorDispatch(saveError)
	}

	return nil
}

// GetAttributesInfo describes product external attributes managed by sale price package
func (it *SalePriceDelegate) GetAttributesInfo() []models.StructAttributeInfo {
	return []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: saleprice.ConstSalePriceDbCollectionName,
			Attribute:  "sale_prices",
			Type:       utils.ConstDataTypeJSON,
			Label:      "SalePrices",
			IsRequired: false,
			IsStatic:   false,
			Group:      "General",
			Editors:    "product_sale_prices",
			Options:    "",
			Default:    "",
			Validators: "",
		},
	}
}
