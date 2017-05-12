package taxcloud

import (
	"strings"

	"github.com/ottemo/foundation/app/actors/tax/taxcloud/model"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetProductTicModel returns default implementation of productTicModel
func GetProductTicModel() (model.InterfaceProductTic, error) {
	foundModel, err := models.GetModel(model.ConstProductTicModelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	productTicModel, ok := foundModel.(model.InterfaceProductTic)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1a37c6c9-4a54-4d60-8566-1a84beab95f4", "model "+foundModel.GetImplementationName()+" is not 'InterfaceProductTic' capable")
	}

	return productTicModel, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultProductTic) GetModelName() string {
	return model.ConstProductTicModelName
}

// GetImplementationName returns default model implementation name
func (it *DefaultProductTic) GetImplementationName() string {
	return "Default" + model.ConstProductTicModelName
}

// New creates new model
func (it *DefaultProductTic) New() (models.InterfaceModel, error) {
	return &DefaultProductTic{}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// Get return model attribute by name
func (it *DefaultProductTic) Get(attribute string) interface{} {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id":
		return it.GetID()

	case ConstTicIdAttribute:
		return it.GetTicID()

	case "product_id":
		return it.GetProductID()
	}

	return nil
}

// Set sets attribute value to object or returns error
func (it *DefaultProductTic) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id":
		if err := it.SetID(utils.InterfaceToString(value)); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bb61a3d6-445a-47eb-bad5-38f693269fb5", err.Error())
		}

	case ConstTicIdAttribute:
		if err := it.SetTicID(utils.InterfaceToInt(value)); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "84b58094-c810-482c-b159-7c9dc4641035", err.Error())
		}

	case "product_id":
		if err := it.SetProductID(utils.InterfaceToString(value)); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ec54f364-9fff-4572-b57c-49bb140d9b62", err.Error())
		}
	}

	return nil
}

// FromHashMap converts object represented by hash map to object
func (it *DefaultProductTic) FromHashMap(input map[string]interface{}) error {
	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap converts object data to hash map presentation
func (it *DefaultProductTic) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.GetID()
	result[ConstTicIdAttribute] = it.GetTicID()
	result["product_id"] = it.GetProductID()

	return result
}

// GetAttributesInfo describes model attributes
func (it *DefaultProductTic) GetAttributesInfo() []models.StructAttributeInfo {
	return []models.StructAttributeInfo{
		{
			Model:      model.ConstProductTicModelName,
			Collection: model.ConstProductTicCollectionName,
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
		{
			Model:      model.ConstProductTicModelName,
			Collection: model.ConstProductTicCollectionName,
			Attribute:  ConstTicIdAttribute,
			Type:       utils.ConstDataTypeInteger,
			Label:      "Taxability Information Code",
			IsRequired: false,
			IsStatic:   false,
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "",
		},
		{
			Model:      model.ConstProductTicModelName,
			Collection: model.ConstProductTicCollectionName,
			Attribute:  "product_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Product",
			Group:      "General",
			Editors:    "product_selector",
			Options:    "",
			Default:    "",
		},
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceStorable implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// SetID sets database storage id for current object
func (it *DefaultProductTic) SetID(id string) error {
	it.id = id
	return nil
}

// GetID returns database storage id of current object
func (it *DefaultProductTic) GetID() string {
	return it.id
}

// Save function check model and save it to storage
func (it *DefaultProductTic) Save() error {
	productTicCollection, err := db.GetCollection(model.ConstProductTicCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	newID, err := productTicCollection.Save(it.ToHashMap())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := it.SetID(newID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "18f29f77-15fb-4ac2-8152-050f54881ef6", err.Error())
	}

	return nil
}

// Load loads model from storage
func (it *DefaultProductTic) Load(id string) error {
	dbProductTicCollection, err := db.GetCollection(model.ConstProductTicCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := dbProductTicCollection.LoadByID(id)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Delete deletes model from storage
func (it *DefaultProductTic) Delete() error {
	dbCollection, err := db.GetCollection(model.ConstProductTicCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// SetProductID : productID setter
func (it *DefaultProductTic) SetProductID(productID string) error {
	it.productID = productID
	return nil
}

// GetProductID : productID getter
func (it *DefaultProductTic) GetProductID() string {
	return it.productID
}

// SetTicID : ticID setter
func (it *DefaultProductTic) SetTicID(ticID int) error {
	it.ticID = ticID
	return nil
}

// GetTicID : ticID getter
func (it *DefaultProductTic) GetTicID() int {
	return it.ticID
}
