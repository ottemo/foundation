package stripe

import (
	"fmt"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
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
	isTokenable := true
	fmt.Println("stripe - isTokenable called", isTokenable)
	return isTokenable
}

func (it *Payment) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	// Set our api key
	stripe.Key = it.ConfigAPIKey()

	// Check if we are just supposed to create a Customer (aka a token)
	action := paymentInfo[checkout.ConstPaymentActionTypeKey]
	isCreateToken := utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken
	if isCreateToken {
		// NOTE orderInstance = nil when creating a token
		fmt.Println("Authorize - isCreateToken", isCreateToken)

		c, err := getCustomer(paymentInfo)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		result := map[string]interface{}{
			"transactionID":      c.ID, // becomes the token id
			"creditCardLastFour": c.Sources.Values[0].Card.LastFour,
			"creditCardType":     getCCBrand(string(c.Sources.Values[0].Card.Brand)),
			"creditCardExp":      formatCardExp(*c.Sources.Values[0].Card),
			// "responseMessage":    responseMessage,
			// "responseResult":     responseResult,
		}
		return result, nil
	}

	// Check if the cc info being passed in is a credit card
	stripeCID := ""
	if ccInfo, present := paymentInfo["cc"]; present {
		if creditCard, ok := ccInfo.(visitor.InterfaceVisitorCard); ok && creditCard != nil {
			stripeCID = creditCard.GetToken()
			fmt.Println("we have a token we can use: ", stripeCID)
		}
	}

	// We don't have a token for the user, error out
	if stripeCID == "" {
		return nil, env.ErrorDispatch(env.ErrorNew(ConstErrorModule, 1, "02128bc6-83d6-4c12-ae90-900a94adb3ad", "Stripe Authorize called without a valid token"))
	}

	// Assemble charge - https://stripe.com/docs/api/go#create_charge
	ch, err := charge.New(&stripe.ChargeParams{
		Amount:   uint64(orderInstance.GetGrandTotal() * 100), // Amount is in cents
		Currency: "usd",
		Customer: stripeCID,
		Email:    utils.InterfaceToString(orderInstance.Get("customer_email")),
	})
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	env.LogEvent(env.LogFields{"api_response": ch}, "charge") //TODO: COMMENT OUT

	// Assemble the response
	orderPaymentInfo := map[string]interface{}{
		"transactionID":     ch.ID,
		"creditCardNumbers": ch.Source.Card.LastFour,
		"creditCardExp":     formatCardExp(*ch.Source.Card),
		"creditCardType":    getCCBrand(string(ch.Source.Card.Brand)),
	}

	return orderPaymentInfo, nil
}

// returns mmyy
func formatCardExp(c stripe.Card) string {
	ccExp := utils.InterfaceToString(c.Month)
	if c.Month < 10 {
		ccExp = "0" + ccExp
	}
	ccExp = ccExp + utils.InterfaceToString(c.Year)[:2]

	return ccExp
}

func getCustomer(paymentInfo map[string]interface{}) (stripe.Customer, error) {

	// Assemble card params
	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
	ccCVC := utils.InterfaceToString(ccInfo["cvc"])
	if ccCVC == "" {
		err := env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
		return stripe.Customer{}, err
	}
	// ccName := orderInstance.GetBillingAddress().GetFirstName() + " " + orderInstance.GetBillingAddress().GetLastName()

	// Email: email, // don't have an easy way to access this
	// TODO: coupons?
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource(&stripe.CardParams{
		Number: utils.InterfaceToString(ccInfo["number"]),
		Month:  utils.InterfaceToString(ccInfo["expire_month"]),
		Year:   utils.InterfaceToString(ccInfo["expire_year"]),
		CVC:    ccCVC, // Optional, highly recommended
		// Name:   ccName, // Optional
		// Address fields can be passed here as well to aid in fraud prevention
	})

	c, err := customer.New(customerParams)
	if err != nil {
		return stripe.Customer{}, err
	}

	env.LogEvent(env.LogFields{"api_response": c}, "customer") // TODO: COMMENT OUT

	if c.Sources.Count > 1 {
		env.LogEvent(env.LogFields{"customer": c}, "stripe customer has multiple cards")
	}

	// dereference the pointer
	return *c, nil
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
