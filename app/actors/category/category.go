package category

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
)

// GetProductsCollection returns category associated products collection instance
func (it *DefaultCategory) GetProductsCollection() product.InterfaceProductCollection {
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return nil
	}

	dbCollection := productCollection.GetDBCollection()
	if dbCollection != nil {
		dbCollection.AddStaticFilter("_id", "in", it.ProductIds)
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

	collection.AddFilter("category_id", "=", categoryID)
	collection.AddFilter("product_id", "=", productID)
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

	collection.AddFilter("category_id", "=", categoryID)
	collection.AddFilter("product_id", "=", productID)
	_, err = collection.Delete()

	return env.ErrorDispatch(err)
}

// GetCollection returns collection of current instance type
func (it *DefaultCategory) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(category.ConstModelNameCategoryCollection)
	if result, ok := model.(category.InterfaceCategoryCollection); ok {
		return result
	}

	return nil
}

// GetDBCollection returns database collection
func (it *DefaultCategoryCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListCategories returns list of category model items
func (it *DefaultCategoryCollection) ListCategories() []category.InterfaceCategory {
	var result []category.InterfaceCategory

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		categoryModel, err := category.GetCategoryModel()
		if err != nil {
			return result
		}
		categoryModel.FromHashMap(recordData)

		result = append(result, categoryModel)
	}

	return result
}

// List enumerates items of model type
func (it *DefaultCategoryCollection) List() ([]models.StructListItem, error) {
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
		categoryModel, err := category.GetCategoryModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		categoryModel.FromHashMap(dbItemData)

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		mediaPath, err := categoryModel.GetMediaPath("image")
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		resultItem.ID = categoryModel.GetID()
		resultItem.Name = categoryModel.GetName()
		resultItem.Image = ""
		resultItem.Desc = categoryModel.GetDescription()

		if categoryModel.GetImage() != "" {
			resultItem.Image = mediaPath + categoryModel.GetImage()
		}

		// serving extra attributes
		//-------------------------
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = categoryModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute allows to obtain additional attributes from  List() function
func (it *DefaultCategoryCollection) ListAddExtraAttribute(attribute string) error {

	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var allowedAttributes []string
	for _, attributeInfo := range categoryModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}
	allowedAttributes = append(allowedAttributes, "parent")

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2509d847-ba1e-48bd-9b29-37edd0cac52b", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3282704a-a048-4de6-b910-b23c753083a9", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds selection filter to List() function
func (it *DefaultCategoryCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// ListFilterReset clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultCategoryCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// ListLimit specifies selection paging
func (it *DefaultCategoryCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}

// updatePath is an internal function used to update "path" attribute of object
func (it *DefaultCategory) updatePath() {
	if it.GetID() == "" {
		it.Path = ""
	} else if it.Parent != nil {
		parentPath, ok := it.Parent.Get("path").(string)
		if ok {
			it.Path = parentPath + "/" + it.GetID()
		}
	} else {
		it.Path = "/" + it.GetID()
	}
}
