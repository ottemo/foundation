package stock

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetProductQty returns stock qty for a requested product-options pair
func (it *DefaultStock) GetProductQty(productID string, options map[string]interface{}) int {

	var qtySetFlag bool
	var minQty int

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		env.ErrorDispatch(err)
		return minQty
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		env.ErrorDispatch(err)
		return minQty
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		env.ErrorDispatch(err)
		return minQty
	}

	// there could be couple matching request - we are looking for minimal value
	for _, dbRecord := range dbRecords {
		if !utils.StrKeysInMap(dbRecord, "qty", "options") {
			env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c4d7d994-3f85-434e-9a72-8d3ab02eb063", "unexpected db result")
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
			env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b0d88054-0678-4894-831b-f332583a2ae7", "unexpected db result")
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
			env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f99c2f4b-ab56-4aa6-9950-00dcb09c15fd", "unexpected db result")
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
