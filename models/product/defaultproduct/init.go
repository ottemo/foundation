package defaultproduct

import (
	"errors"

	"github.com/ottemo/foundation/database"
	"github.com/ottemo/foundation/models"

	"github.com/ottemo/foundation/rest_service"
)

func init() {
	models.RegisterModel("Product", new(DefaultProductModel))
	database.RegisterOnDatabaseStart(SetupModel)

	rest_service.RegisterOnRestServiceStart(SetupAPI)
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
	err := rest_service.GetRestService().RegisterJsonAPI("product", "addAttribute", AddProductAttributeRestAPI)
	if err != nil {
		return err
	}

	err = rest_service.GetRestService().RegisterJsonAPI("product", "createProduct", CreateProductRestAPI)
	if err != nil {
		return err
	}

	err = rest_service.GetRestService().RegisterJsonAPI("product", "loadProduct", LoadProductRestAPI)
	if err != nil {
		return err
	}

	return nil
}
