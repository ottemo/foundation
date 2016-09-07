package saleprice

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/utils"
	contextPkg "github.com/ottemo/foundation/api/context"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceAttributesDelegate implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

func (it *SalePriceDelegate) New(productInstance interface{}) (models.InterfaceAttributesDelegate, error) {
	if productModel, ok := productInstance.(product.InterfaceProduct); ok {
		return &SalePriceDelegate{productInstance: productModel}, nil
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6ac6965a-1f1e-44ae-b854-ad430d5b85a6", "unexpected instance for sale price delegate")
}

func (it *SalePriceDelegate) Get(attribute string) interface{} {
	context := contextPkg.GetContext()
	logDebugHelper("SalePriceDelegate Load "+utils.InterfaceToString(context))
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
				return nil // TODO unhandled error
			}

			salePriceCollectionModel.GetDBCollection().AddFilter("product_id", "=", it.productInstance.GetID())

			salePriceStructListItems, err := salePriceCollectionModel.List()
			if err != nil {
				return nil // TODO unhandled error
			}

			var result []map[string]interface{}
			today := time.Now()
			for _, salePriceStructListItem := range salePriceStructListItems {
				salePriceModel, err := saleprice.GetSalePriceModel()
				if err != nil {
					continue // TODO unhandled error
				}

				err = salePriceModel.Load(salePriceStructListItem.ID)
				if err != nil {
					continue // TODO unhandled error
				}

				if isAdmin || (
					salePriceModel.GetStartDatetime().Before(today) &&
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

func (it *SalePriceDelegate) Set(attribute string, value interface{}) error {
	switch attribute {
	case "sale_prices":
		// TODO save sale prices edited on product editing page through sale price model
	}

	return nil
}

func (it *SalePriceDelegate) GetAttributesInfo() []models.StructAttributeInfo {
	return []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameSalePrices,
			Attribute:  "sale_prices",
			Type:       utils.ConstDataTypeJSON,
			Label:      "SalePrices",
			IsRequired: false,
			IsStatic:   false,
			Group:      "General",
			Editors:    "product_sale_prices", // TODO editor not implemented yet
			Options:    "",
			Default:    "",
			Validators: "",
		},
	}
}
