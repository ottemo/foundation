package stock

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/stock"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// New creates new model
func (it *DefaultStock) New() (models.InterfaceModel, error) {
	return &DefaultStock{}, nil
}

// GetModelName returns model name
func (it *DefaultStock) GetModelName() string {
	return stock.ConstModelNameStock
}

// GetImplementationName returns default model implementation name
func (it *DefaultStock) GetImplementationName() string {
	return "Default" + stock.ConstModelNameStock
}

// newErrorHelper produce new module level error is declared to minimize repeatable code
func newErrorHelper(msg, code string) error {
	return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, code, msg)
}

//--------------------------------------------------------------------------------------------------------
// InterfaceObject implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// Get return model attribute by name
func (it *DefaultStock) Get(attribute string) interface{} {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id":
		return it.GetID()
	case "options":
		return it.GetOptions()
	case "qty":
		return it.GetQty()
	case "product_id":
		return it.GetProductID()
	}

	return nil

}

// Set sets attribute value to object or returns error
func (it *DefaultStock) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
		case "_id":
			if err := it.SetID(utils.InterfaceToString(value)); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5d2016b5-1663-4ba0-9085-361410d0f4ff", err.Error())
			}
		case "options":
			if err := it.SetOptions(utils.InterfaceToString(value)); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5d2016b5-1663-4ba0-9085-361410d0f4ff", err.Error())
			}
		case "product_id":
			if err := it.SetProductID(utils.InterfaceToString(value)); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5d2016b5-1663-4ba0-9085-361410d0f4ff", err.Error())
			}
		case "qty":
			if err := it.SetQty(utils.InterfaceToInt(value)); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5d2016b5-1663-4ba0-9085-361410d0f4ff", err.Error())
			}
	}

	return nil
}

// FromHashMap converts object represented by hash map to object
func (it *DefaultStock) FromHashMap(input map[string]interface{}) error {
	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap converts object data to hash map presentation
func (it *DefaultStock) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.GetID()
	result["options"] = it.GetOptions()
	result["product_id"] = it.GetProductID()
	result["qty"] = it.GetQty()

	return result
}

// GetAttributesInfo describes model attributes
func (it *DefaultStock) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      stock.ConstModelNameStock,
			Collection: stock.ConstModelNameStockCollection,
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
			Model:      stock.ConstModelNameStock,
			Collection: stock.ConstModelNameStockCollection,
			Attribute:  "product_id",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Product ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stock.ConstModelNameStock,
			Collection: stock.ConstModelNameStockCollection,
			Attribute:  "options",
			Type:       db.ConstTypeJSON,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Options",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stock.ConstModelNameStock,
			Collection: stock.ConstModelNameStockCollection,
			Attribute:  "qty",
			Type:       db.ConstTypeInteger,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Qty",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
	}

	return info
}

// Save function check model and save it to storage
func (it *DefaultStock) Save() error {
	// Check model data
	//-----------------

	// Check amount positive
	if it.qty < 0 {
		return newErrorHelper("Qty should be greater not less than 0.", "ccf50f3f-a503-4720-b3a6-2ba1639fb8e7")
	}

	// Check product exists
	_, err := product.LoadProductByID(it.product_id)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Save model to storage
	//----------------------
	stockCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	newID, err := stockCollection.Save(it.ToHashMap())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := it.SetID(newID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ffcfb3fc-efba-44b5-9b0e-695f22b1c5a1", err.Error())
	}

	return nil
}

// GetProductQty returns stock qty for a requested product-options pair
func (it *DefaultStock) GetProductQty(productID string, options map[string]interface{}) int {

	var qtySetFlag bool
	var minQty int

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		_ = env.ErrorDispatch(err)
		return minQty
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		_ = env.ErrorDispatch(err)
		return minQty
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		_ = env.ErrorDispatch(err)
		return minQty
	}

	// there could be couple matching request - we are looking for minimal value
	for _, dbRecord := range dbRecords {
		if !utils.StrKeysInMap(dbRecord, "qty", "options") {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c4d7d994-3f85-434e-9a72-8d3ab02eb063", "unexpected db result")
			break
		}

		recordOptions, ok := dbRecord["options"].(map[string]interface{})

		// skipping un-matching records
		if !ok || !utils.MatchMapAValuesToMapB(recordOptions, options) {
			continue
		}

		qty := utils.InterfaceToInt(dbRecord["qty"])

		if !qtySetFlag || qty < minQty {
			minQty = qty
			qtySetFlag = true
		}
	}

	return minQty
}

