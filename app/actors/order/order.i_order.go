package order

import (
	"fmt"
	"sort"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetItems returns order items for current order
func (it *DefaultOrder) GetItems() []order.InterfaceOrderItem {
	var result []order.InterfaceOrderItem

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

// AddItem adds line item to current order, or returns error
func (it *DefaultOrder) AddItem(productID string, qty int, productOptions map[string]interface{}) (order.InterfaceOrderItem, error) {

	productModel, err := product.LoadProductByID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.ApplyOptions(productOptions)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	newOrderItem := new(DefaultOrderItem)
	newOrderItem.OrderID = it.GetID()

	err = newOrderItem.Set("product_id", productModel.GetID())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("qty", qty)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("options", productModel.GetOptions())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("name", productModel.GetName())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("sku", productModel.GetSku())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("short_description", productModel.GetShortDescription())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("price", productModel.GetPrice())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("weight", productModel.GetWeight())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	it.maxIdx++
	newOrderItem.idx = it.maxIdx
	it.Items[newOrderItem.idx] = newOrderItem

	return newOrderItem, nil
}

// RemoveItem removes line item from current order, or returns error
func (it *DefaultOrder) RemoveItem(itemIdx int) error {
	if orderItem, present := it.Items[itemIdx]; present {

		dbEngine := db.GetDBEngine()
		if dbEngine == nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "54410b67-aff0-418f-ad76-6453a2d6fed6", "can't get DB engine")
		}

		orderItemsCollection, err := dbEngine.GetCollection(ConstCollectionNameOrderItems)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = orderItemsCollection.DeleteByID(orderItem.GetID())
		if err != nil {
			return env.ErrorDispatch(err)
		}

		delete(it.Items, itemIdx)

		return nil
	}

	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1bd2f0f9-a457-43d1-a9db-e05b1aa7e1d2", "can't find index "+utils.InterfaceToString(itemIdx))
}

// NewIncrementID assigns new unique increment id to order
func (it *DefaultOrder) NewIncrementID() error {
	lastIncrementIDMutex.Lock()

	lastIncrementID++
	it.IncrementID = fmt.Sprintf(ConstIncrementIDFormat, lastIncrementID)

	env.GetConfig().SetValue(ConstConfigPathLastIncrementID, lastIncrementID)

	lastIncrementIDMutex.Unlock()

	return nil
}

// GetIncrementID returns increment id of order
func (it *DefaultOrder) GetIncrementID() string {
	return it.IncrementID
}

// SetIncrementID sets increment id to order
func (it *DefaultOrder) SetIncrementID(incrementID string) error {
	it.IncrementID = incrementID

	return nil
}

// CalculateTotals recalculates order Subtotal and GrandTotal
func (it *DefaultOrder) CalculateTotals() error {

	var subtotal float64
	for _, orderItem := range it.Items {
		subtotal += utils.RoundPrice(orderItem.GetPrice() * float64(orderItem.GetQty()))
	}
	it.Subtotal = utils.RoundPrice(subtotal)

	it.GrandTotal = utils.RoundPrice(it.Subtotal + it.ShippingAmount + it.TaxAmount - it.Discount)

	return nil
}

// GetSubtotal returns subtotal of order
func (it *DefaultOrder) GetSubtotal() float64 {
	return it.Subtotal
}

// GetGrandTotal returns grand total of order
func (it *DefaultOrder) GetGrandTotal() float64 {
	return it.GrandTotal
}

// GetDiscountAmount returns discount amount applied to order
func (it *DefaultOrder) GetDiscountAmount() float64 {
	return it.Discount
}

// GetDiscounts returns discount applied to order
func (it *DefaultOrder) GetDiscounts() []order.StructDiscount {
	return it.Discounts
}

// GetTaxAmount returns tax amount applied to order
func (it *DefaultOrder) GetTaxAmount() float64 {
	return it.TaxAmount
}

// GetTaxes returns taxes applied to order
func (it *DefaultOrder) GetTaxes() []order.StructTaxRate {
	return it.Taxes
}

