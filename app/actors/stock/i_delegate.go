package stock

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// --------------------
// Delegate declaration
// --------------------

// Stock delegate adds qty and inventory record to product model, providing possibility to updated them

// New instantiates delegate
func (it *StockDelegate) New(instance interface{}) (models.InterfaceAttributesDelegate, error) {
	if productModel, ok := instance.(product.InterfaceProduct); ok {
		return &StockDelegate{instance: productModel}, nil
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dafe7e34-ca3a-4e5b-b261-e25a6626914d", "unexpected instance for stock delegate")
}

// Get is a getter for external attributes
// This method should not read data from db. It is Load responsibility.
func (it *StockDelegate) Get(attribute string) interface{} {
	switch attribute {
	case "qty":
		return it.Qty
	case "inventory":
		return it.Inventory
	}
	return nil
}

// Set is a setter for external attributes, allow only to set value for current model
func (it *StockDelegate) Set(attribute string, value interface{}) error {
	switch attribute {
	case "qty":
		it.Qty = utils.InterfaceToInt(value)

	case "inventory":
		if haveInventoryOptionsDuplicates(value) {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c222525f-02f1-4bf0-bca5-204362c2d92d", "new inventory contains duplicates")
		}

		inventory := utils.InterfaceToArray(value)
		it.Inventory = []map[string]interface{}{}
		for _, options := range inventory {
			it.Inventory = append(it.Inventory, utils.InterfaceToMap(options))
		}
	}

	return nil
}

// GetAttributesInfo is a specification of external attributes
func (it *StockDelegate) GetAttributesInfo() []models.StructAttributeInfo {
	return []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameStock,
			Attribute:  "qty",
			Type:       utils.ConstDataTypeInteger,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Qty",
			Group:      "General",
			Editors:    "numeric",
			Options:    "",
			Default:    "0",
			Validators: "numeric positive",
		},
		models.StructAttributeInfo{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameStock,
			Attribute:  "inventory",
			Type:       utils.ConstDataTypeJSON,
			Label:      "Inventory",
			IsRequired: false,
			IsStatic:   false,
			Group:      "General",
			Editors:    "json",
			Options:    "",
			Default:    "",
			Validators: "",
		},
	}
}

// Load is a modelInstance.Load() method handler for external attributes, updates qty and inventory values
func (it *StockDelegate) Load(productID string) error {
	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		it.Qty = stockManager.GetProductQty(it.instance.GetID(), it.instance.GetAppliedOptions())
		it.Inventory = stockManager.GetProductOptions(it.instance.GetID())
	}

	return nil
}

// Save is a modelInstance.Save() method handler for external attributes, updates qty and inventory values
// methods toHashMap is called to Save instance so Get methods would be executed before Save
func (it *StockDelegate) Save() error {
	if stockManager := product.GetRegisteredStock(); stockManager != nil {

		// recalculate total qty if options present
		var itQty int
		var optionsPresent = false
		for _, productOptions := range it.Inventory {
			itQty += utils.InterfaceToInt(productOptions["qty"])
			optionsPresent = true
		}

		if optionsPresent {
			it.Qty = itQty
		}

		productID := it.instance.GetID()
		// remove current stock
		err := stockManager.RemoveProductQty(productID, make(map[string]interface{}))
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// set new stock
		err = stockManager.SetProductQty(productID, make(map[string]interface{}), it.Qty)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		for _, productOptions := range it.Inventory {
			options := utils.InterfaceToMap(productOptions["options"])
			qty := utils.InterfaceToInt(productOptions["qty"])

			err = stockManager.SetProductQty(productID, options, qty)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}

	return nil
}

// Delete is a modelInstance.Delete() method handler for external attributes
func (it *StockDelegate) Delete() error {
	// remove qty and inventory values from database
	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		if err := stockManager.RemoveProductQty(it.instance.GetID(), make(map[string]interface{})); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "45a37ef3-6fb7-4c57-a8cd-8a1306eb513e", err.Error())
		}
	}
	return nil
}
