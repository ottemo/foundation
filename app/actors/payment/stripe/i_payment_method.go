package stripe

import (
	"fmt"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

func (it *Payment) GetCode() string {
	return ConstPaymentCode
}

func (it *Payment) GetInternalName() string {
	return ConstPaymentName
}

func (it *Payment) GetName() string {
	return it.ConfigNameInCheckout()
}

func (it *Payment) GetType() string {
	return checkout.ConstPaymentTypeCreditCard
}

func (it *Payment) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return it.ConfigIsEnabled()
}

func (it *Payment) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	fmt.Println("stripe - isTokenable called")
	return true
}

func (it *Payment) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	// Check if we are just supposed to create a customer (create a token)
	action, actionPresent := paymentInfo[checkout.ConstPaymentActionTypeKey]
	if actionPresent && utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken {
		return it.AuthorizeZeroAmount(nil, paymentInfo)
	}

	// Set our api key
	stripe.Key = it.ConfigAPIKey()

	email := utils.InterfaceToString(orderInstance.Get("customer_email"))

	// Check if we have a token for the user
	stripeCID := ""
	if stripeCID == "" {
		// We don't have a token for the user, create one
		cardParams, err := getCardParams(orderInstance, paymentInfo)
		if err != nil {
			env.LogEvent(env.LogFields{"err": err}, "build card params") // TODO: COMMENT OUT
			return nil, env.ErrorDispatch(err)
		}

		customerParams := &stripe.CustomerParams{Email: email} // TODO: coupons?
		customerParams.SetSource(&cardParams)

		c, err := customer.New(customerParams)
		if err != nil {
			env.LogEvent(env.LogFields{"err": err}, "cusotmer") // TODO: COMMENT OUT
			return nil, env.ErrorDispatch(err)
		}
		env.LogEvent(env.LogFields{"api_response": c}, "cusotmer") // TODO: COMMENT OUT

		// TODO: SAVE THE TOKEN ON THE USER

		stripeCID = c.ID
	}

	// Assemble charge - https://stripe.com/docs/api/go#create_charge
	chargeParams := &stripe.ChargeParams{
		Amount:   uint64(orderInstance.GetGrandTotal() * 100), // Amount is in cents
		Currency: "usd",
		Customer: stripeCID,
		Email:    email,
	}
	// chargeParams.AddMeta("email", email)
	// chargeParams.Customer(stripeCID)
	env.LogEvent(env.LogFields{"params": chargeParams}, "chargeParams struct") //TODO: COMMENT OUT

	ch, err := charge.New(chargeParams)
	if err != nil {
		env.LogEvent(env.LogFields{"err": err, "charge": ch}, "charge") // TODO: COMMENT OUT
		return nil, env.ErrorDispatch(err)
	}

	env.LogEvent(env.LogFields{"api_response": ch}, "charge") //TODO: COMMENT OUT

	// Assemble the response information
	orderPaymentInfo := map[string]interface{}{
		"transactionID": ch.ID,
		"ccLastFour":    ch.Source.Card.LastFour,
		"ccMonth":       ch.Source.Card.Month,
		"ccYear":        ch.Source.Card.Year,
		"ccBrand":       getCCBrand(string(ch.Source.Card.Brand)),
	}

	return orderPaymentInfo, nil
}

func getCardParams(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (stripe.CardParams, error) {

	// Assemble cardParams
	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	ccNumber := utils.InterfaceToString(ccInfo["number"])
	ccMonth := utils.InterfaceToString(ccInfo["expire_month"])
	ccYear := utils.InterfaceToString(ccInfo["expire_year"])
	ccCVC := utils.InterfaceToString(ccInfo["cvc"])
	ccName := orderInstance.GetBillingAddress().GetFirstName() + " " + orderInstance.GetBillingAddress().GetLastName()

	// We are enforcing cvc
	if ccCVC == "" {
		err := env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
		return stripe.CardParams{}, env.ErrorDispatch(err)
	}

	cardParams := stripe.CardParams{
		Number: ccNumber,
		Month:  ccMonth,
		Year:   ccYear,
		CVC:    ccCVC,  // Optional, highly recommended
		Name:   ccName, // Optional
		// Address fields can be passed here as well to aid in fraud prevention
	}

	return cardParams, nil
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
