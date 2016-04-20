package category

import (
	"strings"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
)

// GetModelName returns model name
func (it *DefaultCategory) GetModelName() string {
	return category.ConstModelNameCategory
}

// GetImplementationName returns model implementation name
func (it *DefaultCategory) GetImplementationName() string {
	return "Default" + category.ConstModelNameCategory
}

// New returns new instance of model implementation object
func (it *DefaultCategory) New() (models.InterfaceModel, error) {
	return &DefaultCategory{ProductIds: make([]string, 0)}, nil
}

// GetModelName returns model name
func (it *DefaultCategoryCollection) GetModelName() string {
	return category.ConstModelNameCategory
}

// GetImplementationName returns model implementation name
func (it *DefaultCategoryCollection) GetImplementationName() string {
	return "Default" + category.ConstModelNameCategory
}

// New returns new instance of model implementation object
func (it *DefaultCategoryCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultCategoryCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}

////////////////////////////////////////////////////////
////////////////////////////////////////////////////////
//////// Getters and Setters
////////////////////////////////////////////////////////
////////////////////////////////////////////////////////

// GetID returns database storage id of current object
func (it *DefaultCategory) GetID() string {
	return it.id
}

// SetID sets database storage id for current object
func (it *DefaultCategory) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// GetEnabled returns enabled flag for the current category
func (it *DefaultCategory) GetEnabled() bool {
	return it.Enabled
}

// GetName returns current category name
func (it *DefaultCategory) GetName() string {
	return it.Name
}

// GetImage returns the image of the requested category
func (it *DefaultCategory) GetImage() string {
	return it.Image
}

// GetProductIds returns product ids associated to category
func (it *DefaultCategory) GetProductIds() []string {
	return it.ProductIds
}

// GetParent returns parent category of nil
func (it *DefaultCategory) GetParent() category.InterfaceCategory {
	return it.Parent
}

// GetDescription returns the description of the requested category
func (it *DefaultCategory) GetDescription() string {
	return it.Description
}

// GetAttributesInfo returns information about object attributes
func (it *DefaultCategory) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      category.ConstModelNameCategory,
			Collection: ConstCollectionNameCategory,
			Attribute:  "_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      category.ConstModelNameCategory,
			Collection: ConstCollectionNameCategory,
			Attribute:  "enabled",
			Type:       db.ConstTypeBoolean,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Enabled",
			Group:      "General",
			Editors:    "boolean",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      category.ConstModelNameCategory,
			Collection: ConstCollectionNameCategory,
			Attribute:  "name",
			Type:       db.ConstTypeText,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      category.ConstModelNameCategory,
			Collection: ConstCollectionNameCategory,
			Attribute:  "parent_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Parent",
			Group:      "General",
			Editors:    "category_selector",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      category.ConstModelNameCategory,
			Collection: ConstCollectionNameCategory,
			Attribute:  "description",
			Type:       db.ConstTypeText,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Description",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      category.ConstModelNameCategory,
			Collection: ConstCollectionNameCategory,
			Attribute:  "image",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Image",
			Group:      "General",
			Editors:    "image_selector",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      category.ConstModelNameCategory,
			Collection: ConstCollectionNameCategory,
			Attribute:  "products",
			Type:       db.TypeArrayOf(db.ConstTypeID),
			IsRequired: false,
			IsStatic:   true,
			Label:      "Products",
			Group:      "General",
			Editors:    "product_selector",
			Options:    "",
			Default:    "",
		},
	}

	return info
}

// Get returns object attribute value or nil
func (it *DefaultCategory) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.GetID()

	case "enabled":
		return it.GetEnabled()

	case "name":
		return it.GetName()

	case "path":
		if it.Path == "" {
			it.updatePath()
		}
		return it.Path

	case "parent_id":
		if it.Parent != nil {
			return it.Parent.GetID()
		}
		return ""

	case "parent":
		return it.GetParent()

	case "image":
		return it.GetImage()

	case "description":
		return it.GetDescription()

	case "product_ids":
		return it.GetProductIds()
	}

	return nil
}

