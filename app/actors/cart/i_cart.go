package cart

import (
	"sort"
	"strconv"
	"time"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// function calls when cart was changes and subtotal recalculation needed
func (it *DefaultCart) cartChanged() {
	it.Subtotal = 0
}

// function checks that options were set correctly (required elements, option values)
func (it *DefaultCart) checkOptions(productOptions map[string]interface{}, cartItemOptions map[string]interface{}) error {

	// loop 1: checking that all products customer set are available for product
	for optionName, optionValue := range cartItemOptions {

		// checking if product have attribute customer set
		if productOption, present := productOptions[optionName]; present {

			// checking that product option values are strictly predefined
			if productOption, ok := productOption.(map[string]interface{}); ok {
				if _, present := productOption["options"]; present {

					// checking for valid value was set by customer
					// cart option value can be one or multiple values, but should be string there
					var optionValuesToCheck []string
					switch typedOptionValue := optionValue.(type) {
					case string:
						optionValuesToCheck = append(optionValuesToCheck, typedOptionValue)
					case []string:
						optionValuesToCheck = typedOptionValue
					case []interface{}:
						for _, value := range typedOptionValue {
							if value, ok := value.(string); ok {
								optionValuesToCheck = append(optionValuesToCheck, value)
							}
						}
					default:
						return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8e735b08-b2b3-4bb8-aa74-bb8ebfe86369", "unexpected option value for '"+optionName+"' option")
					}

					// checking for option customer set with available for product
					for _, optionValue := range optionValuesToCheck {

						// there could be couple forms of product available options specification
						switch productOptionValues := productOption["options"].(type) {
						case map[string]interface{}:
							if _, present := productOptionValues[optionValue]; !present {
								return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "482efc22-b29e-44e2-9ed6-b78505a56d4d", "invalid value for option '"+optionName+"'")
							}

						case []interface{}:
							found := false
							for _, productOptionValue := range productOptionValues {
								if productOptionValue == optionValue {
									found = true
									break
								}
							}
							if !found {
								return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aad634c1-fe7b-42fb-adb7-58dfc0e3591b", "invalid value for option '"+optionName+"'")
							}

						default:
							if productOptionValues != optionValue {
								return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "eec3b62d-3aad-429c-821c-1b243fd80fd9", "invalid value for option '"+optionName+"'")
							}
						}
					}
				}
			}
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "504c4098-c7fe-418d-9857-adf2a08a8c3f", "unknown option '"+optionName+"'")
		}
	}

	// loop 2: checking that all product required options were set
	for productOption, productOptionValue := range productOptions {
		// required product option means that "productOption["required"]=true" should be set
		if productOptionValue, ok := productOptionValue.(map[string]interface{}); ok {
			if _, present := productOptionValue["required"]; present {
				if value, ok := productOptionValue["required"].(bool); ok && value {

					//checking cart item option for required option existence
					itemOptionValue, present := cartItemOptions[productOption]
					if !present {
						return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7d23e4a5-4f3c-4986-8860-c49008510175", productOption+" was not specified")
					}

					// for multi value options additional check
					switch typedValue := itemOptionValue.(type) {
					case []interface{}:
						if len(typedValue) == 0 {
							return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c03fd072-5aa1-48d9-bb60-f03d23ed6ba4", productOption+" was not specified")
						}
					}

				}
			}
		}
	}

	return nil
}

// AddItem adds item to the current cart
//   - returns added item or nil if error happened
func (it *DefaultCart) AddItem(productID string, qty int, options map[string]interface{}) (cart.InterfaceCartItem, error) {

	//checking qty
	if qty <= 0 {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "653e6163-0775-41c6-9ff5-4d8ac93bed5e", "qty can't be zero or less")
	}

	// options default value if them are not set
	if options == nil {
		options = make(map[string]interface{})
	}

	// looking for already present item in cart
	newItemOptions := utils.EncodeToJSONString(options)
	cartItems := it.GetItems()
	for _, item := range cartItems {
		cartItemOptions := utils.EncodeToJSONString(item.GetOptions())
		if item.GetProductID() == productID && cartItemOptions == newItemOptions {
			err := item.SetQty(item.GetQty() + qty)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
			return item, nil
		}
	}

	// if we are there then new item was added to cart
	// preparing new cart item
	cartItem := &DefaultCartItem{
		idx:       it.maxIdx,
		ProductID: productID,
		Qty:       qty,
		Options:   options,
		Cart:      it,
	}

	// validate cart item before add to cart
	if err := cartItem.ValidateProduct(); err != nil {
		return nil, err
	}

	// adding new item to others
	it.maxIdx++
	cartItem.idx = it.maxIdx
	it.Items[it.maxIdx] = cartItem

	// cart change event broadcast
	it.cartChanged()

	return cartItem, nil
}

// RemoveItem removes item from cart
//   - you need to know index you can get from ListItems()
func (it *DefaultCart) RemoveItem(itemIdx int) error {
	if cartItem, present := it.Items[itemIdx]; present {

		dbEngine := db.GetDBEngine()
		if dbEngine == nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d68c0e08-8bc6-4884-89b0-f8f949266ec4", "can't get DB engine")
		}

		cartItemsCollection, err := dbEngine.GetCollection(ConstCartItemsCollectionName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = cartItemsCollection.DeleteByID(cartItem.GetID())
		if err != nil {
			return env.ErrorDispatch(err)
		}

		delete(it.Items, itemIdx)

		it.cartChanged()

		return nil
	}
	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f72d898f-477d-42b4-9a69-8ccfffeb298d", "can't find index "+strconv.Itoa(itemIdx))
}

