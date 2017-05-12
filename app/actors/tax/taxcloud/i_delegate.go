package taxcloud

import (
	"github.com/ottemo/foundation/app/actors/tax/taxcloud/gotaxcloud"
	"github.com/ottemo/foundation/app/actors/tax/taxcloud/model"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceAttributesDelegate implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// New creates new TicDelegate with associated product
func (it *TicDelegate) New(productInstance interface{}) (models.InterfaceAttributesDelegate, error) {
	if productModel, ok := productInstance.(product.InterfaceProduct); ok {
		return &TicDelegate{productInstance: productModel}, nil
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b213a722-3983-481c-b25a-080555fcdb34", "unexpected instance for TIC delegate")
}

// Get returns product external attributes managed by tax cloud package
func (it *TicDelegate) Get(attribute string) interface{} {
	switch attribute {
	case ConstTicIdAttribute:
		if it.productTicPtr != nil {
			return (*it.productTicPtr).GetTicID()
		} else {
			return ConstDefaultTicID
		}
	}
	return nil
}

// Set stores product external attributes managed by tax cloud package
func (it *TicDelegate) Set(attribute string, value interface{}) error {
	if it.productTicPtr == nil {
		productTicModel, err := GetProductTicModel()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		it.productTicPtr = &productTicModel
	}

	switch attribute {
	case ConstTicIdAttribute:
		(*it.productTicPtr).SetTicID(utils.InterfaceToInt(value))
	}

	return nil
}

// GetAttributesInfo describes product external attributes managed by tax cloud package
func (it *TicDelegate) GetAttributesInfo() []models.StructAttributeInfo {
	return []models.StructAttributeInfo{
		{
			Model:      product.ConstModelNameProduct,
			Collection: model.ConstProductTicCollectionName,
			Attribute:  ConstTicIdAttribute,
			Type:       utils.ConstDataTypeInteger,
			Label:      "Taxability Information Code",
			IsRequired: false,
			IsStatic:   false,
			Group:      "General",
			Editors:    "selector",
			Options:    utils.EncodeToJSONString(getTICsMap()),
			Default:    "",
			Validators: "",
		},
	}
}

// Load get tax cloud information for product from db
func (it *TicDelegate) Load(id string) error {
	collection, err := db.GetCollection(model.ConstProductTicCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddFilter("product_id", "=", it.productInstance.GetID()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f3e006e5-aa63-4ce3-8b93-bba679ffb9ca", err.Error())
	}

	productTicModel, err := GetProductTicModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecords, err := collection.Load()
	if err != nil || len(dbRecords) == 0 {
		it.productTicPtr = nil
	} else {
		if err := productTicModel.FromHashMap(dbRecords[0]); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8d3b2df0-de64-4fd7-894a-8bd69b40d60c", err.Error())
		}

		it.productTicPtr = &productTicModel
	}

	return nil
}

// Save stores tax cloud information for product in db
func (it *TicDelegate) Save() error {
	if it.productTicPtr == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "189a9f90-1156-4ccd-ab8e-dcdfd2f5d941", "unable to save empty Taxability Information")
	}

	(*it.productTicPtr).SetProductID(it.productInstance.GetID())
	if err := (*it.productTicPtr).Save(); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// getTICsMap returns list of Taxability Information codes with descriptions
func getTICsMap() map[int]string {
	if ticsCachePtr == nil {
		result := map[int]string{}

		config := env.GetConfig()
		if config == nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "60175012-f3ba-4b4f-a040-4dbeeaaf827c", "can't obtain config")
			// empty
			return result
		}

		gateway := gotaxcloud.NewGateway(
			utils.InterfaceToString(config.GetValue(ConstConfigPathAPILoginID)),
			utils.InterfaceToString(config.GetValue(ConstConfigPathAPIKey)))

		tics, err := gateway.GetTICs()
		if err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "968e6449-52de-49fd-819e-3e0892effe6d", err.Error())
			// empty
			return result
		}

		for _, tic := range tics.TICs {
			result[tic.TICID] = tic.Description
		}

		ticsCachePtr = &result
	}

	return *ticsCachePtr
}

