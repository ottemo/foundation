package stripe

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
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
	return checkout.ConstPaymentTypeCreditCard
}

func (it *Payment) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return it.ConfigIsEnabled()
}

func (it *Payment) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return false
}

func (it *Payment) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	stripe.Key = it.ConfigAPIKey()

	// Create a charge https://stripe.com/docs/api/go#create_charge
	chargeParams := &stripe.ChargeParams{
		// Amount is in cents
		Amount:   uint64(orderInstance.GetGrandTotal() * 100),
		Currency: "usd",

		//TODO: Recommendation to pass in the customer email
		// Meta: map[string]string{
		// 	email: utils.InterfaceToString(orderInstance.Get("customer_email")),
		// },
	}

	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	ccMonth := utils.InterfaceToString(ccInfo["expire_month"])
	ccYear := utils.InterfaceToString(ccInfo["expire_year"])
	ccNumber := utils.InterfaceToString(ccInfo["number"])
	ccName := orderInstance.GetBillingAddress().GetFirstName() + " " + orderInstance.GetBillingAddress().GetLastName()

	cardParams := stripe.CardParams{
		Month:  ccMonth,
		Year:   ccYear,
		Number: ccNumber,
		Name:   ccName,
		// CVC: "", // Optional, highly recommended, not currently passed up

		// Address1: "",
		// Address2: "",
		// City:     "",
		// State:    "",
		// Zip:      "",
		// Country:  "",
	}

	// Must attach either `customer` or `source` to charge
	chargeParams.SetSource(&cardParams)

	ch, err := charge.New(chargeParams)
	if err != nil {
		env.LogEvent(env.LogFields{"err": err}, "charge error")
		return nil, err
	}

	//TODO: comment out?
	env.LogEvent(env.LogFields{"chargeResponse": ch}, "charge response")

	// Assemble the response information
	//TODO: Normalize keys, and values
	orderPaymentInfo := map[string]interface{}{
		"transactionID": ch.ID,
		"ccLastFour":    ch.Source.Card.LastFour,
		"ccMonth":       ch.Source.Card.Month,
		"ccYear":        ch.Source.Card.Year,
		"ccBrand":       ch.Source.Card.Brand,
	}

	// return the response
	return orderPaymentInfo, nil
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
