package checkout

import (
	"fmt"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/app/models/visitor"
)

// SetShippingAddress sets shipping address for checkout
func (it *DefaultCheckout) SetShippingAddress(address visitor.InterfaceVisitorAddress) error {
	if address == nil {
		it.ShippingAddress = nil
		return nil
	}

	it.ShippingAddress = address.ToHashMap()
	return nil
}

// GetShippingAddress returns checkout shipping address
func (it *DefaultCheckout) GetShippingAddress() visitor.InterfaceVisitorAddress {
	if it.ShippingAddress == nil {
		return nil
	}

	shippingAddress, err := visitor.GetVisitorAddressModel()
	if err != nil {
		_ = env.ErrorDispatch(err)
		return nil
	}

	err = shippingAddress.FromHashMap(it.ShippingAddress)
	if err != nil {
		_ = env.ErrorDispatch(err)
		return nil
	}

	return shippingAddress
}

// SetBillingAddress sets billing address for checkout
func (it *DefaultCheckout) SetBillingAddress(address visitor.InterfaceVisitorAddress) error {
	if address == nil {
		it.BillingAddress = nil
		return nil
	}

	it.BillingAddress = address.ToHashMap()
	return nil
}

// GetBillingAddress returns checkout billing address
func (it *DefaultCheckout) GetBillingAddress() visitor.InterfaceVisitorAddress {
	if it.BillingAddress == nil {
		return nil
	}

	billingAddress, err := visitor.GetVisitorAddressModel()
	if err != nil {
		_ = env.ErrorDispatch(err)
		return nil
	}

	err = billingAddress.FromHashMap(it.BillingAddress)
	if err != nil {
		_ = env.ErrorDispatch(err)
		return nil
	}

	return billingAddress
}

// SetPaymentMethod sets payment method for checkout
func (it *DefaultCheckout) SetPaymentMethod(paymentMethod checkout.InterfacePaymentMethod) error {
	it.PaymentMethodCode = paymentMethod.GetCode()
	return nil
}

// GetPaymentMethod returns checkout payment method
func (it *DefaultCheckout) GetPaymentMethod() checkout.InterfacePaymentMethod {
	if paymentMethods := checkout.GetRegisteredPaymentMethods(); paymentMethods != nil {
		for _, paymentMethod := range paymentMethods {
			if paymentMethod.GetCode() == it.PaymentMethodCode {
				return paymentMethod
			}
		}
	}
	return nil
}

// SetShippingMethod sets payment method for checkout
func (it *DefaultCheckout) SetShippingMethod(shippingMethod checkout.InterfaceShippingMethod) error {
	it.ShippingMethodCode = shippingMethod.GetCode()
	return nil
}

