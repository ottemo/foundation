package product

import (
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
)

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
		productModel.FromHashMap(dbRecordData)

		result = append(result, productModel)
	}

	return result
}
