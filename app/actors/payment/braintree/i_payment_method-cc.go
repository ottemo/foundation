package braintree

import (
	//"github.com/lionelbarrow/braintree-go"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	//"github.com/ottemo/foundation/api"
	//"github.com/ottemo/foundation/app"
	"fmt"
	"github.com/lionelbarrow/braintree-go"
	//"github.com/stripe/stripe-go/customer"
	//"github.com/stripe/stripe-go/card"
	//"github.com/stripe/stripe-go/charge"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetCode returns payment method code for use in business logic
func (it *braintreeCCMethod) GetCode() string {
	return constCCMethodCode
}

// GetInternalName returns the human readable name of the payment method
func (it *braintreeCCMethod) GetInternalName() string {
	return constCCMethodInternalName
}

// GetName returns the user customized name of the payment method
func (it *braintreeCCMethod) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(constCCMethodConfigPathName))
}

// GetType returns type of payment method according to "github.com/ottemo/foundation/app/models/checkout"
func (it *braintreeCCMethod) GetType() string {
	return checkout.ConstPaymentTypeCreditCard
}

// IsAllowed checks for payment method applicability
func (it *braintreeCCMethod) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(constCCMethodConfigPathEnabled))
}

// IsTokenable returns possibility to save token for this payment method
func (it *braintreeCCMethod) IsTokenable(checkoutInstance checkout.InterfaceCheckout) bool {
	return (true)
}

