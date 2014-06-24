package default_address

import (
	"github.com/ottemo/foundation/database"
)

func (dva *DefaultVisitorAddress) GetId() string {
	return dva.id
}

func (dva *DefaultVisitorAddress) SetId(NewId string) error {
	dva.id = NewId
	return nil
}

func (dva *DefaultVisitorAddress) Load(Id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_ADDRESS_COLLECTION_NAME); err == nil {

			if values, err := collection.LoadByID(Id); err == nil {
				if err := dva.FromHashMap(values); err != nil {
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

func (dva *DefaultVisitorAddress) Delete(Id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_ADDRESS_COLLECTION_NAME); err == nil {
			err := collection.DeleteByID(Id)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (dva *DefaultVisitorAddress) Save() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_ADDRESS_COLLECTION_NAME); err == nil {

			//if it.ZipCode== "" {
			//	return errors.New("Zip code for address - required")
			//}

			if newId, err := collection.Save(dva.ToHashMap()); err == nil {
				dva.Set("_id", newId)
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
