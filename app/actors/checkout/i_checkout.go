package checkout

import (
	"fmt"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)

// SetShippingAddress sets shipping address for checkout
func (it *DefaultCheckout) SetShippingAddress(address visitor.I_VisitorAddress) error {
	it.ShippingAddressId = address.GetId()
	return nil
}

// GetShippingAddress returns checkout shipping address
func (it *DefaultCheckout) GetShippingAddress() visitor.I_VisitorAddress {
	shippingAddress, _ := visitor.LoadVisitorAddressByID(it.ShippingAddressId)
	return shippingAddress
}

// SetBillingAddress sets billing address for checkout
func (it *DefaultCheckout) SetBillingAddress(address visitor.I_VisitorAddress) error {
	it.BillingAddressId = address.GetId()
	return nil
}

// GetBillingAddress returns checkout billing address
func (it *DefaultCheckout) GetBillingAddress() visitor.I_VisitorAddress {
	billingAddress, _ := visitor.LoadVisitorAddressByID(it.BillingAddressId)
	return billingAddress
}

// SetPaymentMethod sets payment method for checkout
func (it *DefaultCheckout) SetPaymentMethod(paymentMethod checkout.I_PaymentMethod) error {
	it.PaymentMethodCode = paymentMethod.GetCode()
	return nil
}

// GetPaymentMethod returns checkout payment method
func (it *DefaultCheckout) GetPaymentMethod() checkout.I_PaymentMethod {
	if paymentMethods := checkout.GetRegisteredPaymentMethods(); paymentMethods != nil {
		for _, paymentMethod := range checkout.PaymentMethods {
			if paymentMethod.GetCode() == it.PaymentMethodCode {
				return paymentMethod
			}
		}
	}
	return nil
}

// SetShippingMethod sets payment method for checkout
func (it *DefaultCheckout) SetShippingMethod(shippingMethod checkout.I_ShippingMehod) error {
	it.ShippingMethodCode = shippingMethod.GetCode()
	return nil
}

// GetShippingMethod returns a checkout shipping method
func (it *DefaultCheckout) GetShippingMethod() checkout.I_ShippingMehod {
	if shippingMethods := checkout.GetRegisteredShippingMethods(); shippingMethods != nil {
		for _, shippingMethod := range shippingMethods {
			if shippingMethod.GetCode() == it.ShippingMethodCode {
				return shippingMethod
			}
		}
	}
	return nil
}

// SetShippingRate sets shipping rate for checkout
func (it *DefaultCheckout) SetShippingRate(shippingRate checkout.T_ShippingRate) error {
	it.ShippingRate = shippingRate
	return nil
}

// GetShippingRate returns a checkout shipping rate
func (it *DefaultCheckout) GetShippingRate() *checkout.T_ShippingRate {
	return &it.ShippingRate
}

// SetCart sets cart for checkout
func (it *DefaultCheckout) SetCart(checkoutCart cart.I_Cart) error {
	it.CartId = checkoutCart.GetId()
	return nil
}

// GetCart returns a shopping cart
func (it *DefaultCheckout) GetCart() cart.I_Cart {
	cartInstance, _ := cart.LoadCartById(it.CartId)
	return cartInstance
}

// SetVisitor sets visitor for checkout
func (it *DefaultCheckout) SetVisitor(checkoutVisitor visitor.I_Visitor) error {
	it.VisitorId = checkoutVisitor.GetId()

	if it.BillingAddressId == "" && checkoutVisitor.GetBillingAddress() != nil {
		it.BillingAddressId = checkoutVisitor.GetBillingAddress().GetId()
	}

	if it.ShippingAddressId == "" && checkoutVisitor.GetShippingAddress() != nil {
		it.ShippingAddressId = checkoutVisitor.GetShippingAddress().GetId()
	}

	return nil
}

// GetVisitor return checkout visitor
func (it *DefaultCheckout) GetVisitor() visitor.I_Visitor {
	visitorInstance, _ := visitor.LoadVisitorByID(it.VisitorId)
	return visitorInstance
}

// SetSession sets visitor for checkout
func (it *DefaultCheckout) SetSession(checkoutSession api.I_Session) error {
	it.SessionId = checkoutSession.GetId()
	return nil
}

// GetSession return checkout visitor
func (it *DefaultCheckout) GetSession() api.I_Session {
	return api.GetSessionById(it.SessionId)
}

// GetTaxes collects taxes applied for current checkout
func (it *DefaultCheckout) GetTaxes() (float64, []checkout.T_TaxRate) {

	var amount float64

	if !it.taxesCalculateFlag {
		it.taxesCalculateFlag = true

		it.Taxes = make([]checkout.T_TaxRate, 0)
		for _, tax := range checkout.GetRegisteredTaxes() {
			for _, taxRate := range tax.CalculateTax(it) {
				it.Taxes = append(it.Taxes, taxRate)
				amount += taxRate.Amount
			}
		}

		it.taxesCalculateFlag = false
	} else {
		for _, taxRate := range it.Taxes {
			amount += taxRate.Amount
		}
	}

	return amount, it.Taxes
}

// GetDiscounts collects discounts applied for current checkout
func (it *DefaultCheckout) GetDiscounts() (float64, []checkout.T_Discount) {

	var amount float64

	if !it.discountsCalculateFlag {
		it.discountsCalculateFlag = true

		it.Discounts = make([]checkout.T_Discount, 0)
		for _, discount := range checkout.GetRegisteredDiscounts() {
			for _, discountValue := range discount.CalculateDiscount(it) {
				it.Discounts = append(it.Discounts, discountValue)
				amount += discountValue.Amount
			}
		}

		it.discountsCalculateFlag = false
	} else {
		for _, discount := range it.Discounts {
			amount += discount.Amount
		}
	}

	return amount, it.Discounts
}

