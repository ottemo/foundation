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
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "54410b67aff0418fad766453a2d6fed6", "can't get DB engine")
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

	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1bd2f0f9a45743d1a9dbe05b1aa7e1d2", "can't find index "+utils.InterfaceToString(itemIdx))
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

// GetTaxAmount returns tax amount applied to order
func (it *DefaultOrder) GetTaxAmount() float64 {
	return it.TaxAmount
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
func (it *DefaultOrder) SetStatus(status string) error {
	var err error

	if it.Status == status {
		return nil
	}

	switch status {
	case order.ConstOrderStatusNew:
		err = it.Proceed()
	case order.ConstOrderStatusCancelled:
		err = it.Cancel()
	default:
		it.Status = status
	}

	return err
}

// Proceed subtracts order items from stock, changes status to new, saves order
func (it *DefaultOrder) Proceed() error {

	it.Status = order.ConstOrderStatusNew

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

	err = it.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Cancel returns order items to stock and changing order status to canceled, saves order
func (it *DefaultOrder) Cancel() error {
	it.Status = order.ConstOrderStatusCancelled

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

	return nil
}
