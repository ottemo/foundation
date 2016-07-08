package cart

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/env"

	"time"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
)

// GetID returns id of current cart
func (it *DefaultCart) GetID() string {
	return it.id
}

// SetID sets id for cart
func (it *DefaultCart) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load loads cart information from DB
func (it *DefaultCart) Load(ID string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {

		// loading category
		cartCollection, err := dbEngine.GetCollection(ConstCartCollectionName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		values, err := cartCollection.LoadByID(ID)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// initializing DefaultCart structure
		it.id = utils.InterfaceToString(values["_id"])
		it.Active = utils.InterfaceToBool(values["active"])
		it.VisitorID = utils.InterfaceToString(values["visitor_id"])
		it.SessionID = utils.InterfaceToString(values["session_id"])
		it.UpdatedAt = utils.InterfaceToTime(values["updated_at"])
		it.Info, _ = utils.DecodeJSONToStringKeyMap(values["info"])
		it.CustomInfo = utils.InterfaceToMap(values["custom_info"])
		it.Items = make(map[int]cart.InterfaceCartItem)
		it.maxIdx = 0

		// loading cart items
		cartItemsCollection, err := dbEngine.GetCollection(ConstCartItemsCollectionName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		cartItemsCollection.AddFilter("cart_id", "=", it.GetID())
		cartItems, err := cartItemsCollection.Load()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		for _, cartItemValues := range cartItems {

			cartItem := new(DefaultCartItem)

			cartItem.id = utils.InterfaceToString(cartItemValues["_id"])
			cartItem.idx = utils.InterfaceToInt(cartItemValues["idx"])

			if cartItem.idx > it.maxIdx {
				it.maxIdx = cartItem.idx
			}

			cartItem.Cart = it
			cartItem.ProductID = utils.InterfaceToString(cartItemValues["product_id"])
			cartItem.Qty = utils.InterfaceToInt(cartItemValues["qty"])
			cartItem.Options = utils.InterfaceToMap(cartItemValues["options"])

			// check for product existence and options validation
			if itemProduct := cartItem.GetProduct(); itemProduct != nil &&
				it.checkOptions(itemProduct.GetOptions(), cartItem.GetOptions()) == nil {

				it.Items[cartItem.idx] = cartItem
			} else {
				cartItemsCollection.DeleteByID(utils.InterfaceToString(cartItemValues["_id"]))
			}

		}
	}

	return nil
}

// Delete removes current cart from DB
func (it *DefaultCart) Delete() error {
	if it.GetID() == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "75d83a42-db90-4917-969d-39dd7fb50943", "cart id is not set")
	}

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b168a568-6f15-4612-90e1-f924501c1f0a", "can't get DbEngine")
	}

	// deleting cart items
	cartItemsCollection, err := dbEngine.GetCollection(ConstCartItemsCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = cartItemsCollection.AddFilter("cart_id", "=", it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	_, err = cartItemsCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// deleting cart
	cartCollection, err := dbEngine.GetCollection(ConstCartCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = cartCollection.DeleteByID(it.GetID())

	return env.ErrorDispatch(err)
}

// Save stores current cart in DB
func (it *DefaultCart) Save() error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3ca9e702-7a81-46a9-a950-02c9a246f581", "can't get DbEngine")
	}

	cartCollection, err := dbEngine.GetCollection(ConstCartCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	cartItemsCollection, err := dbEngine.GetCollection(ConstCartItemsCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// packing data before save
	cartStoringValues := make(map[string]interface{})

	cartStoringValues["_id"] = it.GetID()
	cartStoringValues["visitor_id"] = it.VisitorID
	cartStoringValues["session_id"] = it.SessionID
	cartStoringValues["active"] = it.Active
	cartStoringValues["info"] = utils.EncodeToJSONString(it.Info)
	cartStoringValues["custom_info"] = it.CustomInfo
	cartStoringValues["updated_at"] = time.Now()

	newID, err := cartCollection.Save(cartStoringValues)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	it.SetID(newID)

	// storing cart items
	for _, cartItem := range it.GetItems() {
		cartItemStoringValues := make(map[string]interface{})

		cartItemStoringValues["_id"] = cartItem.GetID()
		cartItemStoringValues["idx"] = cartItem.GetIdx()
		cartItemStoringValues["cart_id"] = it.GetID()
		cartItemStoringValues["product_id"] = cartItem.GetProductID()
		cartItemStoringValues["qty"] = cartItem.GetQty()
		cartItemStoringValues["options"] = cartItem.GetOptions()

		newID, err := cartItemsCollection.Save(cartItemStoringValues)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		cartItem.SetID(newID)
	}

	return nil
}
