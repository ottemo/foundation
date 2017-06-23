package actors

import (
	"io"
	"fmt"
	"encoding/csv"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/app/models/product"
)

type inventoryProcessor struct {
	header []string
}

func NewInventoryProcessor() (*inventoryProcessor, error) {
	processor := &inventoryProcessor{}

	return processor, nil
}

func (it *inventoryProcessor) Process(reader io.Reader) error {
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

func (it *inventoryProcessor) prepareHeader(record []string) {
	for _, originalKey := range(record) {
		key := originalKey
		switch originalKey {
		case "UPC Number" :
			key = "sku"
		case "Stock" :
			key = "qty"
		}
		it.header = append(it.header, key)
	}
}

func (it *inventoryProcessor) processRecord(record []string) error {
	item := map[string]string{}
	for idx, key := range(it.header) {
		if idx < len(record) {
			item[key] = record[idx]
		}
	}

	if err := it.updateInventoryBySku(item["sku"], item["qty"]); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func (it *inventoryProcessor) updateInventoryBySku(sku, qty string) error {
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

		return it.updateProductInventory(productID, qty)
	}

	fmt.Println("inventoryProcessor) updateProductInventory === DONE ===")
	return nil
}

func (it *inventoryProcessor) updateProductInventory(productID, qty string) error {
	stockManager := product.GetRegisteredStock()
	if stockManager != nil {
		options := stockManager.GetProductOptions(productID)

		_ = options

		// TODO what if have options, but not configurable?
	}

	return nil
}
