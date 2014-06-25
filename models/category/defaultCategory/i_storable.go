package defaultCategory

import (
	"github.com/ottemo/foundation/database"

	"errors"

	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/product"
)

func (it *DefaultCategory) GetId() string {
	return it.id
}

func (it *DefaultCategory) SetId(NewId string) error {
	it.id = NewId
	return nil
}

func (it *DefaultCategory) Load(Id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {

		// loading category
		categoryCollection, err := dbEngine.GetCollection(CATEGORY_COLLECTION_NAME)
		if err != nil {
			return err
		}

		values, err := categoryCollection.LoadById(Id)
		if err != nil {
			return err
		}

		err = it.FromHashMap(values)
		if err != nil {
			return err
		}

		// loading related products
		junctionCollection, err := dbEngine.GetCollection(CATEGORY_PRODUCT_JUNCTION_COLLECTION_NAME)
		if err != nil {
			return err
		}

		junctionCollection.AddFilter("category_id", "=", it.GetId())
		junctedProducts, err := junctionCollection.Load()
		if err != nil {
			return err
		}

		for _, junctedProduct := range junctedProducts {

			model, err := models.GetModel("Product")
			if err != nil {
				return err
			}

			if productModel, ok := model.(product.IProduct); ok {
				err := productModel.Load(junctedProduct["product_id"].(string))
				if err != nil {
					return err
				}
				it.Products = append(it.Products, productModel)
			} else {
				errors.New("model is not 'IProduct' capable")
			}
		}

	}
	return nil
}

func (it *DefaultCategory) Delete(Id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		//deleting category products join
		junctionCollection, err := dbEngine.GetCollection(CATEGORY_PRODUCT_JUNCTION_COLLECTION_NAME)
		if err != nil {
			return err
		}

		junctionCollection.AddFilter("category_id", "=", it.GetId())
		_, err = junctionCollection.Delete()
		if err != nil {
			return err
		}

		// deleting category
		categoryCollection, err := dbEngine.GetCollection(CATEGORY_COLLECTION_NAME)
		if err != nil {
			return err
		}

		err = categoryCollection.DeleteById(Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (it *DefaultCategory) Save() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {

		storableValues := map[string]interface{}{"name": it.Name}

		categoryCollection, err := dbEngine.GetCollection(CATEGORY_COLLECTION_NAME)
		if err != nil {
			return err
		}

		// saving category products
		junctionCollection, err := dbEngine.GetCollection(CATEGORY_PRODUCT_JUNCTION_COLLECTION_NAME)
		if err != nil {
			return err
		}

		junctionCollection.AddFilter("category_id", "=", it.GetId())
		_, err = junctionCollection.Delete()
		if err != nil {
			return err
		}

		for _, categoryProduct := range it.Products {
			categoryProduct.Save()

			junctionCollection.Save(map[string]interface{}{"category_id": it.GetId(), "product_id": categoryProduct.GetId()})
		}

		// saving collection
		if newId, err := categoryCollection.Save(storableValues); err == nil {
			it.Set("_id", newId)
		} else {
			return err
		}

	}
	return nil
}
