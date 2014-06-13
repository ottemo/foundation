package default_address

import (
	"github.com/ottemo/foundation/database"
)

func (it *DefaultVisitorAddress) GetId() string {
	return it.id
}

func (it *DefaultVisitorAddress) SetId(NewId string) error {
	it.id = NewId
	return nil
}

func (it *DefaultVisitorAddress) Load(Id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection( VISITOR_ADDRESS_COLLECTION_NAME ); err == nil {

			if values, err := collection.LoadById( Id ); err == nil {
				if err := it.FromHashMap(values); err != nil {
					return err
				}
			} else {
				return err
			}

		} else {
			return err
		}
	}
	return nil
}

func (it *DefaultVisitorAddress) Delete(Id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection( VISITOR_ADDRESS_COLLECTION_NAME ); err == nil {
			err := collection.DeleteById(Id)
			if err != nil { return err }
		} else {
			return err
		}
	}
	return nil
}

func (it *DefaultVisitorAddress) Save() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection( VISITOR_ADDRESS_COLLECTION_NAME ); err == nil {

			//if it.ZipCode== "" {
			//	return errors.New("Zip code for address - required")
			//}

			if newId, err := collection.Save( it.ToHashMap() ); err == nil {
				it.Set("_id", newId)
				return err
			} else {
				return err
			}

		} else {
			return err
		}
	}
	return nil
}
