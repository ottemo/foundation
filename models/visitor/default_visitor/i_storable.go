package default_visitor

import (
	"github.com/ottemo/foundation/database"
)

func (it *DefaultVisitor) GetId() string {
	return it.id
}

func (it *DefaultVisitor) SetId(NewId string) error {
	it.id = NewId
	return nil
}

func (it *DefaultVisitor) Load(Id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {

			values, err := collection.LoadById(Id)
			if err != nil {
				return err
			}

			err = it.FromHashMap(values)
			if err != nil {
				return err
			}

		} else {
			return err
		}
	}
	return nil
}

func (it *DefaultVisitor) Delete(id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {
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

func (it *DefaultVisitor) Save() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {

			// prepearing initial hashmap
			storableValues := it.ToHashMap()

			delete(storableValues, "billing_address")
			delete(storableValues, "shipping_address")

			// shipping address save
			if it.ShippingAddress != nil {
				err := it.ShippingAddress.Save()
				if err != nil {
					return err
				}

				storableValues["shipping_address_id"] = it.ShippingAddress.GetId()
			}

			// billing address save
			if it.BillingAddress != nil {
				err := it.BillingAddress.Save()
				if err != nil {
					return err
				}

				storableValues["billing_address_id"] = it.BillingAddress.GetId()
			}

			// saving visitor
			if newId, err := collection.Save(storableValues); err == nil {
				it.Set("_id", newId)
			} else {
				return err
			}

		} else {
			return err
		}
	}
	return nil
}