// GetShippingMethod returns a checkout shipping method
func (it *DefaultCheckout) GetShippingMethod() checkout.InterfaceShippingMethod {
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
func (it *DefaultCheckout) SetShippingRate(shippingRate checkout.StructShippingRate) error {
	it.ShippingRate = shippingRate
	return nil
}

// GetShippingRate returns a checkout shipping rate
func (it *DefaultCheckout) GetShippingRate() *checkout.StructShippingRate {
	return &it.ShippingRate
}

// SetCart sets cart for checkout
func (it *DefaultCheckout) SetCart(checkoutCart cart.InterfaceCart) error {
	if checkoutCart != nil {
		it.CartID = checkoutCart.GetID()
		it.cart = checkoutCart
	} else {
		it.CartID = ""
	}

	return nil
}

// GetCart returns a shopping cart
func (it *DefaultCheckout) GetCart() cart.InterfaceCart {
	if it.cart != nil {
		return it.cart
	}

	if it.CartID != "" {
		if cartInstance, err := cart.LoadCartByID(it.CartID); err == nil {
			it.cart = cartInstance
			return it.cart
		}
	}

	return nil
}

// SetVisitor sets visitor for checkout
func (it *DefaultCheckout) SetVisitor(checkoutVisitor visitor.InterfaceVisitor) error {
	it.VisitorID = checkoutVisitor.GetID()

	if it.BillingAddress == nil && checkoutVisitor.GetBillingAddress() != nil {
		it.BillingAddress = checkoutVisitor.GetBillingAddress().ToHashMap()
	}

	if it.ShippingAddress == nil && checkoutVisitor.GetShippingAddress() != nil {
		it.ShippingAddress = checkoutVisitor.GetShippingAddress().ToHashMap()
	}

	return nil
}

// GetVisitor return checkout visitor
func (it *DefaultCheckout) GetVisitor() visitor.InterfaceVisitor {
	if it.VisitorID == "" {
		return nil
	}

	visitorInstance, err := visitor.LoadVisitorByID(it.VisitorID)
	if err != nil {
		return nil
	}

	return visitorInstance
}

// SetSession sets visitor for checkout
func (it *DefaultCheckout) SetSession(checkoutSession api.InterfaceSession) error {
	it.SessionID = checkoutSession.GetID()
	return nil
}

// GetSession return checkout session
func (it *DefaultCheckout) GetSession() api.InterfaceSession {
	sessionInstance, _ := api.GetSessionByID(it.SessionID, true)
	return sessionInstance
}

// GetTaxAmount returns total amount of taxes for current checkout
func (it *DefaultCheckout) GetTaxAmount() float64 {
	return it.GetItemSpecificTotal(0, checkout.ConstLabelTax)
}

// GetDiscountAmount returns total amount of discounts applied for current checkout
func (it *DefaultCheckout) GetDiscountAmount() float64 {
	return it.GetItemSpecificTotal(0, checkout.ConstLabelDiscount) +
		it.GetItemSpecificTotal(0, checkout.ConstLabelGiftCard) +
		it.GetItemSpecificTotal(0, checkout.ConstLabelSalePriceAdjustment)
}

// GetPriceAdjustments collects price adjustments applied for current checkout
func (it *DefaultCheckout) GetPriceAdjustments(label string) []checkout.StructPriceAdjustment {
	it.GetItemSpecificTotal(0, label) // this function will do initial calculation of checkout if it wasn't done
	var result []checkout.StructPriceAdjustment

	for _, priceAdjustment := range it.priceAdjustments {
		if label == "" || utils.IsInListStr(label, priceAdjustment.Labels) {
			result = append(result, priceAdjustment)
		}
	}

	return result
}

// GetTaxes collects taxes applied for current checkout
func (it *DefaultCheckout) GetTaxes() []checkout.StructPriceAdjustment {

	return it.GetPriceAdjustments(checkout.ConstLabelTax)
}

// GetDiscounts collects discounts applied for current checkout
func (it *DefaultCheckout) GetDiscounts() []checkout.StructPriceAdjustment {

	// currently we will include GC to discount
	result := it.GetPriceAdjustments(checkout.ConstLabelDiscount)
	for _, giftCardPA := range it.GetPriceAdjustments(checkout.ConstLabelGiftCard) {
		result = append(result, giftCardPA)
	}
	for _, salePricePA := range it.GetPriceAdjustments(checkout.ConstLabelSalePriceAdjustment) {
		result = append(result, salePricePA)
	}

	return result
}

// GetSubtotal returns subtotal total for current checkout
func (it *DefaultCheckout) GetSubtotal() float64 {
	return it.GetItemSpecificTotal(0, checkout.ConstLabelSubtotal)
}

// GetItems returns current cart items
func (it *DefaultCheckout) GetItems() []cart.InterfaceCartItem {
	if currentCart := it.GetCart(); currentCart != nil {
		return currentCart.GetItems()
	}

	return nil
}

// GetDiscountableItems returns current cart items that can be discounted (not a gift cards)
func (it *DefaultCheckout) GetDiscountableItems() []cart.InterfaceCartItem {
	if items := it.GetItems(); items != nil {
		var result []cart.InterfaceCartItem
		for _, item := range items {
			// this method should be updated to general product type usage
			if !strings.Contains(item.GetProduct().GetSku(), checkout.GiftCardSkuElement) {
				result = append(result, item)
			}
		}
		return result
	}
	return nil
}

// GetShippingAmount returns shipping price for current checkout
func (it *DefaultCheckout) GetShippingAmount() float64 {
	return it.GetItemSpecificTotal(0, checkout.ConstLabelShipping)
}

// calculateSubtotal it's an element of calculation that provides subtotal amounts
func (it *DefaultCheckout) calculateSubtotal() checkout.StructPriceAdjustment {

	items := it.GetItems()
	result := checkout.StructPriceAdjustment{
		Code:      checkout.ConstLabelSubtotal,
		Name:      checkout.ConstLabelSubtotal,
		Amount:    0,
		IsPercent: false,
		Priority:  checkout.ConstCalculateTargetSubtotal,
		Labels:    []string{checkout.ConstLabelSubtotal},
		PerItem:   map[string]float64{},
	}

	for _, cartItem := range items {
		if cartProduct := cartItem.GetProduct(); cartProduct != nil {
			result.PerItem[utils.InterfaceToString(cartItem.GetIdx())] = utils.RoundPrice(cartProduct.GetPrice() * float64(cartItem.GetQty()))
		}
	}

	return result
}

// calculateShipping it's an element of calculation that provides shipping amounts
func (it *DefaultCheckout) calculateShipping() checkout.StructPriceAdjustment {

	if shippingRate := it.GetShippingRate(); shippingRate != nil {
		return checkout.StructPriceAdjustment{
			Code:      shippingRate.Code,
			Name:      shippingRate.Name,
			Amount:    shippingRate.Price,
			IsPercent: false,
			Priority:  checkout.ConstCalculateTargetShipping,
			Labels:    []string{checkout.ConstLabelShipping},
			PerItem:   nil,
		}
	}

	return checkout.StructPriceAdjustment{}
}

// GetItemTotals return details about totals per item (0 is a cart)
func (it *DefaultCheckout) getItemTotals(idx interface{}) map[string]float64 {
	index := utils.InterfaceToInt(idx)

	// in this case we don't started calculation process so it will be executed first
	if it.calculationDetailTotals == nil {
		it.CalculateAmount(0)
	}

	itemTotals, present := it.calculationDetailTotals[index]
	if !present {

		itemTotals = make(map[string]float64)
		it.calculationDetailTotals[index] = itemTotals
	}

	return itemTotals
}

// GetItemSpecificTotal return current amount value for given item index and label (0 is a cart)
func (it *DefaultCheckout) GetItemSpecificTotal(idx interface{}, label string) float64 {
	if value, present := it.getItemTotals(idx)[label]; present {
		return value
	}

	return 0
}

// applyAmount applies amounts to checkout detail calculation map
func (it *DefaultCheckout) applyAmount(idx interface{}, label string, amount float64) {
	index := utils.InterfaceToInt(idx)
	if index != 0 {
		it.getItemTotals(index)[label] = utils.RoundPrice(amount + it.GetItemSpecificTotal(index, label))
	}

	it.getItemTotals(0)[label] = utils.RoundPrice(amount + it.GetItemSpecificTotal(0, label))

	if label == checkout.ConstLabelGrandTotal {
		it.calculateAmount = utils.RoundPrice(amount + it.calculateAmount)
	}
}

// applyPriceAdjustment used to handle calculation of changes from price adjustment
// and storing to all points with details
func (it *DefaultCheckout) applyPriceAdjustment(priceAdjustment checkout.StructPriceAdjustment) {
	if priceAdjustment.Code == "" {
		return
	}
	var totalPriceAdjustmentAmount float64
	// main part is per items apply (we will handle Amount only if there was no per item value)
	if priceAdjustment.PerItem == nil || len(priceAdjustment.PerItem) == 0 {
		amount := priceAdjustment.Amount
		if priceAdjustment.IsPercent {
			// current grand total will be changed on some percentage
			amount = it.GetItemSpecificTotal(0, checkout.ConstLabelGrandTotal) * amount / 100.0
		}

		// prevent negative values of grand total
		if amount+it.calculateAmount < 0.0 {
			amount = 0.0 - it.calculateAmount
		}

		// affecting grand total of a cart
		amount = utils.RoundPrice(amount)

		it.applyAmount(0, checkout.ConstLabelGrandTotal, amount)
		totalPriceAdjustmentAmount += amount

		// cart details show amount was applied by Types
		for _, label := range priceAdjustment.Labels {
			if label != checkout.ConstLabelGrandTotal {
				it.applyAmount(0, label, amount)
			}
		}

	} else {
		for index, amount := range priceAdjustment.PerItem {
			currentItemTotal := it.GetItemSpecificTotal(index, checkout.ConstLabelGrandTotal)
			if priceAdjustment.IsPercent {
				amount = currentItemTotal * amount / 100.0
			}

			// prevent negative values of grand total per cart
			if amount+it.calculateAmount < 0.0 {
				amount = 0.0 - it.calculateAmount
			}

			// prevent negative values of grand total per item
			if amount+currentItemTotal < 0.0 {
				amount = 0.0 - currentItemTotal
			}

			// adding amount to grand total of current item and full cart
			amount = utils.RoundPrice(amount)
			it.applyAmount(index, checkout.ConstLabelGrandTotal, amount)
			totalPriceAdjustmentAmount += amount

			for _, label := range priceAdjustment.Labels {
				if label != checkout.ConstLabelGrandTotal {
					it.applyAmount(index, label, amount)
				}
			}
		}

	}

	priceAdjustment.Amount = totalPriceAdjustmentAmount
	it.priceAdjustments = append(it.priceAdjustments, priceAdjustment)
}

// CalculateAmount do a calculation of all amounts for checkout
// TODO: make function use calculateTarget as a limit for priority to where it need to be calculated
func (it *DefaultCheckout) CalculateAmount(calculateTarget float64) float64 {

	if !it.calculateFlag {
		it.calculateFlag = true
		it.calculateAmount = 0
		it.priceAdjustments = make([]checkout.StructPriceAdjustment, 0)
		it.calculationDetailTotals = make(map[int]map[string]float64)

		var priceAdjustments []checkout.StructPriceAdjustment
		priceAdjustmentCalls := make(map[float64]func(checkout.InterfaceCheckout, float64) []checkout.StructPriceAdjustment)
		for _, priceAdjustment := range checkout.GetRegisteredPriceAdjustments() {
			for _, priorityValue := range priceAdjustment.GetPriority() {
				priceAdjustmentCalls[priorityValue] = priceAdjustment.Calculate
			}
		}

		basePoints := map[float64]func() checkout.StructPriceAdjustment{
			checkout.ConstCalculateTargetSubtotal: func() checkout.StructPriceAdjustment {
				return it.calculateSubtotal()
			},
			checkout.ConstCalculateTargetShipping: func() checkout.StructPriceAdjustment {
				return it.calculateShipping()
			},
		}

		var minPriority float64
		var maxPriority float64

		minIsSet := false
		maxIsSet := false

		searchMode := true

		// 2 cycle calculation loop
		// 1st loop - search mode, looks for current minimal priority to apply
		// 2nd loop - apply current priority items
		for searchMode || maxIsSet {

			// setting previousPriority since 2nd search
			if searchMode && maxIsSet {
				minPriority = maxPriority
				minIsSet = true
				maxIsSet = false
			}

			// base points lookup (subtotal, shipping)
			for priority, value := range basePoints {

				if searchMode {
					if (!maxIsSet || priority < maxPriority) && (!minIsSet || priority > minPriority) {
						maxPriority = priority
						maxIsSet = true
					}
				} else {
					if priority == maxPriority {
						it.applyPriceAdjustment(value())
					}
				}
			}

			// priceAdjustment calls lookup
			for priority, priceAdjustmentCall := range priceAdjustmentCalls {

				if searchMode {
					if (!maxIsSet || priority < maxPriority) && (!minIsSet || priority > minPriority) {
						maxPriority = priority
						maxIsSet = true
					}
				} else {
					if priority == maxPriority {
						for _, priceAdjustment := range priceAdjustmentCall(it, priority) {
							priceAdjustments = append(priceAdjustments, priceAdjustment)
						}
					}
				}
			}

			// priceAdjustment lookup
			for _, priceAdjustment := range priceAdjustments {

				if searchMode {
					priority := priceAdjustment.Priority
					if (!maxIsSet || priority < maxPriority) && (!minIsSet || priority > minPriority) {
						maxPriority = priceAdjustment.Priority
						maxIsSet = true
					}
				} else {
					if priceAdjustment.Priority == maxPriority {
						it.applyPriceAdjustment(priceAdjustment)
					}
				}
			}

			// cycle mode switcher
			if searchMode {
				searchMode = false
			} else {
				searchMode = true
			}
		}

		infoDetails := map[string]interface{}{}
		for index, details := range it.calculationDetailTotals {
			infoDetails[utils.InterfaceToString(index)] = details
		}

		if err := it.SetInfo("calculation", infoDetails); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3f7a3ad2-724b-4a6c-b62d-7d3eee5df1e8", err.Error())
		}
		if err := it.SetInfo("price_adjustments", it.priceAdjustments); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a1e6050e-4f10-4525-bea2-8b1fd6ed4abc", err.Error())
		}

		it.calculateFlag = false
	}

	return it.calculateAmount
}