// Set sets attribute value to object or returns error
func (it *DefaultCategory) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.SetID(utils.InterfaceToString(value))

	case "enabled":
		it.Enabled = utils.InterfaceToBool(value)

	case "name":
		it.Name = utils.InterfaceToString(value)

	case "parent_id":
		if value, ok := value.(string); ok {
			value = strings.TrimSpace(value)
			if value != "" {
				model, err := models.GetModel("Category")
				if err != nil {
					return env.ErrorDispatch(err)
				}
				categoryModel, ok := model.(category.InterfaceCategory)
				if !ok {
					return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "39b6496a-4145-4b16-9f67-ca6375fd8b1f", "unsupported category model "+model.GetImplementationName())
				}

				err = categoryModel.Load(value)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				selfID := it.GetID()
				if selfID != "" {
					parentPath, ok := categoryModel.Get("path").(string)
					if categoryModel.GetID() != selfID && ok && !strings.Contains(parentPath, selfID) {
						it.Parent = categoryModel
					} else {
						return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0ae64841-1123-4add-8250-c4f324ad8eab", "category can't have sub-category or itself as parent")
					}
				} else {
					it.Parent = categoryModel
				}
			} else {
				it.Parent = nil
			}
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "04ac194b-1912-4520-9087-b0248b9ea758", "unsupported id specified")
		}
		it.updatePath()

	case "parent":
		switch value := value.(type) {
		case category.InterfaceCategory:
			it.Parent = value
		case string:
			it.Set("parent_id", value)
		default:
			env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2156d563-932b-4de7-a615-7d473717a3bd", "unsupported 'parent' value")
		}
		// path should be changed as well
		it.updatePath()

	case "image":
		it.Image = utils.InterfaceToString(value)

	case "description":
		it.Description = utils.InterfaceToString(value)

	case "products":
		switch typedValue := value.(type) {

		case []interface{}:
			for _, listItem := range typedValue {
				productID, ok := listItem.(string)
				if ok {
					productModel, err := product.LoadProductByID(productID)
					if err != nil {
						return env.ErrorDispatch(err)
					}

					it.ProductIds = append(it.ProductIds, productModel.GetID())
				}
			}

		case []product.InterfaceProduct:
			it.ProductIds = make([]string, 0)
			for _, productItem := range typedValue {
				it.ProductIds = append(it.ProductIds, productItem.GetID())
			}

		default:
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "84284b03-0a29-4036-aa2d-b35768884b63", "unsupported 'products' value")
		}
	}
	return nil
}

// FromHashMap fills object attributes from map[string]interface{}
func (it *DefaultCategory) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.LogError(err)
		}
	}

	return nil
}

// ToHashMap represents object as map[string]interface{}
func (it *DefaultCategory) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.GetID()

	result["enabled"] = it.Get("enabled")
	result["description"] = it.Get("description")

	result["image"] = it.Get("image")

	result["parent_id"] = it.Get("parent_id")
	result["name"] = it.Get("name")
	result["product_ids"] = it.Get("product_ids")
	result["path"] = it.Get("path")

	return result
}

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

// AddMedia adds new media assigned to category
func (it *DefaultCategory) AddMedia(mediaType string, mediaName string, content []byte) error {
	categoryID := it.GetID()
	if categoryID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "85650715-3acf-4e47-a365-c6e8911d9118", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Save(it.GetModelName(), categoryID, mediaType, mediaName, content)
}

// RemoveMedia removes media assigned to category
func (it *DefaultCategory) RemoveMedia(mediaType string, mediaName string) error {
	categoryID := it.GetID()
	if categoryID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "87bb383a-cf35-48e0-9d50-ad517ed2e8f9", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Remove(it.GetModelName(), categoryID, mediaType, mediaName)
}

// ListMedia lists media assigned to category
func (it *DefaultCategory) ListMedia(mediaType string) ([]string, error) {
	var result []string

	categoryID := it.GetID()
	if categoryID == "" {
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "45b1ebde-3dd0-4c6c-9960-fddd89f4907f", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	return mediaStorage.ListMedia(it.GetModelName(), categoryID, mediaType)
}

// GetMedia returns content of media assigned to category
func (it *DefaultCategory) GetMedia(mediaType string, mediaName string) ([]byte, error) {
	categoryID := it.GetID()
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5f5d3c33-de82-4580-a6e7-f5c45e9281e5", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaStorage.Load(it.GetModelName(), categoryID, mediaType, mediaName)
}

// GetMediaPath returns relative location of media assigned to category in media storage
func (it *DefaultCategory) GetMediaPath(mediaType string) (string, error) {
	categoryID := it.GetID()
	if categoryID == "" {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0055f93a-5d10-41db-8d93-ea2bb4bee216", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaStorage.GetMediaPath(it.GetModelName(), categoryID, mediaType)
}

// Load loads object information from database storage
func (it *DefaultCategory) Load(ID string) error {

	// loading category
	categoryCollection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := categoryCollection.LoadByID(ID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.updatePath()

	// loading category product ids
	junctionCollection, err := db.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	junctionCollection.AddFilter("category_id", "=", it.GetID())
	junctedProducts, err := junctionCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, junctionRecord := range junctedProducts {
		it.ProductIds = append(it.ProductIds, utils.InterfaceToString(junctionRecord["product_id"]))
	}

	return nil
}

// Delete removes current object from database storage
func (it *DefaultCategory) Delete() error {
	//deleting category products join
	junctionCollection, err := db.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = junctionCollection.AddFilter("category_id", "=", it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	_, err = junctionCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// deleting category
	categoryCollection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = categoryCollection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Save stores current object to database storage
func (it *DefaultCategory) Save() error {

	storingValues := it.ToHashMap()

	delete(storingValues, "products")

	categoryCollection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// saving category
	if newID, err := categoryCollection.Save(storingValues); err == nil {
		if it.GetID() != newID {
			it.SetID(newID)
			it.updatePath()
			it.Save()
		}
	} else {
		return env.ErrorDispatch(err)
	}

	// saving category products assignment
	junctionCollection, err := db.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// deleting old assigned products
	junctionCollection.AddFilter("category_id", "=", it.GetID())
	_, err = junctionCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// adding new assignments
	for _, categoryProductID := range it.ProductIds {
		junctionCollection.Save(map[string]interface{}{"category_id": it.GetID(), "product_id": categoryProductID})
	}

	return nil
}
