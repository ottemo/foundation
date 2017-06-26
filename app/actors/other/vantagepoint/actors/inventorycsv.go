package actors

import (
	"encoding/csv"
	"fmt"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"io"
)

type inventorySCV struct {
	header []string
}

func NewInventoryProcessor() (*inventorySCV, error) {
	processor := &inventorySCV{}

	return processor, nil
}

func (it *inventorySCV) Process(reader io.Reader) error {
	csvReader := csv.NewReader(reader)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return env.ErrorDispatch(err)
		}

		if it.header == nil {
			it.prepareHeader(record)
		} else {
			it.processRecord(record)
		}
	}

	return nil
}

func (it *inventorySCV) prepareHeader(record []string) {
	for _, originalKey := range record {
		key := originalKey
		switch originalKey {
		case "UPC Number":
			key = "sku"
		case "Stock":
			key = "qty"
		}
		it.header = append(it.header, key)
	}
}

func (it *inventorySCV) processRecord(record []string) error {
	item := map[string]string{}
	for idx, key := range it.header {
		if idx < len(record) {
			item[key] = record[idx]
		}
	}

	qty := utils.InterfaceToInt(item["qty"])

	if err := it.updateInventoryBySku(item["sku"], qty); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func (it *inventorySCV) updateInventoryBySku(sku string, qty int) error {
	fmt.Println("inventoryProcessor) updateProductInventory === START ===")
	collection, err := product.GetProductCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.ListFilterAdd("sku", "=", sku); err != nil {
		return env.ErrorDispatch(err)
	}

	products := collection.ListProducts()
	fmt.Println("products", utils.InterfaceToString(products))

	if len(products) > 1 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a8bf1294-539f-4cad-adb0-362b878e30eb", "morethen one product with sku "+sku)
	} else if len(products) == 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d491f656-d477-4e7b-9912-2682b12ac34b", "no products with sku "+sku)
	} else {
		productID := products[0].GetID()

		if err := it.updateProductInventory(productID, qty); err != nil {
			return env.ErrorDispatch(err)
		}

		if err := it.updateInventoryForOptions(productID, qty); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	fmt.Println("inventoryProcessor) updateProductInventory === DONE ===")
	return nil
}

func (it *inventorySCV) updateProductInventory(productID string, qty int) error {
	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "824e8aae-0021-40d7-b974-96a3fdbf8486", "stock is undefined")
	} else {
		options := stockManager.GetProductOptions(productID)
		fmt.Println("OPTIONS ===", productID, options)

		// unable to update product with options, because at the time of writing there were
		// no clear mapping of imported options to existing ones
		if len(options) > 1 {
			msg := fmt.Sprintf("product [%s] have more than one options set", productID)
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2d65a05e-661c-439d-abaa-3f4f90f9f2a4", msg)
		}

		stockManager.SetProductQty(productID, map[string]interface{}{}, qty)
	}

	return nil
}

func (it *inventorySCV) updateInventoryForOptions(optionProductID string, qty int) error {
	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "824e8aae-0021-40d7-b974-96a3fdbf8486", "stock is undefined")
	}

	collection, err := product.GetProductCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.ListFilterAdd("options", "LIKE", optionProductID); err != nil {
		return env.ErrorDispatch(err)
	}

	products := collection.ListProducts()
	fmt.Println("OPTION PRODUCTS", utils.InterfaceToString(products))

	// need nested loops because of "options" nature
	for _, foundProduct := range products {
		options := foundProduct.GetOptions()

		selectedOptions := map[string]interface{}{}

		for _, optionInterface := range options {
			option := utils.InterfaceToMap(optionInterface)
			optionKey := utils.InterfaceToString(option["key"])
			optionOptionsInterface := option["options"]
			optionOptions := utils.InterfaceToMap(optionOptionsInterface)

			for _, optionsOptionsInterface := range optionOptions {
				optionsOptions := utils.InterfaceToMap(optionsOptionsInterface)
				optionsOptionKey := optionsOptions["key"]
				_ids := optionsOptions["_ids"]

				if utils.IsInArray(optionProductID, _ids) {
					selectedOptions[optionKey] = optionsOptionKey
				}
			}
		}

		oldQty := stockManager.GetProductQty(foundProduct.GetID(), selectedOptions)
		fmt.Println("OLD QTY === ", foundProduct.GetID(), selectedOptions)

		if err := stockManager.UpdateProductQty(foundProduct.GetID(), selectedOptions, qty-oldQty); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}
