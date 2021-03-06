package stripe

import (
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"

	"github.com/ottemo/commerce/app/models/checkout"
	"github.com/ottemo/commerce/app/models/order"
	"github.com/ottemo/commerce/app/models/visitor"
)

// GetCode will return the Stripe payment method code
func (it *Payment) GetCode() string {
	return ConstPaymentCode
}

// GetInternalName returns the internal payment method name for Stripe
func (it *Payment) GetInternalName() string {
	return ConstPaymentName
}

// GetName returns the payment method name used for Stripe in checkout
func (it *Payment) GetName() string {
	return it.ConfigNameInCheckout()
}

// GetType returns the credit card type used for payment
func (it *Payment) GetType() string {
	return checkout.ConstPaymentTypeCreditCard
}

// IsAllowed is a flag to check if the Stripe payment method is enabled in the current store
func (it *Payment) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return it.ConfigIsEnabled()
}

// IsTokenable is a flag to indicate if the Stripe payment method supports tokens
func (it *Payment) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return true
}

// Authorize is the method used to validate a visitor's card and associated address data
// - it also allows us to create a token for the card
// - the visitor's card is also authorized for the amount of the order in anticipation of fulfillment
func (it *Payment) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	// Set our api key, applies to any http calls
	stripe.Key = it.ConfigAPIKey()

	// Check if we are just supposed to create a Customer (aka a token)
	action := paymentInfo[checkout.ConstPaymentActionTypeKey]
	isCreateToken := utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken
	if isCreateToken {
		// NOTE: `orderInstance = nil` when creating a token

		// 1. Get our customer token
		extra := utils.InterfaceToMap(paymentInfo["extra"])
		visitorID := utils.InterfaceToString(extra["visitor_id"])
		stripeCID := getStripeCustomerToken(visitorID)
		if stripeCID == "" {
			email := utils.InterfaceToString(extra["email"])
			// 2. We don't have a stripe id on file, make a new customer
			c, err := customer.New(&stripe.CustomerParams{
				Email: &email,
				// TODO: coupons?
			})
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			stripeCID = c.ID
		}

		// 3. Create a card
		ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
		ccInfo["billing_name"] = extra["billing_name"]
		cp, err := getCardParams(ccInfo, stripeCID)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		ca, err := card.New(cp)
		// env.LogEvent(env.LogFields{"api_response": ca, "err": err}, "card")
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// This response looks like our normal authorize response
		// but this map is translated into other keys to store a token
		result := map[string]interface{}{
			"transactionID":      ca.ID,                        // token_id
			"creditCardLastFour": ca.Last4,                     // number
			"creditCardType":     getCCBrand(string(ca.Brand)), // type
			"creditCardExp":      formatCardExp(*ca),           // expiration_date
			"customerID":         stripeCID,                    // customer_id
		}

		return result, nil
	}

	// Charging: https://stripe.com/docs/api/go#create_charge
	var ch *stripe.Charge
	ccInfo := paymentInfo["cc"]

	// Token Charge
	// - we have a Customer, and a Card
	// - create a Charge with the Card as the Source
	// - must reference Customer
	// - email is stored on the Customer
	if creditCard, ok := ccInfo.(visitor.InterfaceVisitorCard); ok && creditCard != nil {
		var err error
		cardID := creditCard.GetToken()
		stripeCID := creditCard.GetCustomerID()

		if cardID == "" || stripeCID == "" {
			err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "02128bc6-83d6-4c12-ae90-900a94adb3ad", "looks like we want to charge a token, but we don't have the fields we need")
			return nil, env.ErrorDispatch(err)
		}

		currency := "usd"
		amount := int64(orderInstance.GetGrandTotal() * 100)
		customer := stripeCID

		chParams := stripe.ChargeParams{
			Currency: &currency,
			Amount:   &amount,   // Amount is in cents
			Customer: &customer, // Mandatory
		}
		if err := chParams.SetSource(cardID); err != nil {
			_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "329ddd35-8fdc-4681-9a02-06290a405073", err.Error())
		}

		ch, err = charge.New(&chParams)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

	} else {
		// Regular Charge
		// - don't create a customer, or store a token
		// - email is stored on the charge's meta hashmap
		var err error

		currency := "usd"
		amount := int64(orderInstance.GetGrandTotal() * 100)

		chargeParams := stripe.ChargeParams{
			Currency: &currency,
			Amount:   &amount, // Amount is in cents
		}
		chargeParams.AddMetadata("email", utils.InterfaceToString(orderInstance.Get("customer_email")))

		// Must attach either `customer` or `source` to charge
		// source can be either a `token` or `cardParams`
		ccInfo := utils.InterfaceToMap(paymentInfo["cc"])

		if ba := orderInstance.GetBillingAddress(); ba != nil {
			ccInfo["billing_name"] = ba.GetFirstName() + " " + ba.GetLastName()
		}

		cp, err := getCardParams(ccInfo, "")
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		if err := chargeParams.SetSource(cp); err != nil {
			_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "ed9f2315-7355-4e9f-959e-e2ac1f737e9f", err.Error())
		}

		ch, err = charge.New(&chargeParams)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// Assemble the response
	orderPaymentInfo := map[string]interface{}{
		"transactionID":     ch.ID,
		"creditCardNumbers": ch.Source.Card.Last4,
		"creditCardExp":     formatCardExp(*ch.Source.Card),
		"creditCardType":    getCCBrand(string(ch.Source.Card.Brand)),
	}

	return orderPaymentInfo, nil
}

