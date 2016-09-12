package saleprice

// SalePriceDelegate type implements:
//	- InterfaceAttributesDelegate

import (
	"time"

	contextPkg "github.com/ottemo/foundation/api/context"
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
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6ac6965a-1f1e-44ae-b854-ad430d5b85a6", "unexpected instance for sale price delegate")
}

// Get returns product external attributes managed by sale price package
func (it *SalePriceDelegate) Get(attribute string) interface{} {
	context := contextPkg.GetContext()
	isAdmin := false
	if context != nil {
		if isAdminContext, ok := context["is_admin"]; ok {
			isAdmin = utils.InterfaceToBool(isAdminContext)
		}
	}

	switch attribute {
	case "sale_prices":
		if it.SalePrices == nil {
			salePriceCollectionModel, err := saleprice.GetSalePriceCollectionModel()
			if err != nil {
				logWarnHelper("Can not get sale price collection [1ee355c3-20f4-4706-9723-9fe6c7e1bda4]")
				return nil
			}

			salePriceCollectionModel.GetDBCollection().AddFilter("product_id", "=", it.productInstance.GetID())

			salePriceStructListItems, err := salePriceCollectionModel.List()
			if err != nil {
				logWarnHelper("Can not get sale prices list [9d77e24b-e45f-426a-8b7d-dd859271b0d2]")
				return nil
			}

			var result []map[string]interface{}
			today := time.Now()
			for _, salePriceStructListItem := range salePriceStructListItems {
				salePriceModel, err := saleprice.GetSalePriceModel()
				if err != nil {
					logWarnHelper("Can not get sale price model [d5f43503-d73c-4d60-a349-2668ae37c6b0]")
					continue
				}

				err = salePriceModel.Load(salePriceStructListItem.ID)
				if err != nil {
					logWarnHelper("Can not load sale price model [dd08dffe-6147-4d96-8306-c6b60dcb704f]")
					continue
				}

				if isAdmin || (salePriceModel.GetStartDatetime().Before(today) &&
					today.Before(salePriceModel.GetEndDatetime())) {
					result = append(result, salePriceModel.ToHashMap())
				}
			}

			it.SalePrices = result
		}

		return it.SalePrices
	}
	return nil
}

// Set saves product external attributes managed by sale price package
func (it *SalePriceDelegate) Set(attribute string, value interface{}) error {
	switch attribute {
	case "sale_prices":
		// TODO: save sale prices edited on product editing page through sale price model
		if value != nil {
			salePriceCollectionModel, err := saleprice.GetSalePriceCollectionModel()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			salePriceCollectionModel.GetDBCollection().AddFilter("product_id", "=", it.productInstance.GetID())

			salePriceStructListItems, err := salePriceCollectionModel.List()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			for _, salePriceStructListItem := range salePriceStructListItems {
				salePriceModel, err := saleprice.GetSalePriceModel()
				if err != nil {
					return env.ErrorDispatch(err)
				}

				err = salePriceModel.Load(salePriceStructListItem.ID)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				err = salePriceModel.Delete()
				if err != nil {
					return env.ErrorDispatch(err)
				}
			}

			newSalePrices := utils.InterfaceToArray(value)
			for _, salePrice := range newSalePrices {
				salePriceHashMap := utils.InterfaceToMap(salePrice)
				delete(salePriceHashMap, "_id")

				salePriceModel, err := saleprice.GetSalePriceModel()
				if err != nil {
					return env.ErrorDispatch(err)
				}

				err = salePriceModel.FromHashMap(salePriceHashMap)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				err = salePriceModel.Save()
				if err != nil {
					return env.ErrorDispatch(err)
				}

				it.SalePrices = append(it.SalePrices, salePriceModel.ToHashMap())
			}
		}
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
