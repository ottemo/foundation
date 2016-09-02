package saleprice

import (
	"strings"

	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
)

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

func (it *DefaultSalePrice) GetModelName() string {
	return ConstModelNameSalePrice
}

func (it *DefaultSalePrice) GetImplementationName() string {
	return "Default" + ConstModelNameSalePrice
}

func (it *DefaultSalePrice) New() (models.InterfaceModel, error) {
	return &DefaultSalePrice{}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

func (it *DefaultSalePrice) Get(attribute string) interface{} {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id":
		return it.GetID()

	case "amount":
		return it.GetAmount()

	case "end_datetime":
		return it.GetEndDatetime()

	case "product_id":
		return it.GetProductID()

	case "start_datetime":
		return it.GetStartDatetime()
	}

	return nil

}

// Set sets attribute value to object or returns error
func (it *DefaultSalePrice) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id":
		it.SetID(utils.InterfaceToString(value))

	case "amount":
		it.SetAmount(utils.InterfaceToFloat64(value))

	case "end_datetime":
		it.SetEndDatetime(utils.InterfaceToTime(value))

	case "product_id":
		it.SetProductID(utils.InterfaceToString(value))

	case "start_datetime":
		it.SetStartDatetime(utils.InterfaceToTime(value))
	}

	return nil
}

func (it *DefaultSalePrice) FromHashMap(input map[string]interface{}) error {
	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

func (it *DefaultSalePrice) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.GetID()
	result["amount"] = it.GetAmount()
	result["end_datetime"] = it.GetEndDatetime()
	result["product_id"] = it.GetProductID()
	result["start_datetime"] = it.GetStartDatetime()

	return result
}

func (it *DefaultSalePrice) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      saleprice.ConstModelNameSalePrice,
			Collection: ConstCollectionNameSalePrices,
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
			Model:      saleprice.ConstModelNameSalePrice,
			Collection: ConstCollectionNameSalePrices,
			Attribute:  "amount",
			Type:       db.ConstTypeMoney,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Amount",
			Group:      "General",
			Editors:    "money",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      saleprice.ConstModelNameSalePrice,
			Collection: ConstCollectionNameSalePrices,
			Attribute:  "end_datetime",
			Type:       db.ConstTypeDatetime,
			IsRequired: true,
			IsStatic:   true,
			Label:      "End Datetime",
			Group:      "General",
			Editors:    "datetime",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      saleprice.ConstModelNameSalePrice,
			Collection: ConstCollectionNameSalePrices,
			Attribute:  "product_id",
			Type:       db.ConstTypeID,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Product ID",
			Group:      "General",
			Editors:    "product_selector",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      saleprice.ConstModelNameSalePrice,
			Collection: ConstCollectionNameSalePrices,
			Attribute:  "start_datetime",
			Type:       db.ConstTypeDatetime,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Start Datetime",
			Group:      "General",
			Editors:    "datetime",
			Options:    "",
			Default:    "",
		},
	}

	return info
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceStorable implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// SetID sets database storage id for current object
func (it *DefaultSalePrice) SetID(id string) error {
	it.id = id
	return nil
}

// GetID returns database storage id of current object
func (it *DefaultSalePrice) GetID() string {
	return it.id
}

// Save function check model and save it to storage
func (it *DefaultSalePrice) Save() error {
	logDebugHelper("(it *DefaultSalePrice) Save " + utils.InterfaceToString(it.ToHashMap()))
	// Check model data
	//-----------------

	// Check amount positive
	if it.GetAmount() <= 0 {
		return newErrorHelper("Amount should be greater than 0.", "ccf50f3f-a503-4720-b3a6-2ba1639fb8e7")
	}

	// Check start date before end date
	if !it.GetStartDatetime().Before(it.GetEndDatetime()) {
		return newErrorHelper("Start Datetime should be before End Datetime.", "668c3bd4-1a10-417a-aa68-2ec13e559a11")
	}

	// Check product exists
	productModel, err := product.LoadProductByID(it.GetProductID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Check amount < product price
	if it.GetAmount() < productModel.GetPrice() {
		return newErrorHelper("Amount should be less than product price.", "e30a767c-08a3-484f-9453-106290e99050")
	}

	//TODO: check period is not overlapped with other periods for product if exists

	// Save model to storage
	//----------------------
	dbCollection, err := db.GetCollection(ConstCollectionNameSalePrices)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	newID, err := dbCollection.Save(it.ToHashMap())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.SetID(newID)

	return nil
}

func (it *DefaultSalePrice) Load(id string) error {
	dbSalePriceCollection, err := db.GetCollection(ConstCollectionNameSalePrices)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := dbSalePriceCollection.LoadByID(id)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func (it *DefaultSalePrice) Delete() error {
	dbCollection, err := db.GetCollection(ConstCollectionNameSalePrices)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

