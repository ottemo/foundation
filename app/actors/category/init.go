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
	if err := models.RegisterModel(category.ConstModelNameCategory, categoryInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5cd7a0d2-007b-4a2d-8559-2f3d1df4e441", err.Error())
	}

	categoryCollectionInstance := new(DefaultCategoryCollection)
	var _ category.InterfaceCategoryCollection = categoryCollectionInstance
	if err := models.RegisterModel(category.ConstModelNameCategoryCollection, categoryCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b0664a4e-c269-4987-98d2-e6957827d338", err.Error())
	}

	db.RegisterOnDatabaseStart(categoryInstance.setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func (it *DefaultCategory) setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("enabled", db.ConstTypeBoolean, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "cc1a6ba4-d16f-4444-9e0c-d946667b3915", err.Error())
	}
	if err := collection.AddColumn("parent_id", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "06458f1e-fc62-49d9-86fd-39e900c2df9a", err.Error())
	}
	if err := collection.AddColumn("path", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "12d1697a-8ead-42ec-a59e-c042b62c1458", err.Error())
	}
	if err := collection.AddColumn("name", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6a2d91ad-71fc-4dcb-b3c7-779f98e3c935", err.Error())
	}
	if err := collection.AddColumn("description", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4dcaa63d-19dd-490a-b8b4-a87536cdc2dc", err.Error())
	}
	if err := collection.AddColumn("image", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f4f83dc6-7eae-45ab-aeac-9b3edc6decb4", err.Error())
	}

	collection, err = db.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("category_id", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "79063fb7-f828-4576-9b90-80047e67e010", err.Error())
	}
	if err := collection.AddColumn("product_id", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "73dd8523-649f-486e-9f02-d329194c14b7", err.Error())
	}

	return nil
}