// returns string of mmyy
func formatCardExp(c stripe.Card) string {
	exp := utils.InterfaceToString(c.ExpMonth)

	// pad with a zero
	if c.ExpMonth < 10 {
		exp = "0" + exp
	}

	// append the last two year digits
	y := utils.InterfaceToString(c.ExpYear)
	if len(y) == 4 {
		exp = exp + y[:2]
	} else {
		err := env.ErrorNew(ConstErrorModule, 1, "0a17b25a-4155-487a-82ad-dfb4b654eba8", "unexpected year length coming back from stripe "+y)
		_ = env.ErrorDispatch(err)
	}

	return exp
}

// getCardParams Assemble the stripe.CardParams based on the ccInfo we have
// - validates cvc
// - optionally sets customer id
// - optionally sets name from ccInfo["billing_name"]
func getCardParams(ccInfo map[string]interface{}, stripeCID string) (*stripe.CardParams, error) {

	ccCVC := utils.InterfaceToString(ccInfo["cvc"])
	if ccCVC == "" {
		err := env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
		return &stripe.CardParams{}, err
	}

	cardNumber := utils.InterfaceToString(ccInfo["number"])
	expMonth := utils.InterfaceToString(ccInfo["expire_month"])
	expYear := utils.InterfaceToString(ccInfo["expire_year"])
	billingName := utils.InterfaceToString(ccInfo["billing_name"])

	cp := &stripe.CardParams{
		Number:   &cardNumber,
		ExpMonth: &expMonth ,
		ExpYear:  &expYear,
		CVC:      &ccCVC, // Optional, highly recommended

		// might not be passed in
		Customer: &stripeCID,
		Name:     &billingName, // Optional

		// Address fields can be passed here as well to aid in fraud prevention
	}

	return cp, nil
}

// Delete saved card from the payment system.
func (it *Payment) DeleteSavedCard(token visitor.InterfaceVisitorCard) (interface{}, error) {
	// Set our api key, applies to any http calls
	stripe.Key = it.ConfigAPIKey()

	customer := token.GetCustomerID()

	card, err := card.Del(
		token.GetToken(),
		&stripe.CardParams{Customer: &customer},
	)

	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, 1, "05199a06-7bd4-49b6-9fb0-0f1589a9cd74", err.Error())
	}

	return card, nil
}

// Capture is the payment method used to capture authorized funds.  **This method is for future use**
func (it *Payment) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, 1, "0cedd382-392d-4e06-86a9-83276c92de13", "called but not implemented")
}

// Refund is the payment method used to refund a visitor on behalf of a merchant. **This method is for future use**
func (it *Payment) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, 1, "c8768719-80ab-453d-b52e-513dfb4aab22", "called but not implemented")
}

// Void is the payment method used to cancel a visitor transaction before funds have been collected.  **This method is for future use**
func (it *Payment) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, 1, "4194a950-18fd-4b0d-96e6-e33e930f4320", "called but not implemented")
}
