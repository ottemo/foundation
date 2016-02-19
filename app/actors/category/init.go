package category

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	categoryInstance := new(DefaultCategory)
	var _ category.InterfaceCategory = categoryInstance
	models.RegisterModel(category.ConstModelNameCategory, categoryInstance)

	categoryCollectionInstance := new(DefaultCategoryCollection)
	var _ category.InterfaceCategoryCollection = categoryCollectionInstance
	models.RegisterModel(category.ConstModelNameCategoryCollection, categoryCollectionInstance)

	db.RegisterOnDatabaseStart(categoryInstance.setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func (it *DefaultCategory) setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("enabled", db.ConstTypeBoolean, true)
	collection.AddColumn("parent_id", db.ConstTypeID, true)
	collection.AddColumn("path", db.ConstTypeVarchar, true)
	collection.AddColumn("name", db.ConstTypeVarchar, true)
	collection.AddColumn("description", db.ConstTypeVarchar, true)
	collection.AddColumn("image", db.ConstTypeVarchar, true)

	collection, err = db.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("category_id", db.ConstTypeID, true)
	collection.AddColumn("product_id", db.ConstTypeID, true)

	return nil
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