// Authorize makes payment method authorize operations
//  - just create token if set in paymentInfo
//  - otherwise create transaction
func (it *braintreeCCMethod) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	// paymentResult

	braintreeInstance := braintree.New(
		utils.InterfaceToString(env.ConfigGetValue(constGeneralConfigPathEnvironment)),
		utils.InterfaceToString(env.ConfigGetValue(constGeneralConfigPathMerchantID)),
		utils.InterfaceToString(env.ConfigGetValue(constGeneralConfigPathPublicKey)),
		utils.InterfaceToString(env.ConfigGetValue(constGeneralConfigPathPrivateKey)),
	)

	action := paymentInfo[checkout.ConstPaymentActionTypeKey]
	isCreateToken := utils.InterfaceToString(action) == checkout.ConstPaymentActionTypeCreateToken

	if isCreateToken {
		// NOTE: `orderInstance = nil` when creating a token

		// 1. Get our customer token
		extra := utils.InterfaceToMap(paymentInfo["extra"])
		visitorID := utils.InterfaceToString(extra["visitor_id"])
		braintreeClientID := getBraintreeCustomerToken(visitorID)

		if braintreeClientID == "" {
			// 2. We don't have a braintree client id on file, make a new customer
			customerPtr, err := braintreeInstance.Customer().Create(&braintree.Customer{
				Email: utils.InterfaceToString(extra["email"]), // TODO: add more info (is it required)
			})
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			braintreeClientID = customerPtr.Id
		}

		// 3. Create a card
		creditCardInfo := utils.InterfaceToMap(paymentInfo["cc"])
		creditCardInfo["billing_name"] = extra["billing_name"]

		creditCardCVC := utils.InterfaceToString(creditCardInfo["cvc"])
		if creditCardCVC == "" {
			return nil, env.ErrorDispatch(env.ErrorNew(constErrorModule, constErrorLevel, "bd0a78bf-065a-462b-92c7-d5a1529797c4", "CVC field was left empty"))
		}

		newCreditCard := &braintree.CreditCard{
			CustomerId:      braintreeClientID,
			Number:          utils.InterfaceToString(creditCardInfo["number"]),
			ExpirationYear:  utils.InterfaceToString(creditCardInfo["expire_year"]),
			ExpirationMonth: utils.InterfaceToString(creditCardInfo["expire_month"]),
			CVV:             creditCardCVC,
			Options: &braintree.CreditCardOptions{
				VerifyCard: true,
			},
		}

		ca, err := braintreeInstance.CreditCard().Create(newCreditCard)
		// env.LogEvent(env.LogFields{"api_response": ca, "err": err}, "card")
		if err != nil {
			fmt.Println("\nERROR bt.CreditCard().Create(cp)\n", err, "\n")
			return nil, env.ErrorDispatch(err)
		}

		fmt.Println("\n--- ca: ", ca, "\n\n", utils.InterfaceToString(ca))

		// This response looks like our normal authorize response
		// but this map is translated into other keys to store a token
		tokenCreationResult := map[string]interface{}{
			"transactionID":      ca.Token,                      // token_id
			"creditCardLastFour": ca.Last4,                      // number
			"creditCardType":     ca.CardType,                   // type
			"creditCardExp":      formatCardExpirationDate(*ca), // expiration_date
			"customerID":         ca.CustomerId,                 // customer_id
		}

		fmt.Println("\n--- result: ", tokenCreationResult, "\n\n", utils.InterfaceToString(tokenCreationResult))

		return tokenCreationResult, nil
	}

	// Charging: https://stripe.com/docs/api/go#create_charge
	//var ch *stripe.Charge
	var tr *braintree.Transaction
	ccInfo := paymentInfo["cc"]
	//ccInfoMap := utils.InterfaceToMap(ccInfo)

	// Token Charge
	// - we have a Customer, and a Card
	// - create a Charge with the Card as the Source
	// - must reference Customer
	// - email is stored on the Customer
	if creditCard, ok := ccInfo.(visitor.InterfaceVisitorCard); ok && creditCard != nil {
		fmt.Println("\n--- creditCard, ok: ", utils.InterfaceToString(creditCard), "\n")
		var err error
		cardID := creditCard.GetToken()
		stripeCID := creditCard.GetCustomerID()

		if cardID == "" || stripeCID == "" {
			fmt.Println("cardID == '' || stripeCID == ''")
			err := env.ErrorNew(constErrorModule, env.ConstErrorLevelStartStop, "02128bc6-83d6-4c12-ae90-900a94adb3ad", "looks like we want to charge a token, but we don't have the fields we need")
			return nil, env.ErrorDispatch(err)
		}

		//chParams := stripe.ChargeParams{
		//	Currency: "usd",
		//	Amount:   uint64(orderInstance.GetGrandTotal() * 100), // Amount is in cents
		//	Customer: stripeCID,                                   // Mandatory
		//}
		//ccCVC := utils.InterfaceToString(ccInfoMap["cvc"])
		//if ccCVC == "" {
		//	fmt.Println("ccCVC == ''")
		//	err := env.ErrorNew(ConstErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
		//	return nil, env.ErrorDispatch(err)
		//}

		cc, err := braintreeInstance.CreditCard().Find(cardID)
		if err != nil {
			fmt.Println("\n--- Can not find cc.")
			return nil, env.ErrorDispatch(err)
		}
		fmt.Println("\n--- found creditCard, ok: ", utils.InterfaceToString(cc), "\n")

		//cc := &braintree.CreditCard{
		//	//Number:         utils.InterfaceToString(ccInfoMap["number"]),
		//	//ExpirationYear:utils.InterfaceToString(ccInfoMap["expire_year"]),
		//	//ExpirationMonth: utils.InterfaceToString(ccInfoMap["expire_month"]),
		//	//CVV:            ccCVC,
		//	//ExpirationYear:  "25",
		//	//CustomerId:stripeCID,
		//	Token:cardID,
		//	Options: &braintree.CreditCardOptions{
		//		VerifyCard: true,
		//		//FailOnDuplicatePaymentMethod: true,
		//	},
		//}

		tx := &braintree.Transaction{
			Type: "sale",
			//Amount: uint64(orderInstance.GetGrandTotal() * 100),
			Amount: braintree.NewDecimal(int64(orderInstance.GetGrandTotal()*100), 2),
			//CustomerID: utils.InterfaceToString(requestData["x_customer_id"]),
			//PaymentMethodNonce: utils.InterfaceToString(requestData["nonce"]),
			CustomerID: stripeCID,
			//CreditCard: cc,
			PaymentMethodToken: cardID,

			Options: &braintree.TransactionOptions{
				SubmitForSettlement: true,
				StoreInVault:        true,
			},
		}
		//chParams.SetSource(cardID)
		fmt.Println("\n--- tx: ", tx, "\n\n", utils.InterfaceToString(tx))

		//ch, err = charge.New(&chParams)
		tr, err = braintreeInstance.Transaction().Create(tx)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		fmt.Println("\n--- tr: ", tr, "\n\n", utils.InterfaceToString(tr))

	} else {
		//fmt.Println("Regular Charge STOP")
		//return nil, env.ErrorDispatch(*new(error))
		//// Regular Charge
		//// - don't create a customer, or store a token
		//// - email is stored on the charge's meta hashmap
		var err error
		//chargeParams := stripe.ChargeParams{
		//	Currency: "usd",
		//	Amount:   uint64(orderInstance.GetGrandTotal() * 100), // Amount is in cents
		//}

		//// Must attach either `customer` or `source` to charge
		//// source can be either a `token` or `cardParams`
		ccInfo := utils.InterfaceToMap(paymentInfo["cc"])
		//
		if ba := orderInstance.GetBillingAddress(); ba != nil {
			ccInfo["billing_name"] = ba.GetFirstName() + " " + ba.GetLastName()
		}
		//
		ccCVC := utils.InterfaceToString(ccInfo["cvc"])
		if ccCVC == "" {
			err := env.ErrorNew(constErrorModule, 1, "15edae76-1d3e-4e7a-a474-75ffb61d26cb", "CVC field was left empty")
			return nil, env.ErrorDispatch(err)
		}

		cp := &braintree.CreditCard{
			//CustomerId: 	stripeCID,
			Number:          utils.InterfaceToString(ccInfo["number"]),
			ExpirationYear:  utils.InterfaceToString(ccInfo["expire_year"]),
			ExpirationMonth: utils.InterfaceToString(ccInfo["expire_month"]),
			CVV:             ccCVC,
			//Options: &braintree.CreditCardOptions{
			//	VerifyCard: true,
			//},
		}
		fmt.Println("\n--- cp: ", cp, "\n\n", utils.InterfaceToString(cp))
		//chargeParams.SetSource(cp)

		tx := &braintree.Transaction{
			Type: "sale",
			//Amount: uint64(orderInstance.GetGrandTotal() * 100),
			Amount: braintree.NewDecimal(int64(orderInstance.GetGrandTotal()*100), 2),
			//CustomerID: utils.InterfaceToString(requestData["x_customer_id"]),
			//PaymentMethodNonce: utils.InterfaceToString(requestData["nonce"]),
			//CustomerID:stripeCID,
			//Customer:&braintree.Customer{
			//	Email:	utils.InterfaceToString(orderInstance.Get("customer_email")),
			//},
			CreditCard: cp,
			//PaymentMethodToken:cardID,

			Options: &braintree.TransactionOptions{
				SubmitForSettlement: true,
				//StoreInVault: true,
			},
		}

		//chargeParams.AddMeta("email", utils.InterfaceToString(orderInstance.Get("customer_email")))
		//
		//
		//ch, err = charge.New(&chargeParams)
		//if err != nil {
		//	return nil, env.ErrorDispatch(err)
		//}
		fmt.Println("\n--- tx: ", tx, "\n\n", utils.InterfaceToString(tx))

		//ch, err = charge.New(&chParams)
		tr, err = braintreeInstance.Transaction().Create(tx)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		fmt.Println("\n--- tr: ", tr, "\n\n", utils.InterfaceToString(tr))
	}

	// Assemble the response
	fmt.Println("\n--- Assemble the response\n")
	//fmt.Println("\n--- --- tr.CreditCard.Token\n", tr.CreditCard.Token)
	//fmt.Println("\n--- --- tr.CreditCard.Last4\n", tr.CreditCard.Last4)
	//fmt.Println("\n--- --- *tr.CreditCard\n", *tr.CreditCard)
	//fmt.Println("\n--- --- tr.CreditCard.CardType\n", tr.CreditCard.CardType)
	//fmt.Println("\n--- --- tr.Customer.Id\n", tr.Customer.Id)
	paymentResult := map[string]interface{}{
		"transactionID":      tr.CreditCard.Token,
		"creditCardLastFour": tr.CreditCard.Last4,
		"creditCardExp":      formatCardExpirationDate(*tr.CreditCard),
		"creditCardType":     tr.CreditCard.CardType,
		"customerID":         tr.Customer.Id,
	}
	fmt.Println("\n--- orderPaymentInfo: ", paymentResult, "\n\n", utils.InterfaceToString(paymentResult))

	return paymentResult, nil
}

// Capture makes payment method capture operation
// - at time of implementation this method is not used anywhere
func (it *braintreeCCMethod) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(constErrorModule, constErrorLevel, "772bc737-f025-4c81-a85a-c10efb67e1b3", " Capture method not implemented")
}

// Refund will return funds on the given order
// - at time of implementation this method is not used anywhere
func (it *braintreeCCMethod) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(constErrorModule, constErrorLevel, "26febf8b-7e26-44d4-bfb4-e9b29126fe5a", "Refund method not implemented")
}

// Void will mark the order and capture as void
// - at time of implementation this method is not used anywhere
func (it *braintreeCCMethod) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(constErrorModule, constErrorLevel, "561e0cc8-3bee-4ec4-bf80-585fa566abd4", "Void method not implemented")
}