// SetQty sets new qty for particular item in cart
//   - you need to it's index, use ListItems() for that
func (it *DefaultCart) SetQty(itemIdx int, qty int) error {
	cartItem, present := it.Items[itemIdx]
	if present {
		err := cartItem.SetQty(qty)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		it.cartChanged()

		return nil
	}
	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "77005096-3b3f-46bc-8c67-4b7df255fdaf", "there is no item with idx="+strconv.Itoa(itemIdx))
}

// GetSubtotal returns subtotal for cart items
func (it *DefaultCart) GetSubtotal() float64 {

	if it.Subtotal == 0 {
		for _, cartItem := range it.Items {
			if cartProduct := cartItem.GetProduct(); cartProduct != nil {
				if err := cartProduct.ApplyOptions(cartItem.GetOptions()); err != nil {
					_ = env.ErrorDispatch(err)
				}
				it.Subtotal += cartProduct.GetPrice() * float64(cartItem.GetQty())
			}
		}
	}

	return it.Subtotal
}

// GetItems enumerates current cart items sorted by item idx
func (it *DefaultCart) GetItems() []cart.InterfaceCartItem {

	var result []cart.InterfaceCartItem

	var keys []int
	for key := range it.Items {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {
		result = append(result, it.Items[key])
	}

	return result
}

// GetVisitorID returns visitor id this cart belongs to
func (it *DefaultCart) GetVisitorID() string {
	return it.VisitorID
}

// SetVisitorID sets new owner of cart
func (it *DefaultCart) SetVisitorID(visitorID string) error {
	it.VisitorID = visitorID
	return nil
}

// GetVisitor returns visitor model represents owner or current cart or nil if visitor was not set to cart
func (it *DefaultCart) GetVisitor() visitor.InterfaceVisitor {
	visitorInstance, err := visitor.LoadVisitorByID(it.VisitorID)
	if err != nil {
		_ = env.ErrorDispatch(err)
	}
	return visitorInstance
}

// SetCartInfo assigns some information to current cart
func (it *DefaultCart) SetCartInfo(infoAttribute string, infoValue interface{}) error {
	if it.Info == nil {
		it.Info = make(map[string]interface{})
	}

	it.Info[infoAttribute] = infoValue

	return nil
}

// GetCartInfo returns current cart info assigned
func (it *DefaultCart) GetCartInfo() map[string]interface{} {
	return it.Info
}

func (it *DefaultCart) SetCustomInfo(info map[string]interface{}) {
	it.CustomInfo = info
}
func (it *DefaultCart) GetCustomInfo() map[string]interface{} {
	return it.CustomInfo
}

// MakeCartForVisitor loads cart for given visitor from database or creates new one
func (it *DefaultCart) MakeCartForVisitor(visitorID string) error {
	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4c6da71f-e0ea-4aae-8560-026c08bfc097", "can't get DB Engine")
	}

	cartCollection, err := dbEngine.GetCollection(ConstCartCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := cartCollection.AddFilter("visitor_id", "=", visitorID); err != nil {
		return env.ErrorDispatch(err)
	}
	if err := cartCollection.AddFilter("active", "=", true); err != nil {
		return env.ErrorDispatch(err)
	}
	rowsData, err := cartCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(rowsData) < 1 {
		newModel, err := it.New()
		if err != nil {
			return env.ErrorDispatch(err)
		}
		newCart := newModel.(*DefaultCart)
		if err := newCart.SetVisitorID(visitorID); err != nil {
			return env.ErrorDispatch(err)
		}
		if err := newCart.Activate(); err != nil {
			return env.ErrorDispatch(err)
		}
		if err := newCart.Save(); err != nil {
			return env.ErrorDispatch(err)
		}

		*it = *newCart
	} else {
		err := it.Load(rowsData[0]["_id"].(string))
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// Activate makes cart active
//   - only one cart can be active for particular visitor
func (it *DefaultCart) Activate() error {
	it.Active = true
	return nil
}

// Deactivate makes cart un-active
//   - so new cart will be created on next request
func (it *DefaultCart) Deactivate() error {
	it.Active = false
	return nil
}

// IsActive returns active flag status of cart
func (it *DefaultCart) IsActive() bool {
	return it.Active
}

// ValidateCart returns nil of cart is valid, error otherwise
func (it *DefaultCart) ValidateCart() error {
	for _, cartItem := range it.GetItems() {
		if err := cartItem.ValidateProduct(); err != nil {
			return err
		}
	}
	return nil
}

// GetSessionID returns session id last time used for cart
func (it *DefaultCart) GetSessionID() string {
	return it.SessionID
}

// SetSessionID sets session id associated to cart
func (it *DefaultCart) SetSessionID(sessionID string) error {
	it.SessionID = sessionID
	return nil
}

// GetLastUpdateTime returns cart last update time
func (it *DefaultCart) GetLastUpdateTime() time.Time {
	return it.UpdatedAt
}