// GetGrandTotal returns grand total for current checkout
func (it *DefaultCheckout) GetGrandTotal() float64 {
	return it.CalculateAmount(0)
}

// SetInfo sets additional info for checkout - any values related to checkout process
func (it *DefaultCheckout) SetInfo(key string, value interface{}) error {
	if value == nil {
		if _, present := it.Info[key]; present {
			delete(it.Info, key)
		}
	} else {
		it.Info[key] = value
	}

	return nil
}

// GetInfo returns additional checkout info value or nil,
//   - use "*" as a key to get all keys and values currently set
func (it *DefaultCheckout) GetInfo(key string) interface{} {
	if key == "*" {
		return it.Info
	}

	if value, present := it.Info[key]; present {
		return value
	}

	return nil
}

// SetOrder sets order for current checkout
func (it *DefaultCheckout) SetOrder(checkoutOrder order.InterfaceOrder) error {
	it.OrderID = checkoutOrder.GetID()
	return nil
}

// GetOrder returns current checkout related order or nil if not created yet
func (it *DefaultCheckout) GetOrder() order.InterfaceOrder {
	if it.OrderID != "" {
		orderInstance, err := order.LoadOrderByID(it.OrderID)
		if err == nil {
			return orderInstance
		}
	}
	return nil
}

