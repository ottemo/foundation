package defaultproduct

import (
	"github.com/ottemo/foundation/database"
)

func (it *DefaultProductModel) GetId() string {
	return it.id
}

func (it *DefaultProductModel) SetId(NewId string) error {
	it.id = NewId
	return nil
}

func (it *DefaultProductModel) Load(loadId string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			if values, err := collection.LoadById(loadId); err == nil {
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

func (it *DefaultProductModel) Delete(id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			err := collection.DeleteById(Id)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (it *DefaultProductModel) Save() error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			if newId, err := collection.Save(it.ToHashMap()); err == nil {
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
