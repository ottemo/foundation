package saleprice

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GetSalePriceModel retrieves current InterfaceSalePrice model implementation
func GetSalePriceModel() (InterfaceSalePrice, error) {
	model, err := models.GetModel(ConstModelNameSalePrice)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	salePriceModel, ok := model.(InterfaceSalePrice)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1ef4f461-e355-4989-9d79-75b060d8153c", "model "+model.GetImplementationName()+" is not 'InterfaceSalePrice' capable")
	}

	return salePriceModel, nil
}