// GetGrandTotal return grand total for current checkout: [cart subtotal] + [shipping rate] + [taxes] - [discounts]
func (it *DefaultCheckout) GetGrandTotal() float64 {
	var amount float64

	currentCart := it.GetCart()
	if currentCart != nil {
		amount += currentCart.GetSubtotal()
	}

	if shippingRate := it.GetShippingRate(); shippingRate != nil {
		amount += shippingRate.Price
	}

	taxAmount, _ := it.GetTaxes()
	amount += taxAmount

	discountAmount, _ := it.GetDiscounts()
	amount -= discountAmount

	return amount
}

// SetInfo sets additional info for checkout - any values related to checkout process
func (it *DefaultCheckout) SetInfo(key string, value interface{}) error {
	it.Info[key] = value

	return nil
}

// GetInfo returns additional checkout info value or nil
func (it *DefaultCheckout) GetInfo(key string) interface{} {
	if value, present := it.Info[key]; present {
		return value
	}
	return nil
}

// SetOrder sets order for current checkout
func (it *DefaultCheckout) SetOrder(checkoutOrder order.I_Order) error {
	it.OrderId = checkoutOrder.GetId()
	return nil
}

// GetOrder returns current checkout related order or nil if not created yet
func (it *DefaultCheckout) GetOrder() order.I_Order {
	if it.OrderId != "" {
		orderInstance, err := order.LoadOrderById(it.OrderId)
		if err == nil {
			return orderInstance
		}
	}
	return nil
}

// Submit creates the order with provided information
func (it *DefaultCheckout) Submit() (interface{}, error) {

	if it.GetBillingAddress() == nil {
		return nil, env.ErrorNew("Billing address is not set")
	}

	if it.GetShippingAddress() == nil {
		return nil, env.ErrorNew("Shipping address is not set")
	}

	if it.GetPaymentMethod() == nil {
		return nil, env.ErrorNew("Payment method is not set")
	}

	if it.GetShippingMethod() == nil {
		return nil, env.ErrorNew("Shipping method is not set")
	}

	currentCart := it.GetCart()
	if currentCart == nil {
		return nil, env.ErrorNew("Cart is not specified")
	}

	cartItems := currentCart.GetItems()
	if len(cartItems) == 0 {
		return nil, env.ErrorNew("Cart have no products inside")
	}

	// making new order if needed
	//---------------------------
	currentTime := time.Now()

	checkoutOrder := it.GetOrder()
	if checkoutOrder == nil {
		newOrder, err := order.GetOrderModel()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		newOrder.Set("created_at", currentTime)

		checkoutOrder = newOrder
	}

	// updating order information
	//---------------------------
	checkoutOrder.Set("updated_at", currentTime)

	checkoutOrder.Set("status", "new")
	if currentVisitor := it.GetVisitor(); currentVisitor != nil {
		checkoutOrder.Set("visitor_id", currentVisitor.GetId())

		checkoutOrder.Set("customer_email", currentVisitor.GetEmail())
		checkoutOrder.Set("customer_name", currentVisitor.GetFullName())
	}

	billingAddress := it.GetBillingAddress().ToHashMap()
	checkoutOrder.Set("billing_address", billingAddress)

	shippingAddress := it.GetShippingAddress().ToHashMap()
	checkoutOrder.Set("shipping_address", shippingAddress)

	checkoutOrder.Set("cart_id", currentCart.GetId())
	checkoutOrder.Set("payment_method", it.GetPaymentMethod().GetCode())
	checkoutOrder.Set("shipping_method", it.GetShippingMethod().GetCode()+"/"+it.GetShippingRate().Code)

	discountAmount, _ := it.GetDiscounts()
	taxAmount, _ := it.GetTaxes()

	checkoutOrder.Set("discount", discountAmount)
	checkoutOrder.Set("tax_amount", taxAmount)
	checkoutOrder.Set("shipping_amount", it.GetShippingRate().Price)

	generateDescriptionFlag := false
	orderDescription := utils.InterfaceToString(it.GetInfo("order_description"))
	if orderDescription == "" {
		generateDescriptionFlag = true
	}

	for _, cartItem := range cartItems {
		orderItem, err := checkoutOrder.AddItem(cartItem.GetProductId(), cartItem.GetQty(), cartItem.GetOptions())
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if generateDescriptionFlag {
			if orderDescription != "" {
				orderDescription += ", "
			}
			orderDescription += fmt.Sprintf("%dx %s", cartItem.GetQty(), orderItem.GetName())
		}
	}
	checkoutOrder.Set("description", orderDescription)

	err := checkoutOrder.CalculateTotals()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = checkoutOrder.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	it.SetOrder(checkoutOrder)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// trying to process payment
	//--------------------------
	paymentInfo := make(map[string]interface{})
	paymentInfo["sessionId"] = it.GetSession().GetId()

	result, err := it.GetPaymentMethod().Authorize(checkoutOrder, paymentInfo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// TODO: should be different way to do that, result should be some interface or just error but not this
	if result != nil {
		return result, nil
	}

	// assigning new order increment id after success payment
	//-------------------------------------------------------
	checkoutOrder.NewIncrementId()

	checkoutOrder.Set("status", "pending")

	err = it.CheckoutSuccess(checkoutOrder, it.GetSession())
	if err != nil {
		return nil, err
	}

	err = it.SendOrderConfirmationMail()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// cleanup checkout information
	//-----------------------------
	currentCart.Deactivate()
	currentCart.Save()

	it.GetSession().Set(cart.SESSION_KEY_CURRENT_CART, nil)
	it.GetSession().Set(checkout.SESSION_KEY_CURRENT_CHECKOUT, nil)

	return checkoutOrder.ToHashMap(), nil
}
