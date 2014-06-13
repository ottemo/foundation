package default_address

import (
	"errors"
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/database"
)

func init(){
	instance := new(DefaultVisitorAddress)

	models.RegisterModel("VisitorAddress", instance )
	database.RegisterOnDatabaseStart( instance.SetupModel )
}


func (it *DefaultVisitorAddress) SetupModel() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection( VISITOR_ADDRESS_COLLECTION_NAME ); err == nil {
			collection.AddColumn("street", "text", false)
			collection.AddColumn("city", "text", false)
			collection.AddColumn("state", "text", false)
			collection.AddColumn("phone", "text", false)
			collection.AddColumn("zip_code", "text", false)
		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
