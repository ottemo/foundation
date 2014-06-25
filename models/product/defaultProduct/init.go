package defaultProduct

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/database"
	"github.com/ottemo/foundation/models"
)

func init() {
	models.RegisterModel("Product", new(ProductModel))
	database.RegisterOnDatabaseStart(SetupModel)

	api.RegisterOnServiceStart(SetupAPI)
}

func SetupModel() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			collection.AddColumn("sku", "text", true)
			collection.AddColumn("name", "text", true)
		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}

func SetupAPI() error {
	err := api.GetService().RegisterJsonAPI("product", "POST", "addAttribute", AddProductAttributeRestAPI)
	if err != nil {
		return err
	}

	err = api.GetService().RegisterJsonAPI("product", "POST", "createProduct", CreateProductRestAPI)
	if err != nil {
		return err
	}

	err = api.GetService().RegisterJsonAPI("product", "GET", "loadProduct", LoadProductRestAPI)
	if err != nil {
		return err
	}

	return nil
}
