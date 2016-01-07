package zeropay

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
)

// GetName returns payment method name
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

// Authorize makes payment method authorize operation
func (it *ZeroAmountPayment) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	if orderInstance.GetGrandTotal() > 0 {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0a1de8e4-94f8-4e8b-92fd-378b92d7d9fa", "Order amount above zero, please specify allowe payment method.")
	}
	return nil, nil
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
