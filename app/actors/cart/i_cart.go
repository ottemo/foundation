package cart

import (
	"sort"
	"strconv"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
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
					optionValuesToCheck := make([]string, 0)
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
						return env.ErrorNew("unexpected option value for '" + optionName + "' option")
					}

					// checking for option customer set with available for product
					for _, optionValue := range optionValuesToCheck {

						// there could be couple forms of product available options specification
						switch productOptionValues := productOption["options"].(type) {
						case map[string]interface{}:
							if _, present := productOptionValues[optionValue]; !present {
								return env.ErrorNew("invalid value for option '" + optionName + "'")
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
								return env.ErrorNew("invalid value for option '" + optionName + "'")
							}

						default:
							if productOptionValues != optionValue {
								return env.ErrorNew("invalid value for option '" + optionName + "'")
							}
						}
					}
				}
			}
		} else {
			return env.ErrorNew("unknown option '" + optionName + "'")
		}
	}

	// loop 2: checking that all product required options were set
	for productOption, productOptionValue := range productOptions {
		// required product option means that "productOption["required"]=true" should be set
		if productOptionValue, ok := productOptionValue.(map[string]interface{}); ok {
			if _, present := productOptionValue["required"]; present {
				if value, ok := productOptionValue["required"].(bool); ok && value {

					//checking cart item option for required option existence
					if itemOptionValue, present := cartItemOptions[productOption]; !present {
						return env.ErrorNew(productOption + " was not specified")
					} else {
						// for multi value options additional check
						switch typedValue := itemOptionValue.(type) {
						case []interface{}:
							if len(typedValue) == 0 {
								return env.ErrorNew(productOption + " was not specified")
							}
						}
					}

				}
			}
		}
	}

	return nil
}

// adds item to the current cart
//   - returns added item or nil if error happened
func (it *DefaultCart) AddItem(productId string, qty int, options map[string]interface{}) (cart.I_CartItem, error) {

	//checking qty
	if qty <= 0 {
		return nil, env.ErrorNew("qty can't be zero or less")
	}

	// checking product existence
	reqProduct, err := product.LoadProductById(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// options default value if them are not set
	if options == nil {
		options = make(map[string]interface{})
	}

	// preparing new cart item
	cartItem := &DefaultCartItem{
		idx:       it.maxIdx,
		ProductId: reqProduct.GetId(),
		Qty:       qty,
		Options:   options,
		Cart:      it,
	}

	// checking for right options
	if err := it.checkOptions(reqProduct.GetOptions(), cartItem.Options); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// adding new item to others
	it.maxIdx += 1
	cartItem.idx = it.maxIdx
	it.Items[it.maxIdx] = cartItem

	// cart change event broadcast
	it.cartChanged()

	return cartItem, nil
}

// removes item from cart
//   - you need to know index you can get from ListItems()
func (it *DefaultCart) RemoveItem(itemIdx int) error {
	if cartItem, present := it.Items[itemIdx]; present {

		dbEngine := db.GetDBEngine()
		if dbEngine == nil {
			return env.ErrorNew("can't get DB engine")
		}

		cartItemsCollection, err := dbEngine.GetCollection(CART_ITEMS_COLLECTION_NAME)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = cartItemsCollection.DeleteById(cartItem.GetId())
		if err != nil {
			return env.ErrorDispatch(err)
		}

		delete(it.Items, itemIdx)

		it.cartChanged()

		return nil
	} else {
		return env.ErrorNew("can't find index " + strconv.Itoa(itemIdx))
	}
}

// sets new qty for particular item in cart
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
	} else {
		return env.ErrorNew("there is no item with idx=" + strconv.Itoa(itemIdx))
	}
}

// returns subtotal for cart items
func (it *DefaultCart) GetSubtotal() float64 {

	if it.Subtotal == 0 {
		for _, cartItem := range it.Items {
			if cartProduct := cartItem.GetProduct(); cartProduct != nil {
				cartProduct.ApplyOptions(cartItem.GetOptions())
				it.Subtotal += cartProduct.GetPrice() * float64(cartItem.GetQty())
			}
		}
	}

	return it.Subtotal
}

// enumerates current cart items sorted by item idx
func (it *DefaultCart) GetItems() []cart.I_CartItem {

	result := make([]cart.I_CartItem, 0)

	keys := make([]int, 0)
	for key, _ := range it.Items {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {
		result = append(result, it.Items[key])
	}

	return result
}

// returns visitor id this cart belongs to
func (it *DefaultCart) GetVisitorId() string {
	return it.VisitorId
}

// sets new owner of cart
func (it *DefaultCart) SetVisitorId(visitorId string) error {
	it.VisitorId = visitorId
	return nil
}

// returns visitor model represents owner or current cart or nil if visitor was not set to cart
func (it *DefaultCart) GetVisitor() visitor.I_Visitor {
	visitor, _ := visitor.LoadVisitorByID(it.VisitorId)
	return visitor
}

// assigns some information to current cart
func (it *DefaultCart) SetCartInfo(infoAttribute string, infoValue interface{}) error {
	if it.Info == nil {
		it.Info = make(map[string]interface{})
	}

	it.Info[infoAttribute] = infoValue

	return nil
}

// returns current cart info assigned
func (it *DefaultCart) GetCartInfo() map[string]interface{} {
	return it.Info
}

// loads cart information from DB for visitor
func (it *DefaultCart) MakeCartForVisitor(visitorId string) error {
	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("can't get DB Engine")
	}

	cartCollection, err := dbEngine.GetCollection(CART_COLLECTION_NAME)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	cartCollection.AddFilter("visitor_id", "=", visitorId)
	cartCollection.AddFilter("active", "=", true)
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
		newCart.SetVisitorId(visitorId)
		newCart.Activate()
		newCart.Save()

		*it = *newCart
	} else {
		err := it.Load(rowsData[0]["_id"].(string))
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// makes cart active
//   - only one cart can be active for particular visitor
func (it *DefaultCart) Activate() error {
	it.Active = true
	return nil
}

// makes cart un-active
//   - so new cart will be created on next request
func (it *DefaultCart) Deactivate() error {
	it.Active = false
	return nil
}

// returns active flag status of cart
func (it *DefaultCart) IsActive() bool {
	return it.Active
}
