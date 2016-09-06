package saleprice

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/product"
)

func setupConfig() error {
	// TODO suggest sale price module registered and enabled

	productModel, err := product.GetProductModel()
	if err != nil {
		env.LogError(err)
	}

	if err = productModel.AddExternalAttributes(salePriceDelegate); err != nil {
		env.LogError(err)
	}

	return nil
}
