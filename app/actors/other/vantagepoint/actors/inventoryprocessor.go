package actors

import (
	"io"
	"fmt"
	"encoding/csv"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
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
	fmt.Println(utils.InterfaceToString(item))

	return nil
}
