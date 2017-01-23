package stock

import (
	"github.com/ottemo/foundation/utils"
)

// haveInventoryOptionsDuplicates checks inventory for duplicates
func haveInventoryOptionsDuplicates(inventory interface{}) bool {
	var inventoryArray = utils.InterfaceToArray(inventory)

	for idxA, itemA := range inventoryArray {
		var itemMapA = utils.InterfaceToMap(itemA)
		if _, present := itemMapA["options"]; present {
			var optionsA = utils.InterfaceToMap(itemMapA["options"])
			for idxB, itemB := range inventoryArray {
				var itemMapB = utils.InterfaceToMap(itemB)
				if idxB > idxA {
					if _, present := itemMapB["options"]; present {
						var optionsB = utils.InterfaceToMap(itemMapB["options"])
						var isEqual = utils.MatchMapAValuesToMapB(optionsA, optionsB) && utils.MatchMapAValuesToMapB(optionsB, optionsA)

						if isEqual {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
