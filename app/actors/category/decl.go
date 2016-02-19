// Package category is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/category" package
package category

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Package global constants
const (
	ConstCollectionNameCategory                = "category"
	ConstCollectionNameCategoryProductJunction = "category_product"

	ConstErrorModule = "category"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstCategoryMediaTypeImage = "image"
)

// DefaultCategory is a default implementer of InterfaceCategory
type DefaultCategory struct {
	id string

	Enabled     bool
	Name        string
	Description string
	Image       string
	Parent      category.InterfaceCategory
	Path        string
	ProductIds  []string
}

// DefaultCategoryCollection is a default implementer of InterfaceCategoryCollection
type DefaultCategoryCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}

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
