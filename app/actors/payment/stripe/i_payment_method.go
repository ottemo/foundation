package stripe

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	// stripe "github.com/stripe/stripe-go"
)

func (it *Payment) GetCode() string {
	return ConstPaymentCode
}

func (it *Payment) GetInternalName() string {
	return ConstPaymentName
}

func (it *Payment) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathName))
}

func (it *Payment) GetType() string {
	// TODO: NOT SURE
	// checkout.ConstPaymentTypeCreditCard
	return ""
}

func (it *Payment) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return it.ConfigIsEnabled()
}

func (it *Payment) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return false
}

func (it *Payment) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func (it *Payment) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, 1, "05199a06-7bd4-49b6-9fb0-0f1589a9cd74", "called but not implemented")
}

func (it *Payment) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, 1, "c8768719-80ab-453d-b52e-513dfb4aab22", "called but not implemented")
}

func (it *Payment) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, 1, "4194a950-18fd-4b0d-96e6-e33e930f4320", "called but not implemented")
}
