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
