package braintree

import (
	//"github.com/ottemo/foundation/api"
	//"github.com/ottemo/foundation/env"
	//"fmt"
	//"github.com/ottemo/foundation/utils"
	//"github.com/lionelbarrow/braintree-go"
	//"github.com/ottemo/foundation/app/models/checkout"
	//"strings"
	//"github.com/ottemo/foundation/app"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	//service := api.GetRestService()
	//service.POST("braintree/submit", APISubmit)

	return nil
}

// APISubmit processes Braintree receipt response
// can be used for redirecting customer to it on exit from authorize.net
//   - "x_session" should be specified in request contents with id of existing session
//   - refer to http://www.authorize.net/support/DirectPost_guide.pdf for other fields receipt response should contain
//func APISubmit(context api.InterfaceApplicationContext) (interface{}, error) {
//	requestData, err := api.GetRequestContentAsMap(context)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Println("requestData: ", requestData)
//	fmt.Println("requestData: ", utils.InterfaceToString(requestData))
//
//	bt := braintree.New(
//		braintree.Sandbox,
//		"ddxtcwf5n3hvtz3g",
//		"cfj6fzzrkc898mm6",
//		"24d8738ee7bc4331bbc3bac79f2a54c2",
//	)
//
//	// Find or create card
//	//var cc = &braintree.CreditCard{
//	//	Number:          "4111111111111111",
//	//	//Number:          "4000111111111115",
//	//	//Number:		 "3566002020360505",
//	//	//Number:		 "3566002020360505",
//	//	PaymentMethodNonce: utils.InterfaceToString(requestData["nonce"]),
//	//	ExpirationMonth: "05",
//	//	ExpirationYear:  "25",
//	//	Options: &braintree.CreditCardOptions{
//	//		VerifyCard: true,
//	//		//FailOnDuplicatePaymentMethod: true,
//	//	},
//	//}
//
//	//fmt.Println("\ncc: ", cc, "\n\n", utils.InterfaceToString(cc))
//
//	//bt.CreditCard().Find(token)
//
//
//	var decimal = braintree.NewDecimal(100, 2)
//	err = decimal.UnmarshalText([]byte(utils.InterfaceToString(requestData["x_amount"])))
//	if err != nil {
//		return nil, env.ErrorDispatch(err)
//	}
//
//	tr, err := bt.Transaction().Create(&braintree.Transaction{
//		Type: "sale",
//		Amount: decimal,
//		//CustomerID: utils.InterfaceToString(requestData["x_customer_id"]),
//		PaymentMethodNonce: utils.InterfaceToString(requestData["nonce"]),
//		Options: &braintree.TransactionOptions{
//			SubmitForSettlement: true,
//			StoreInVault: true,
//		},
//		//CreditCard: &braintree.CreditCard{
//		//	//Number:          "4111111111111111",
//		//	//Number:          "4000111111111115",
//		//	//Number:		 "3566002020360505",
//		//	//Number:		 "3566002020360505",
//		//	//CVV:             "123",
//		//	//ExpirationMonth: "05",
//		//	//ExpirationYear:  "25",
//		//	Options: &braintree.CreditCardOptions{
//		//		VerifyCard: true,
//		//		//FailOnDuplicatePaymentMethod: true,
//		//	},
//		//},
//
//		//CreditCard: &braintree.CreditCard{
//		//	Number: "4111111111111111",
//		//	ExpirationDate: "05/25",
//		//},
//	})
//	fmt.Println("transaction", tr)
//	fmt.Println("transaction", utils.InterfaceToString(tr))
//	if err != nil {
//		return nil, env.ErrorDispatch(err)
//	}
//
//	//tr, err = bt.Transaction().Settle(tr.Id)
//	//if err != nil {
//	//	fmt.Println("TRANSACTION not SETTLED")
//	//	return nil, env.ErrorDispatch(err)
//	//}
//
//	fmt.Println("TRANSACTION STATUS: ", tr.Status)
//
//	//------------------------------------------
//	// Here we have CC token
//	// Here we have Customer + ID token
//	//------------------------------------------
//
//	//var visitorCreditCard visitor.InterfaceVisitorCard
//	//if visitorCreditCard != nil && visitorCreditCard.GetID() != "" {
//	//	//orderPaymentInfo["creditCardID"] = visitorCreditCard.GetID()
//	//
//	//	visitorCreditCard.Set("token_id", orderTransactionID)
//	//	visitorCreditCard.Set("token_updated", time.Now())
//	//	visitorCreditCard.Save()
//	//}
//
//
//	//return nil, env.ErrorDispatch(*new(error))
//
//	session, err := api.GetSessionByID(utils.InterfaceToString(requestData["x_session"]), false)
//	if session == nil {
//		return nil, env.ErrorNew(constErrorModule, env.ConstErrorLevelAPI, "48f70911-836f-41ba-9ed9-b2afcb7ca462", "Wrong session ID")
//	}
//	context.SetSession(session)
//
//	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
//	if err != nil {
//		return nil, env.ErrorDispatch(err)
//	}
//
//	checkoutOrder := currentCheckout.GetOrder()
//
//	currentCart := currentCheckout.GetCart()
//	if currentCart == nil {
//		return nil, env.ErrorNew(constErrorModule, env.ConstErrorLevelAPI, "6244e778-a837-4425-849b-fbce26d5b095", "Cart is not specified")
//	}
//	if checkoutOrder != nil {
//
//		orderMap, err := currentCheckout.SubmitFinish(requestData)
//		if err != nil {
//			env.LogError(env.ErrorNew(constErrorModule, env.ConstErrorLevelAPI, "54296509-fc83-447d-9826-3b7a94ea1acb", "Can't proceed submiting order from Authorize relay"))
//		}
//
//		redirectURL := "" //utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathDPMReceiptURL))
//		if strings.TrimSpace(redirectURL) == "" {
//			redirectURL = app.GetStorefrontURL("")
//		}
//
//		env.Log(constLogStorage, env.ConstLogPrefixInfo, "TRANSACTION APPROVED: "+
//			"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
//			"OrderID - "+checkoutOrder.GetID()+", "+
//			"Card  - "+utils.InterfaceToString(requestData["cardType"])+" "+utils.InterfaceToString(requestData["lastTwo"])+", "+
//			"Total - "+utils.InterfaceToString(requestData["x_amount"])+", "+
//			"Transaction ID - "+tr.CreditCard.Token)
//			//"Transaction ID - "+utils.InterfaceToString(requestData["x_trans_id"]))
//
//		return api.StructRestRedirect{Result: orderMap, Location: redirectURL, DoRedirect: true}, err
//	}
//
//
//	fmt.Println("requestData: ", requestData)
//	fmt.Println("requestData: ", utils.InterfaceToString(requestData))
//
//	return nil, env.ErrorNew(constErrorModule, env.ConstErrorLevelAPI, "28ad0ab3-6505-4fc6-9c9f-d07fb556190e", "can't process Braintree response")
//}
