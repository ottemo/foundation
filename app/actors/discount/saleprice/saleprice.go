package saleprice

// DefaultSalePrice type implements:
// 	- InterfaceSalePrice
// 	- InterfaceModel
// 	- InterfaceObject
// 	- InterfaceStorable

import (
	"strings"
	"time"

	"github.com/ottemo/foundation/app/models"
	salepriceModel "github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultSalePrice) GetModelName() string {
	return salepriceModel.ConstModelNameSalePrice
}

// GetImplementationName returns default model implementation name
func (it *DefaultSalePrice) GetImplementationName() string {
	return "Default" + salepriceModel.ConstModelNameSalePrice
}

// New creates new model
func (it *DefaultSalePrice) New() (models.InterfaceModel, error) {
	return &DefaultSalePrice{}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// Get return model attribute by name
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

// FromHashMap converts object represented by hash map to object
func (it *DefaultSalePrice) FromHashMap(input map[string]interface{}) error {
	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap converts object data to hash map presentation
func (it *DefaultSalePrice) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.GetID()
	result["amount"] = it.GetAmount()
	result["end_datetime"] = it.GetEndDatetime()
	result["product_id"] = it.GetProductID()
	result["start_datetime"] = it.GetStartDatetime()

	return result
}

// GetAttributesInfo describes model attributes
func (it *DefaultSalePrice) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      salepriceModel.ConstModelNameSalePrice,
			Collection: salepriceModel.ConstSalePriceDbCollectionName,
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
			Model:      salepriceModel.ConstModelNameSalePrice,
			Collection: salepriceModel.ConstSalePriceDbCollectionName,
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
			Model:      salepriceModel.ConstModelNameSalePrice,
			Collection: salepriceModel.ConstSalePriceDbCollectionName,
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
			Model:      salepriceModel.ConstModelNameSalePrice,
			Collection: salepriceModel.ConstSalePriceDbCollectionName,
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
			Model:      salepriceModel.ConstModelNameSalePrice,
			Collection: salepriceModel.ConstSalePriceDbCollectionName,
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
	// Check model data
	//-----------------

	// Truncate datetimes by hour
	it.SetStartDatetime(it.GetStartDatetime().Truncate(time.Hour))
	it.SetEndDatetime(it.GetEndDatetime().Truncate(time.Hour))

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
	if it.GetAmount() >= productModel.GetPrice() {
		return newErrorHelper("Amount should be less than product price.", "e30a767c-08a3-484f-9453-106290e99050")
	}

	// Check period is not overlapped with other periods for product if exists
	salePriceCollection, err := db.GetCollection(salepriceModel.ConstSalePriceDbCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = salePriceCollection.AddFilter("product_id", "=", it.GetProductID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = salePriceCollection.AddFilter("start_datetime", "<=", it.GetEndDatetime())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = salePriceCollection.AddFilter("end_datetime", ">=", it.GetStartDatetime())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	salePrices, err := salePriceCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}
	if len(salePrices) > 0 {
		return newErrorHelper(
			"New datetime range is overlapped with other range already specified for product.",
			"ae75884a-b239-4739-a8cf-fa9e1d9d5ee1")
	}

	// Save model to storage
	//----------------------
	err = salePriceCollection.ClearFilters()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	newID, err := salePriceCollection.Save(it.ToHashMap())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.SetID(newID)

	return nil
}

// Load loads model from storage
func (it *DefaultSalePrice) Load(id string) error {
	dbSalePriceCollection, err := db.GetCollection(salepriceModel.ConstSalePriceDbCollectionName)
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

// Delete deletes model from storage
func (it *DefaultSalePrice) Delete() error {
	dbCollection, err := db.GetCollection(salepriceModel.ConstSalePriceDbCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceSalePrice implementation (package "github.com/ottemo/foundation/app/models/discount/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// SetAmount : amount setter
func (it *DefaultSalePrice) SetAmount(amount float64) error {
	it.amount = amount
	return nil
}

// GetAmount : amount getter
func (it *DefaultSalePrice) GetAmount() float64 {
	return it.amount
}

// SetEndDatetime : endDatetime setter
func (it *DefaultSalePrice) SetEndDatetime(endDatetime time.Time) error {
	it.endDatetime = endDatetime
	return nil
}

// GetEndDatetime : endDatetime getter
func (it *DefaultSalePrice) GetEndDatetime() time.Time {
	return it.endDatetime
}

// SetProductID : productID setter
func (it *DefaultSalePrice) SetProductID(productID string) error {
	it.productID = productID
	return nil
}

// GetProductID : productID getter
func (it *DefaultSalePrice) GetProductID() string {
	return it.productID
}

// SetStartDatetime : startDatetime setter
func (it *DefaultSalePrice) SetStartDatetime(startDatetime time.Time) error {
	it.startDatetime = startDatetime
	return nil
}

// GetStartDatetime : startDatetime getter
func (it *DefaultSalePrice) GetStartDatetime() time.Time {
	return it.startDatetime
}
