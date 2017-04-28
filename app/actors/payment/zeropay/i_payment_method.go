package zeropay

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetInternalName returns the name of the payment method
func (it ZeroAmountPayment) GetInternalName() string {
	return ConstPaymentName
}

// GetName returns the user customized name of the payment method
func (it *ZeroAmountPayment) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathName))
}

// GetCode returns payment method code
func (it *ZeroAmountPayment) GetCode() string {
	return ConstPaymentZeroPaymentCode
}

// GetType returns type of payment method
func (it *ZeroAmountPayment) GetType() string {
	return checkout.ConstPaymentTypeSimple
}

// IsAllowed checks for method applicability
func (it *ZeroAmountPayment) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	if checkoutInstance.GetGrandTotal() > 0 || !utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEnabled)) {
		return false
	}
	return true
}

// IsTokenable checks for method applicability
func (it *ZeroAmountPayment) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return false
}

// Authorize makes payment method authorize operation
func (it *ZeroAmountPayment) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	if orderInstance.GetGrandTotal() > 0 {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0a1de8e4-94f8-4e8b-92fd-378b92d7d9fa", "Order amount above zero, please specify allowe payment method.")
	}
	return nil, nil
}

// Delete saved card from the payment system.  **This method is for future use**
func (it *ZeroAmountPayment) DeleteSavedCard(token visitor.InterfaceVisitorCard) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8596a836-0351-475f-9ca5-4fab87c7a74a", "Not implemented")
}

// Capture makes payment method capture operation
func (it *ZeroAmountPayment) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}

// Refund makes payment method refund operation
func (it *ZeroAmountPayment) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}

// Void makes payment method void operation
func (it *ZeroAmountPayment) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}
