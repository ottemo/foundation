package saleprice

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
)

// setupDB prepares system database for package usage
func (it *DefaultSalePrice) setupDB() error {
	dbSalePriceCollection, err := db.GetCollection(ConstCollectionNameSalePrices)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	modelAttributesInfo := salePriceModel.GetAttributesInfo()
	for _, attributeInfo := range(modelAttributesInfo) {
		//TODO: Strut. ID field should be filtered somewhat else.
		if attributeInfo.Attribute != "_id" {
			// TODO: AddColumn has parameter indexed, which is not
			// predicted in StructAttributeInfo
			dbSalePriceCollection.AddColumn(
				attributeInfo.Attribute,
				attributeInfo.Type,
				true)
		}
	}

	return nil
}
