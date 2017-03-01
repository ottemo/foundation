package category

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
)

// GetEnabled returns enabled flag for the current category
func (it *DefaultCategory) GetEnabled() bool {
	return it.Enabled
}

// GetName returns current category name
func (it *DefaultCategory) GetName() string {
	return it.Name
}

// GetProductIds returns product ids associated to category
func (it *DefaultCategory) GetProductIds() []string {
	return it.ProductIds
}

// GetProductsCollection returns category associated products collection instance
func (it *DefaultCategory) GetProductsCollection() product.InterfaceProductCollection {
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return nil
	}

	dbCollection := productCollection.GetDBCollection()
	if dbCollection != nil {
		if err := dbCollection.AddStaticFilter("_id", "in", it.ProductIds); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "208cb734-953c-4921-add8-71101154d312", err.Error())
		}
	}

	return productCollection
}

// GetProducts returns a set of category associated products
func (it *DefaultCategory) GetProducts() []product.InterfaceProduct {
	var result []product.InterfaceProduct

	for _, productID := range it.ProductIds {
		productModel, err := product.LoadProductByID(productID)
		if err == nil {
			result = append(result, productModel)
		}
	}

	return result
}

// GetParent returns parent category of nil
func (it *DefaultCategory) GetParent() category.InterfaceCategory {
	return it.Parent
}

// GetDescription returns the description of the requested category
func (it *DefaultCategory) GetDescription() string {
	return it.Description
}

// GetImage returns the image of the requested category
func (it *DefaultCategory) GetImage() string {
	return it.Image
}

// AddProduct associates given product with category
func (it *DefaultCategory) AddProduct(productID string) error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "642ed88a-6d8b-48a1-9b3c-feac54c4d9a3", "Can't obtain DBEngine")
	}

	collection, err := dbEngine.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	categoryID := it.GetID()
	if categoryID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "67e7fe19-2ca8-4199-9a7c-94f997d88098", "category ID is not set")
	}
	if productID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e2a7b643-e1b0-46c8-88ad-de2447407875", "product ID is not set")
	}

	if err := collection.AddFilter("category_id", "=", categoryID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0dadcd29-a5cd-4a7c-82e5-4b87279f7e03", err.Error())
	}
	if err := collection.AddFilter("product_id", "=", productID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0d27e020-0515-4b0d-9674-30bd2da5023c", err.Error())
	}
	cnt, err := collection.Count()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if cnt == 0 {
		_, err := collection.Save(map[string]interface{}{"category_id": categoryID, "product_id": productID})
		if err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "623ff72f-6221-4acd-bdf4-e5b765fcd3db", "junction already exists")
	}

	return nil
}

// RemoveProduct un-associates given product with category
func (it *DefaultCategory) RemoveProduct(productID string) error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "92859011-3646-478b-9265-e2fb919e42b3", "Can't obtain DBEngine")
	}

	collection, err := dbEngine.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	categoryID := it.GetID()
	if categoryID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5180a734-0a5e-46ec-9fa2-840a2b1aa6ce", "category ID is not set")
	}
	if productID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "70b5aa6b-dadd-4be8-b8b9-d6f41a7cf237", "product ID is not set")
	}

	if err := collection.AddFilter("category_id", "=", categoryID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b7f797c0-2e71-422c-bbe8-23c1d551a585", err.Error())
	}
	if err := collection.AddFilter("product_id", "=", productID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b43af18d-1bd3-436a-a21e-b29d468a2131", err.Error())
	}
	_, err = collection.Delete()

	return env.ErrorDispatch(err)
}
