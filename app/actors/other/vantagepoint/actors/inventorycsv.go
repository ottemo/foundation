package actors

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/product"
)

// --------------------------------------------------------------------------------------------------------------------

type inventoryCSV struct {
	env   EnvInterface

	header []string
}

func NewInventoryCSV(env EnvInterface) (*inventoryCSV, error) {
	processor := &inventoryCSV{
		env:   env,
	}

	return processor, nil
}

func (it *inventoryCSV) Process(reader io.Reader) error {
	csvReader := csv.NewReader(reader)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return it.env.ErrorDispatch(err)
		}

		if it.header == nil {
			it.prepareHeader(record)
		} else {
			if err := it.processRecord(record); err != nil {
				return it.env.ErrorDispatch(err)
			}
		}
	}

	return nil
}

func (it *inventoryCSV) prepareHeader(record []string) {
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

func (it *inventoryCSV) processRecord(record []string) error {
	item := map[string]string{}
	for idx, key := range it.header {
		if idx < len(record) {
			item[key] = record[idx]
		}
	}

	qty := utils.InterfaceToInt(item["qty"])

	if err := it.updateInventoryBySku(item["sku"], qty); err != nil {
		return it.env.ErrorDispatch(err)
	}

	return nil
}

func (it *inventoryCSV) updateInventoryBySku(sku string, qty int) error {
	collection, err := product.GetProductCollectionModel()
	if err != nil {
		return it.env.ErrorDispatch(err)
	}

	if err := collection.ListFilterAdd("sku", "=", sku); err != nil {
		return it.env.ErrorDispatch(err)
	}

	products := collection.ListProducts()

	if len(products) > 1 {
		return it.env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a8bf1294-539f-4cad-adb0-362b878e30eb", "morethen one product with sku "+sku)
	} else if len(products) == 0 {
		return it.env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d491f656-d477-4e7b-9912-2682b12ac34b", "no products with sku "+sku)
	} else {
		productID := products[0].GetID()

		if err := it.updateProductInventory(productID, qty); err != nil {
			return it.env.ErrorDispatch(err)
		}

		if err := it.updateInventoryForOptions(productID, qty); err != nil {
			return it.env.ErrorDispatch(err)
		}
	}

	return nil
}

func (it *inventoryCSV) updateProductInventory(productID string, qty int) error {
	stock := product.GetRegisteredStock()
	if stock == nil {
		return it.env.ErrorNew(ConstErrorModule, ConstErrorLevel, "824e8aae-0021-40d7-b974-96a3fdbf8486", "stock is undefined")
	}

	options := stock.GetProductOptions(productID)

	// unable to update product with options, because at the time of writing there were
	// no clear mapping of imported options to existing ones
	if len(options) > 1 {
		msg := fmt.Sprintf("product [%s] have more than one options set", productID)
		return it.env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2d65a05e-661c-439d-abaa-3f4f90f9f2a4", msg)
	}

	return stock.SetProductQty(productID, map[string]interface{}{}, qty)
}

func (it *inventoryCSV) updateInventoryForOptions(optionProductID string, qty int) error {
	stock := product.GetRegisteredStock()
	if stock == nil {
		return it.env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f3656381-c997-47bc-a3e6-10485a2b6d4d", "stock is undefined")
	}

	collection, err := product.GetProductCollectionModel()
	if err != nil {
		return it.env.ErrorDispatch(err)
	}

	if err := collection.ListFilterAdd("options", "LIKE", optionProductID); err != nil {
		return it.env.ErrorDispatch(err)
	}

	products := collection.ListProducts()

	// need nested loops because of "options" nature
	for _, foundProduct := range products {
		options := foundProduct.GetOptions()

		selectedOptions := map[string]interface{}{}

		for _, optionInterface := range options {
			option := utils.InterfaceToMap(optionInterface)

			if !utils.StrKeysInMap(option, "key", "options") {
				continue
			}

			optionKey := utils.InterfaceToString(option["key"])
			optionOptionsInterface := option["options"]
			optionOptions := utils.InterfaceToMap(optionOptionsInterface)

			for _, optionsOptionsInterface := range optionOptions {
				optionsOptions := utils.InterfaceToMap(optionsOptionsInterface)

				if !utils.StrKeysInMap(option, "key", "_ids") {
					continue
				}

				optionsOptionKey := optionsOptions["key"]
				_ids := optionsOptions["_ids"]

				if utils.IsInArray(optionProductID, _ids) {
					selectedOptions[optionKey] = optionsOptionKey
				}
			}
		}

		oldQty := stock.GetProductQty(foundProduct.GetID(), selectedOptions)

		if err := stock.UpdateProductQty(foundProduct.GetID(), selectedOptions, qty-oldQty); err != nil {
			return it.env.ErrorDispatch(err)
		}
	}

	return nil
}
