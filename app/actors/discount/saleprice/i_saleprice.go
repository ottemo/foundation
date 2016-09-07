package saleprice

import (
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// setupDB prepares system database for package usage
func (it *DefaultSalePrice) setupDB() error {
	dbSalePriceCollection, err := db.GetCollection(saleprice.ConstModelNameSalePriceCollection)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	modelAttributesInfo := salePriceModel.GetAttributesInfo()
	for _, attributeInfo := range modelAttributesInfo {
		if attributeInfo.Attribute != "_id" {
			dbSalePriceCollection.AddColumn(
				attributeInfo.Attribute,
				attributeInfo.Type,
				true)
		}
	}

	return nil
}
