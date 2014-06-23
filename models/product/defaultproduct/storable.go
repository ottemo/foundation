package defaultproduct

import (
	"github.com/ottemo/foundation/database"
)

func (dpm *DefaultProductModel) GetId() string {
	return dpm.id
}

func (dpm *DefaultProductModel) SetId(NewId string) error {
	dpm.id = NewId
	return nil
}

func (dpm *DefaultProductModel) Load(loadId string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			if values, err := collection.LoadById(loadId); err == nil {
				if err := dpm.FromHashMap(values); err != nil {
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

func (dpm *DefaultProductModel) Delete(id string) error {
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

func (dpm *DefaultProductModel) Save() error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			if newId, err := collection.Save(it.ToHashMap()); err == nil {
				dpm.Set("_id", newId)
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