// IsSubscription provide a flag if there subscription items in it
func (it *DefaultCheckout) IsSubscription() bool {
	return subscription.ContainsSubscriptionItems(it)
}

// Submit creates the order with provided information
func (it *DefaultCheckout) Submit() (interface{}, error) {

	if it.GetBillingAddress() == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "080db3c0-dbb5-4398-b1f1-4c3fefef79b4", "Billing address is not set")
	}

	if it.GetShippingAddress() == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1c069d54-2847-46cb-bccd-76fc13d229ea", "Shipping address is not set")
	}

	if it.GetPaymentMethod() == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c0628038-3e06-47e9-9252-480351d903c0", "Payment method is not set")
	}

	if it.GetShippingMethod() == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e22877fe-248d-4b5e-ad2f-10843cb9890c", "Shipping method is not set")
	}

	currentCart := it.GetCart()
	if currentCart == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ea211827-2025-46ca-be47-d09841021890", "Cart is not specified")
	}

	cartItems := currentCart.GetItems()
	if len(cartItems) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "84ad4db5-29e9-430c-aecf-b675cbafbcec", "Cart is empty")
	}

	currentVisitor := it.GetVisitor()
	if (it.VisitorID == "" || currentVisitor == nil) && utils.InterfaceToString(it.GetInfo("customer_email")) == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c5f53ede-63b7-40ea-952d-4d4c04337563", "customer e-mail was not specified")
	}

	if err := currentCart.ValidateCart(); err != nil {
		return nil, env.ErrorDispatch(err)
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

		if err := newOrder.Set("created_at", currentTime); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8718b96a-3a5f-4ac2-95ad-e19cabcedda9", err.Error())
		}

		checkoutOrder = newOrder
	}

	// updating order information
	//---------------------------
	if err := checkoutOrder.Set("session_id", it.GetInfo("session_id")); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "92434482-8676-4b0c-be94-710f98f70a60", err.Error())
	}
	if err := checkoutOrder.Set("updated_at", currentTime); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bdece7e9-8667-48b1-8783-9d6d0c29f06e", err.Error())
	}

	if err := checkoutOrder.SetStatus(order.ConstOrderStatusNew); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7b5d801d-a553-4e18-a8f4-520e05f138d1", err.Error())
	}

	// If the visitor is logged in that should dictate the email and name
	if currentVisitor != nil {
		if err := checkoutOrder.Set("visitor_id", currentVisitor.GetID()); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "85d8fe53-4853-428a-9dfa-3ad594647472", err.Error())
		}

		if err := checkoutOrder.Set("customer_email", currentVisitor.GetEmail()); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b426f813-d607-4dfc-b7cd-e780945ea2a2", err.Error())
		}
		if err := checkoutOrder.Set("customer_name", currentVisitor.GetFullName()); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a5e7127a-9c09-492e-9e4f-684031f9cba7", err.Error())
		}
	} else {
		// Visitor is not logged in, maybe we were passed some extra data
		if it.GetInfo("customer_email") != nil {
			orderCustomerEmail := utils.InterfaceToString(it.GetInfo("customer_email"))
			if err := checkoutOrder.Set("customer_email", orderCustomerEmail); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0b8126f1-3b33-4f1a-8cfe-9ddcc915b290", err.Error())
			}
		}
		if it.GetInfo("customer_name") != nil {
			orderCustomerName := utils.InterfaceToString(it.GetInfo("customer_name"))
			if err := checkoutOrder.Set("customer_name", orderCustomerName); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ffe69f29-27b6-4b53-820f-9ba7ca6570ae", err.Error())
			}
		}
	}

	billingAddress := it.GetBillingAddress().ToHashMap()
	if err := checkoutOrder.Set("billing_address", billingAddress); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c7f3cae2-ed76-4002-8d7c-5031cf55c6c6", err.Error())
	}

	shippingAddress := it.GetShippingAddress().ToHashMap()
	if err := checkoutOrder.Set("shipping_address", shippingAddress); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c005d525-4c06-4d65-933b-6a928dd93bac", err.Error())
	}

	shippingInfo := utils.InterfaceToMap(checkoutOrder.Get("shipping_info"))
	shippingInfo["shipping_method_name"] = it.GetShippingMethod().GetName() + "/" + it.GetShippingRate().Name
	if notes := utils.InterfaceToString(it.GetInfo("notes")); notes != "" {
		shippingInfo["notes"] = notes
	}
	if err := checkoutOrder.Set("shipping_info", shippingInfo); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e6de55fb-4a18-4aee-b474-1f0ecf742789", err.Error())
	}
	if err := checkoutOrder.Set("shipping_method", it.GetShippingMethod().GetCode()+"/"+it.GetShippingRate().Code); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "78ef6507-709c-4aa9-a3ac-d13364128b42", err.Error())
	}

	if err := checkoutOrder.Set("cart_id", currentCart.GetID()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "792a8468-4e49-43f6-9670-8f1c2fa122c0", err.Error())
	}

	paymentMethod := it.GetPaymentMethod()

	// call for recalculating of all amounts including taxes and discounts
	// cause they can be not calculated at this point in case of direct post
	it.CalculateAmount(checkout.ConstCalculateTargetGrandTotal)

	customInfo := utils.InterfaceToMap(checkoutOrder.Get("custom_info"))
	customInfo["calculation"] = it.Info["calculation"]
	if err := checkoutOrder.Set("custom_info", customInfo); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "937c4411-c18d-4cf2-a0de-17c0ccb9d42b", err.Error())
	}

	if notes := utils.InterfaceToArray(it.GetInfo("order_notes")); len(notes) > 0 {
		if err := checkoutOrder.Set("notes", notes); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ebc2b8ad-914f-4922-b5b1-97fe237a840e", err.Error())
		}
	}

	if !paymentMethod.IsAllowed(it) {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7a5490ee-daa3-42b4-a84a-dade12d103e8", "Payment method not allowed")
	}

	if err := checkoutOrder.Set("payment_method", paymentMethod.GetCode()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "275dbdd6-7721-4176-bea0-1d208fed88c4", err.Error())
	}
	paymentInfo := utils.InterfaceToMap(checkoutOrder.Get("payment_info"))
	paymentInfo["payment_method_name"] = it.GetPaymentMethod().GetInternalName()
	paymentInfo["gift_cards_charged_amount"] = it.GetItemSpecificTotal(0, checkout.ConstLabelGiftCard)
	if err := checkoutOrder.Set("payment_info", paymentInfo); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "606c7244-e577-40e2-a23a-456127035fc4", err.Error())
	}

	discounts := it.GetDiscounts()
	discountAmount := it.GetDiscountAmount()
	if err := checkoutOrder.Set("discount", discountAmount); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7f124c39-3928-4525-966d-58ebd75f4739", err.Error())
	}
	if err := checkoutOrder.Set("discounts", discounts); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c06c9a96-ff20-4a2e-a793-1911a75827a4", err.Error())
	}

	taxes := it.GetTaxes()
	taxAmount := it.GetTaxAmount()
	if err := checkoutOrder.Set("tax_amount", taxAmount); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9e58b4bb-fe0c-458a-aa2f-e0a87bbdddd7", err.Error())
	}
	if err := checkoutOrder.Set("taxes", taxes); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "97a4c133-a8f6-4f59-b996-b5b9baeca46a", err.Error())
	}

	if err := checkoutOrder.Set("shipping_amount", it.GetShippingAmount()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c592435a-71fc-45cb-bd7a-18790fe616a6", err.Error())
	}

	// remove order items, and add new from current cart with new description
	err := checkoutOrder.RemoveAllItems()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderDescription := ""

	for _, cartItem := range cartItems {
		orderItem, err := checkoutOrder.AddItem(cartItem.GetProductID(), cartItem.GetQty(), cartItem.GetOptions())
		if err := orderItem.Set("idx", cartItem.GetIdx()); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b1423471-070a-430e-9af5-1143133afba4", err.Error())
		}
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if orderDescription != "" {
			orderDescription += ", "
		}
		orderDescription += fmt.Sprintf("%dx %s", cartItem.GetQty(), orderItem.GetName())
	}

	if err := checkoutOrder.Set("description", orderDescription); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f30e65ba-e0a4-4198-960a-0a0d1392dbb4", err.Error())
	}

	err = checkoutOrder.CalculateTotals()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = checkoutOrder.SetStatus(order.ConstOrderStatusPending)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if err := it.SetOrder(checkoutOrder); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "96f10bcc-4f5e-49a3-b043-678eadd7d327", err.Error())
	}
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// trying to process payment
	//--------------------------
	paymentDetails := make(map[string]interface{})
	if currentSession := it.GetSession(); currentSession != nil {
		paymentDetails["sessionID"] = currentSession.GetID()
	}
	paymentDetails["cc"] = it.GetInfo("cc")
	paymentDetails["extra"] = map[string]interface{}{
		"email":        checkoutOrder.Get("customer_email"),
		"billing_name": checkoutOrder.GetBillingAddress().GetFirstName() + " " + checkoutOrder.GetBillingAddress().GetLastName(),
	}

	result, err := paymentMethod.Authorize(checkoutOrder, paymentDetails)
	if err != nil {
		if err := checkoutOrder.SetStatus(order.ConstOrderStatusNew); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "93aa7713-8dc5-4db8-81e9-24c219641908", err.Error())
		}
		return nil, env.ErrorDispatch(err)
	}

	// Payment method require to return as a result:
	// redirect (with completing of checkout after payment processing)
	// or payment info for order
	switch value := result.(type) {
	case api.StructRestRedirect:
		return result, nil

	case map[string]interface{}:
		return it.SubmitFinish(value)
	}

	return it.SubmitFinish(nil)
}

