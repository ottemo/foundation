package product

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// init makes package self-initialization routine
func init() {
	productInstance := new(DefaultProduct)
	var _ product.InterfaceProduct = productInstance
	if err := models.RegisterModel(product.ConstModelNameProduct, productInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1f606bff-caa7-4425-995a-91f6a06d2759", err.Error())
	}

	collectionInstance := new(DefaultProductCollection)
	var _ product.InterfaceProductCollection = collectionInstance
	if err := models.RegisterModel(product.ConstModelNameProductCollection, collectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7878fd53-3cc8-4454-b03c-3e75249079e6", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var shouldFillVisibleField = !collection.HasColumn("visible")

	if err := collection.AddColumn("enabled", db.ConstTypeBoolean, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c06e9370-bc54-4752-9f68-ea6656d69230", err.Error())
	}
	if err := collection.AddColumn("sku", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "742f35ff-059c-498f-a441-2ecb7901394b", err.Error())
	}
	if err := collection.AddColumn("name", db.ConstTypeVarchar, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a7d441c2-106e-496d-a617-4390556ed17c", err.Error())
	}
	if err := collection.AddColumn("short_description", db.ConstTypeVarchar, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d4ceda3b-b1c4-46b6-90f5-5e060045f766", err.Error())
	}
	if err := collection.AddColumn("description", db.ConstTypeText, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d33824eb-5731-4737-b417-74de7d9c9f68", err.Error())
	}
	if err := collection.AddColumn("default_image", db.ConstTypeVarchar, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a3a9e157-d1d4-4dd7-9eaf-b081ce1b357a", err.Error())
	}
	if err := collection.AddColumn("price", db.ConstTypeMoney, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1e50bacb-9ab5-4ae2-8065-78476745c244", err.Error())
	}
	if err := collection.AddColumn("weight", db.ConstTypeFloat, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0c772222-0464-4682-ac3e-81e0ba7f0045", err.Error())
	}
	if err := collection.AddColumn("options", db.ConstTypeJSON, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "47ae696b-e798-4d3c-a423-f76bb7d7e7c8", err.Error())
	}
	if err := collection.AddColumn("visible", db.ConstTypeBoolean, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "55badc97-99e4-4162-8683-e892295b37d7", err.Error())
	}

	if err := collection.AddColumn("related_pids", db.TypeArrayOf(db.ConstTypeID), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0b63db43-4cb0-4f9e-85f6-d8850dadb4c9", err.Error())
	}

	if shouldFillVisibleField {
		env.Log(ConstErrorModule, env.ConstLogPrefixInfo, "Field 'visible' have been added. Make all products visible.")
		if err:= fillVisibleField(); err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		env.Log(ConstErrorModule, env.ConstLogPrefixInfo, "'visible' value need no update.")
	}

	return nil
}

// fillVisibleField makes all product visible on adding "visible" column
func fillVisibleField() error {
	// get product collection
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// update products
	for _, currentProduct := range productCollection.ListProducts() {
		err = currentProduct.Set("visible", true)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err := currentProduct.Save()
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}
