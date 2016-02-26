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
		env.LogError(err)
		return nil
	}

	err = shippingAddress.FromHashMap(it.ShippingAddress)
	if err != nil {
		env.LogError(err)
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
		env.LogError(err)
		return nil
	}

	err = billingAddress.FromHashMap(it.BillingAddress)
	if err != nil {
		env.LogError(err)
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
	return it.GetItemSpecificTotal(0, checkout.ConstLabelDiscount) + it.GetItemSpecificTotal(0, checkout.ConstLabelGiftCard)
}

// GetPriceAdjustments collects price adjustments applied for current checkout
func (it *DefaultCheckout) GetPriceAdjustments(label string) []checkout.StructPriceAdjustment {
	var result []checkout.StructPriceAdjustment

	for _, priceAdjustment := range it.priceAdjustments {
		if label == "" || utils.IsInListStr(label, priceAdjustment.Types) {
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

// GetShippingAmount returns shipping price for current checkout
func (it *DefaultCheckout) GetShippingAmount() float64 {
	return it.GetItemSpecificTotal(0, checkout.ConstLabelShipping)
}

// calculateSubtotal it's an element of calculation that provides subtotal amounts
func (it *DefaultCheckout) calculateSubtotal() checkout.StructPriceAdjustment {

	items := it.GetItems()
	result := checkout.StructPriceAdjustment{
		Code:      checkout.ConstLabelSubtotal,
		Label:     checkout.ConstLabelSubtotal,
		Amount:    0,
		IsPercent: false,
		Priority:  checkout.ConstCalculateTargetSubtotal,
		Types:     []string{checkout.ConstLabelSubtotal},
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
			Label:     shippingRate.Name,
			Amount:    shippingRate.Price,
			IsPercent: false,
			Priority:  checkout.ConstCalculateTargetShipping,
			Types:     []string{checkout.ConstLabelShipping},
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

// applyAmount applies amounts to checkout detail calculaltion map
func (it *DefaultCheckout) applyAmount(idx interface{}, label string, amount float64) {
	amount = utils.RoundPrice(amount)
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
			amount = it.GetItemSpecificTotal(0, checkout.ConstLabelGrandTotal) * priceAdjustment.Amount / 100
		}

		// prevent negative values of grand total
		if amount+it.calculateAmount < 0 {
			amount = it.calculateAmount
		}

		// affecting grand total of a cart
		it.applyAmount(0, checkout.ConstLabelGrandTotal, amount)
		totalPriceAdjustmentAmount += utils.RoundPrice(amount)

		// cart details show amount was applied by Types
		for _, label := range priceAdjustment.Types {
			if label != checkout.ConstLabelGrandTotal {
				it.applyAmount(0, label, amount)
			}
		}

	} else {
		for index, amount := range priceAdjustment.PerItem {
			currentItemTotal := it.GetItemSpecificTotal(index, checkout.ConstLabelGrandTotal)
			if priceAdjustment.IsPercent {
				amount = currentItemTotal * priceAdjustment.Amount / 100
			}

			// prevent negative values of grand total per cart
			if amount+it.calculateAmount < 0 {
				amount = it.calculateAmount
			}

			// prevent negative values of grand total per item
			if amount+currentItemTotal < 0 {
				amount = it.calculateAmount
			}

			// adding amount to grand total of current item and full cart
			it.applyAmount(index, checkout.ConstLabelGrandTotal, amount)
			totalPriceAdjustmentAmount += utils.RoundPrice(amount)

			for _, label := range priceAdjustment.Types {
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
		for _, priceAdjustment := range checkout.GetRegisteredPriceAdjustments() {
			for _, priceAdjustmentElement := range priceAdjustment.Calculate(it) {
				priceAdjustments = append(priceAdjustments, priceAdjustmentElement)
			}
		}

		basePoints := map[float64]func() checkout.StructPriceAdjustment{
			checkout.ConstCalculateTargetSubtotal: func() checkout.StructPriceAdjustment { return it.calculateSubtotal() },
			checkout.ConstCalculateTargetShipping: func() checkout.StructPriceAdjustment { return it.calculateShipping() },
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

		newOrder.Set("created_at", currentTime)

		checkoutOrder = newOrder
	}

	// updating order information
	//---------------------------
	checkoutOrder.Set("session_id", it.GetInfo("session_id"))
	checkoutOrder.Set("updated_at", currentTime)

	checkoutOrder.SetStatus(order.ConstOrderStatusNew)

	// If the visitor is logged in that should dictate the email and name
	if currentVisitor != nil {
		checkoutOrder.Set("visitor_id", currentVisitor.GetID())

		checkoutOrder.Set("customer_email", currentVisitor.GetEmail())
		checkoutOrder.Set("customer_name", currentVisitor.GetFullName())
	} else {
		// Visitor is not logged in, maybe we were passed some extra data
		if it.GetInfo("customer_email") != nil {
			orderCustomerEmail := utils.InterfaceToString(it.GetInfo("customer_email"))
			checkoutOrder.Set("customer_email", orderCustomerEmail)
		}
		if it.GetInfo("customer_name") != nil {
			orderCustomerName := utils.InterfaceToString(it.GetInfo("customer_name"))
			checkoutOrder.Set("customer_name", orderCustomerName)
		}
	}

	billingAddress := it.GetBillingAddress().ToHashMap()
	checkoutOrder.Set("billing_address", billingAddress)

	shippingAddress := it.GetShippingAddress().ToHashMap()
	checkoutOrder.Set("shipping_address", shippingAddress)

	shippingInfo := utils.InterfaceToMap(checkoutOrder.Get("shipping_info"))
	shippingInfo["shipping_method_name"] = it.GetShippingMethod().GetName() + "/" + it.GetShippingRate().Name
	if notes := utils.InterfaceToString(it.GetInfo("notes")); notes != "" {
		shippingInfo["notes"] = notes
	}
	checkoutOrder.Set("shipping_info", shippingInfo)
	checkoutOrder.Set("shipping_method", it.GetShippingMethod().GetCode()+"/"+it.GetShippingRate().Code)

	checkoutOrder.Set("cart_id", currentCart.GetID())

	paymentMethod := it.GetPaymentMethod()

	// call for recalculating of all amounts including taxes and discounts
	// cause they can be not calculated at this point in case of direct post
	it.CalculateAmount(checkout.ConstCalculateTargetGrandTotal)

	if !paymentMethod.IsAllowed(it) {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7a5490ee-daa3-42b4-a84a-dade12d103e8", "Payment method not allowed")
	}

	checkoutOrder.Set("payment_method", paymentMethod.GetCode())
	paymentInfo := utils.InterfaceToMap(checkoutOrder.Get("payment_info"))
	paymentInfo["payment_method_name"] = it.GetPaymentMethod().GetName()
	paymentInfo["gift_cards_charged_amount"] = it.GetItemSpecificTotal(0, checkout.ConstLabelGiftCard)
	checkoutOrder.Set("payment_info", paymentInfo)

	discounts := it.GetDiscounts()
	discountAmount := it.GetDiscountAmount()
	checkoutOrder.Set("discount", discountAmount)
	checkoutOrder.Set("discounts", discounts)

	taxes := it.GetTaxes()
	taxAmount := it.GetTaxAmount()
	checkoutOrder.Set("tax_amount", taxAmount)
	checkoutOrder.Set("taxes", taxes)

	checkoutOrder.Set("shipping_amount", it.GetShippingAmount())

	// remove order items, and add new from current cart with new description
	for index := range checkoutOrder.GetItems() {
		err := checkoutOrder.RemoveItem(index + 1)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	orderDescription := ""

	for _, cartItem := range cartItems {
		orderItem, err := checkoutOrder.AddItem(cartItem.GetProductID(), cartItem.GetQty(), cartItem.GetOptions())
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if orderDescription != "" {
			orderDescription += ", "
		}
		orderDescription += fmt.Sprintf("%dx %s", cartItem.GetQty(), orderItem.GetName())

	}

	checkoutOrder.Set("description", orderDescription)

	err := checkoutOrder.CalculateTotals()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = checkoutOrder.SetStatus(order.ConstOrderStatusPending)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	it.SetOrder(checkoutOrder)
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

	result, err := paymentMethod.Authorize(checkoutOrder, paymentDetails)
	if err != nil {
		checkoutOrder.SetStatus(order.ConstOrderStatusNew)
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
	checkoutOrder.SetStatus(order.ConstOrderStatusProcessed)

	// updating payment info of order with info given from payment method
	if paymentInfo != nil {
		currentPaymentInfo := utils.InterfaceToMap(checkoutOrder.Get("payment_info"))
		for key, value := range paymentInfo {
			currentPaymentInfo[key] = value
		}

		checkoutOrder.Set("payment_info", currentPaymentInfo)
	}

	checkoutOrder.Set("created_at", time.Now())
	if utils.InterfaceToString(checkoutOrder.Get("session_id")) == "" {
		checkoutOrder.Set("session_id", it.GetInfo("session_id"))
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