// SubmitFinish finishes processing of submit (required for payment methods to finish with this call?)
func (it *DefaultCheckout) SubmitFinish(paymentInfo map[string]interface{}) (interface{}, error) {

	checkoutOrder := it.GetOrder()
	if checkoutOrder == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6372e487-7d43-4da7-a08d-a4d743baa83c", "Order not present in checkout")
	}

	// track for order status to prevent double executing of checkout success
	previousOrderStatus := checkoutOrder.GetStatus()
	if err := checkoutOrder.SetStatus(order.ConstOrderStatusProcessed); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f4007de4-3426-46da-a911-e70133a5d02c", err.Error())
	}

	// updating payment info of order with info given from payment method
	if paymentInfo != nil {
		currentPaymentInfo := utils.InterfaceToMap(checkoutOrder.Get("payment_info"))
		for key, value := range paymentInfo {
			currentPaymentInfo[key] = value
		}

		if err := checkoutOrder.Set("payment_info", currentPaymentInfo); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6aa02b33-100e-4f75-b54b-6251b9391c7c", err.Error())
		}
	}

	if err := checkoutOrder.Set("created_at", time.Now()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ba3c0e93-600b-49fd-b8c6-b7a132452edf", err.Error())
	}
	if utils.InterfaceToString(checkoutOrder.Get("session_id")) == "" {
		if err := checkoutOrder.Set("session_id", it.GetInfo("session_id")); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4d2106e2-e32c-435d-864c-5dada3970aff", err.Error())
		}
	}

	err := checkoutOrder.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := checkoutOrder.ToHashMap()
	var orderItems []map[string]interface{}

	for _, orderItem := range checkoutOrder.GetItems() {
		options := make(map[string]interface{})

		for optionName, optionKeys := range orderItem.GetOptions() {
			optionMap := utils.InterfaceToMap(optionKeys)
			options[optionName] = optionMap["value"]
		}
		orderItems = append(orderItems, map[string]interface{}{
			"name":    orderItem.GetName(),
			"options": options,
			"sku":     orderItem.GetSku(),
			"qty":     orderItem.GetQty(),
			"price":   orderItem.GetPrice()})
	}

	result["items"] = orderItems

	// return order in map if order has already been processed or marked completed by merchant
	if previousOrderStatus == order.ConstOrderStatusProcessed || previousOrderStatus == order.ConstOrderStatusCancelled {
		return result, nil
	}

	return result, it.CheckoutSuccess(checkoutOrder, it.GetSession())
}
