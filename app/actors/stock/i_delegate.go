package stock

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// --------------------
// Delegate declaration
// --------------------

// Stock delegate adds qty and inventory record to product model, providing possibility to updated them

// New instantiates delegate
func (it *StockDelegate) New(instance interface{}) (models.InterfaceAttributesDelegate, error) {
	env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) New")
	if productModel, ok := instance.(product.InterfaceProduct); ok {
		return &StockDelegate{instance: productModel}, nil
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dafe7e34-ca3a-4e5b-b261-e25a6626914d", "unexpected instance for stock delegate")
}

// Get is a getter for external attributes
func (it *StockDelegate) Get(attribute string) interface{} {
	env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Get "+attribute)
	switch attribute {
	case "qty":
		//if stockManager := product.GetRegisteredStock(); stockManager != nil {
		//	it.Qty = stockManager.GetProductQty(it.instance.GetID(), it.instance.GetAppliedOptions())
		//}
		env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Get "+attribute+"="+utils.InterfaceToString(it.Qty))
		return it.Qty
	case "inventory":
		//if it.Inventory == nil {
		//	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		//		it.Inventory = stockManager.GetProductOptions(it.instance.GetID())
		//	}
		//}

		env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Get "+attribute+"="+utils.InterfaceToString(it.Inventory))
		return it.Inventory
	}
	return nil
}

// Set is a setter for external attributes, allow only to set value for current model
func (it *StockDelegate) Set(attribute string, value interface{}) error {
	env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Set "+attribute+"="+utils.InterfaceToString(value))
	switch attribute {
	case "qty":
		env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Set 1 "+utils.InterfaceToString(it.Qty))
		env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Set 2 "+utils.InterfaceToString(utils.InterfaceToInt(value)))
		it.Qty = utils.InterfaceToInt(value)
		env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Set 3 "+utils.InterfaceToString(it.Qty))

	case "inventory":
		inventory := utils.InterfaceToArray(value)
		for _, options := range inventory {
			it.Inventory = append(it.Inventory, utils.InterfaceToMap(options))
		}
	}

	return nil
}

// GetAttributesInfo is a specification of external attributes
func (it *StockDelegate) GetAttributesInfo() []models.StructAttributeInfo {
	env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) GetAttributesInfo")
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
	env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Save: "+utils.InterfaceToString(it));
	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		productID := it.instance.GetID()
		env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Save: productID "+utils.InterfaceToString(productID));
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

		env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Save Inventory start");
		for _, productOptions := range it.Inventory {
			env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Save Inventory options "+utils.InterfaceToString(productOptions));
			options := utils.InterfaceToMap(productOptions["options"])
			qty := utils.InterfaceToInt(productOptions["qty"])

			err = stockManager.SetProductQty(productID, options, qty)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}
		env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Save Inventory done");
	}

	return nil
}

// Delete is a modelInstance.Delete() method handler for external attributes
func (it *StockDelegate) Delete() error {
	env.Log("errors.log", env.ConstLogPrefixDebug, "StockDelegate) Delete")
	// remove qty and inventory values from database
	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		stockManager.RemoveProductQty(it.instance.GetID(), make(map[string]interface{}))
	}
	return nil
}
