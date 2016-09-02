package saleprice

import "github.com/ottemo/foundation/app/models/discount/saleprice"

// ListSalePrices returns list of sale price model items
func (it *DefaultSalePriceCollection) ListSalePrices() []saleprice.InterfaceSalePrice {
	var result []saleprice.InterfaceSalePrice

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		salePriceModel, err := saleprice.GetSalePriceModel()
		if err != nil {
			return result
		}
		salePriceModel.FromHashMap(recordData)

		result = append(result, salePriceModel)
	}

	return result
}