// GetShippingAmount returns order shipping cost
func (it *DefaultOrder) GetShippingAmount() float64 {
	return it.ShippingAmount
}

// GetShippingMethod returns shipping method for order
func (it *DefaultOrder) GetShippingMethod() string {
	return it.ShippingMethod
}

// GetPaymentMethod returns payment method used for order
func (it *DefaultOrder) GetPaymentMethod() string {
	return it.PaymentMethod
}

// GetShippingAddress returns shipping address for order
func (it *DefaultOrder) GetShippingAddress() visitor.InterfaceVisitorAddress {
	addressModel, _ := visitor.GetVisitorAddressModel()
	addressModel.FromHashMap(it.ShippingAddress)

	return addressModel
}

// GetBillingAddress returns billing address for order
func (it *DefaultOrder) GetBillingAddress() visitor.InterfaceVisitorAddress {
	addressModel, _ := visitor.GetVisitorAddressModel()
	addressModel.FromHashMap(it.BillingAddress)

	return addressModel
}

// GetStatus returns current order status
func (it *DefaultOrder) GetStatus() string {
	return it.Status
}

// SetStatus changes status for current order
//   - if status change no supposing stock operations, order instance will not be saved automatically
func (it *DefaultOrder) SetStatus(newStatus string) error {
	var err error

	// cases with no actions
	if it.Status == newStatus || newStatus == "" {
		return nil
	}

	// changing status
	oldStatus := it.Status
	it.Status = newStatus

	// if order new status is "new" or "cancelled" - returning items to stock, otherwise taking them from
	if newStatus == order.ConstOrderStatusCancelled || newStatus == order.ConstOrderStatusNew {

		if oldStatus != order.ConstOrderStatusNew && oldStatus != order.ConstOrderStatusCancelled && oldStatus != "" {
			err = it.Rollback()
		}

	} else {

		// taking items from stock
		if oldStatus == order.ConstOrderStatusCancelled || oldStatus == order.ConstOrderStatusNew || oldStatus == "" {
			err = it.Proceed()
		}
	}

	return env.ErrorDispatch(err)
}

// Proceed subtracts order items from stock, changes status to new if status was not set yet, saves order
func (it *DefaultOrder) Proceed() error {

	if it.Status == "" {
		it.Status = order.ConstOrderStatusNew
	}

	var err error
	stockManager := product.GetRegisteredStock()
	if stockManager != nil {
		for _, orderItem := range it.GetItems() {
			options := orderItem.GetOptions()

			for optionName, optionValue := range options {
				if optionValue, ok := optionValue.(map[string]interface{}); ok {
					if value, present := optionValue["value"]; present {
						options := map[string]interface{}{optionName: value}

						err := stockManager.UpdateProductQty(orderItem.GetProductID(), options, -1*orderItem.GetQty())
						if err != nil {
							return env.ErrorDispatch(err)
						}

					}
				}
			}

		}
	}

	// checking order's incrementID, if not set - assigning new one
	if it.GetIncrementID() == "" {
		err = it.NewIncrementID()
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	err = it.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	eventData := map[string]interface{}{"order": it}
	env.Event("order.proceed", eventData)

	return nil
}

// Rollback returns order items to stock, changing order status to cancelled if status was not set yet, saves order
func (it *DefaultOrder) Rollback() error {

	if it.Status == "" {
		it.Status = order.ConstOrderStatusCancelled
	}

	var err error
	stockManager := product.GetRegisteredStock()
	if stockManager != nil {
		for _, orderItem := range it.GetItems() {
			options := orderItem.GetOptions()

			for optionName, optionValue := range options {
				if optionValue, ok := optionValue.(map[string]interface{}); ok {
					if value, present := optionValue["value"]; present {
						options := map[string]interface{}{optionName: value}

						err := stockManager.UpdateProductQty(orderItem.GetProductID(), options, orderItem.GetQty())
						if err != nil {
							return env.ErrorDispatch(err)
						}

					}
				}
			}

		}
	}

	err = it.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	eventData := map[string]interface{}{"order": it}
	env.Event("order.rollback", eventData)

	return nil
}