// GetProductOptions returns list of existing product options
func (it *DefaultStock) GetProductOptions(productID string) []map[string]interface{} {

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		env.LogError(err)
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		env.LogError(err)
	}

	productOptions, err := dbCollection.Load()
	if err != nil {
		env.LogError(err)
	}

	for _, productOption := range productOptions {
		if _, present := productOption["_id"]; present {
			delete(productOption, "_id")
		}
		if _, present := productOption["product_id"]; present {
			delete(productOption, "product_id")
		}

	}

	return productOptions
}

// RemoveProductQty removes database records matching given product-options pair
func (it *DefaultStock) RemoveProductQty(productID string, options map[string]interface{}) error {

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// collecting record ids to remove
	var stockIDsToRemove []string
	for _, dbRecord := range dbRecords {
		if recordOptions, present := dbRecord["options"]; present {
			recordOptions, ok := recordOptions.(map[string]interface{})

			// skipping un-matching records
			if !ok || !utils.MatchMapAValuesToMapB(options, recordOptions) {
				continue
			}

			if stockRecordID, present := dbRecord["_id"]; present {
				stockIDsToRemove = append(stockIDsToRemove, utils.InterfaceToString(stockRecordID))
			}
		}
	}

	// removing database records
	if len(stockIDsToRemove) > 0 {
		err := dbCollection.AddFilter("_id", "in", stockIDsToRemove)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		_, err = dbCollection.Delete()
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// SetProductQty updates stock qty for a product-options pair to be exact given value
//   - use UpdateProductQty in all cases you can as it is more safer option
func (it *DefaultStock) SetProductQty(productID string, options map[string]interface{}, qty int) error {

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("options", "=", options)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// looking for records matching request
	recordsProcessed := 0
	for _, dbRecord := range dbRecords {
		if !utils.StrKeysInMap(dbRecord, "_id", "qty", "options") {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b0d88054-0678-4894-831b-f332583a2ae7", "unexpected db result")
			break
		}

		recordOptions, ok := dbRecord["options"].(map[string]interface{})

		// skipping un-matching records
		if !ok || !utils.MatchMapAValuesToMapB(options, recordOptions) {
			continue
		}

		dbRecord["qty"] = qty
		_, err = dbCollection.Save(dbRecord)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		recordsProcessed++
	}

	// no records was - adding new
	if recordsProcessed == 0 {
		_, err := dbCollection.Save(map[string]interface{}{"product_id": productID, "options": options, "qty": qty})
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// UpdateProductQty updates stock qty for a product-options pair on delta value which can be positive or negative
func (it *DefaultStock) UpdateProductQty(productID string, options map[string]interface{}, deltaQty int) error {

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// looking for records matching request
	recordsProcessed := 0
	for _, dbRecord := range dbRecords {
		if !utils.StrKeysInMap(dbRecord, "_id", "qty", "options") {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f99c2f4b-ab56-4aa6-9950-00dcb09c15fd", "unexpected db result")
			break
		}

		recordOptions, ok := dbRecord["options"].(map[string]interface{})

		// skipping un-matching records
		if !ok || !utils.MatchMapAValuesToMapB(recordOptions, options) {
			continue
		}

		// TODO: should work through collection update function
		qty := utils.InterfaceToInt(dbRecord["qty"])
		dbRecord["qty"] = qty + deltaQty
		_, err := dbCollection.Save(dbRecord)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		recordsProcessed++
	}

	if recordsProcessed == 0 {
		return env.ErrorDispatch(env.ErrorNew(ConstErrorModule, 1, "62214642-d384-4474-b45c-e5fe8424fc3a", "Was given a set of options that didn't match any stock options in the db"))
	}

	return nil
}

// Load loads model from storage
func (it *DefaultStock) Load(id string) error {
	dbStockCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := dbStockCollection.LoadByID(id)
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
func (it *DefaultStock) Delete() error {
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// SetID sets database storage id for current object
func (it *DefaultStock) SetID(id string) error {
	it.id = id
	return nil
}

// GetID returns database storage id of current object
func (it *DefaultStock) GetID() string {
	return it.id
}

// SetProductID sets database storage id for current object
func (it *DefaultStock) SetProductID(product_id string) error {
	it.product_id = product_id
	return nil
}

// GetProductID returns database storage id of current object
func (it *DefaultStock) GetProductID() string {
	return it.product_id
}

// SetOptions sets database storage id for current object
func (it *DefaultStock) SetOptions(options string) error {
	it.options = options
	return nil
}

// GetOptions returns database storage id of current object
func (it *DefaultStock) GetOptions() string {
	return it.options
}

// SetOptions sets database storage id for current object
func (it *DefaultStock) SetQty(qty int) error {
	it.qty = qty
	return nil
}

// GetOptions returns database storage id of current object
func (it *DefaultStock) GetQty() int {
	return it.qty
}
