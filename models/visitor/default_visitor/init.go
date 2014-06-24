package default_visitor

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/database"
	"github.com/ottemo/foundation/models"
)

func init() {
	instance := new(DefaultVisitor)

	models.RegisterModel("Visitor", instance)
	database.RegisterOnDatabaseStart(instance.SetupModel)

	api.RegisterOnEndPointStart(instance.SetupAPI)
}

func (it *DefaultVisitor) SetupModel() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {
			collection.AddColumn("email", "text", true)
			collection.AddColumn("first_name", "text", false)
			collection.AddColumn("last_name", "text", false)
			collection.AddColumn("billing_address_id", "int", false)
			collection.AddColumn("shipping_address_id", "int", false)
		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}

func (it *DefaultVisitor) SetupAPI() error {
	err := api.GetEndPoint().RegisterJsonAPI("visitor", "create", it.CreateVisitorAPI)
	if err != nil {
		return err
	}

	err = api.GetEndPoint().RegisterJsonAPI("visitor", "update", it.UpdateVisitorAPI)
	if err != nil {
		return err
	}

	err = api.GetEndPoint().RegisterJsonAPI("visitor", "load", it.LoadVisitorAPI)
	if err != nil {
		return err
	}

	return nil
}
