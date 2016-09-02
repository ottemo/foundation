package saleprice

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

func (it *DefaultSalePriceCollection) GetModelName() string {
	return saleprice.ConstModelNameSalePriceCollection
}

func (it *DefaultSalePriceCollection) GetImplementationName() string {
	return "Default" + saleprice.ConstModelNameSalePriceCollection
}

func (it *DefaultSalePriceCollection) New() (models.InterfaceModel, error) {
	return &DefaultSalePriceCollection{}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceCollection implementation (package "github.com/ottemo/foundation/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

func (it *DefaultSalePriceCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

func (it *DefaultSalePriceCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	// loading data from DB
	//---------------------
	dbItems, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	// converting db record to StructListItem
	//-----------------------------------
	for _, dbItemData := range dbItems {
		salePriceModel, err := saleprice.GetSalePriceModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		salePriceModel.FromHashMap(dbItemData)

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = salePriceModel.GetID()
		resultItem.Name = salePriceModel.GetProductID()+", "+
			utils.InterfaceToString(salePriceModel.GetAmount())+", "+
			utils.InterfaceToString(salePriceModel.GetStartDatetime())+", "+
			utils.InterfaceToString(salePriceModel.GetEndDatetime())
		resultItem.Image = ""
		resultItem.Desc = "For product ["+salePriceModel.GetProductID()+"], "+
			" set sale price ["+utils.InterfaceToString(salePriceModel.GetAmount())+"], "+
			" from ["+utils.InterfaceToString(salePriceModel.GetStartDatetime())+"], "+
			" to ["+utils.InterfaceToString(salePriceModel.GetEndDatetime())+"]"

		result = append(result, *resultItem)
	}

	return result, nil
}

func (it *DefaultSalePriceCollection) ListAddExtraAttribute(attribute string) error {
	// TODO: just for now no external attributes
	return nil
}

func (it *DefaultSalePriceCollection) ListFilterAdd(attribute string, operator string, value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

//func (it *DefaultSalePriceCollection) ListFilterReset() error {
//
//}

//func (it *DefaultSalePriceCollection) ListLimit(offset int, limit int) error {
//
//}


