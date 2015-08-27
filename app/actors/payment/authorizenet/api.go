package authorizenet

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("authorizenet/receipt", api.ConstRESTOperationCreate, APIReceipt)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("authorizenet/relay", api.ConstRESTOperationCreate, APIRelay)
	if err != nil {
		return err
	}

	return nil
}

// APIReceipt processes Authorize.net receipt response
//   - "x_session" should be specified in request contents with id of existing session
//   - refer to http://www.authorize.net/support/DirectPost_guide.pdf for other fields receipt response should contain
func APIReceipt(context api.InterfaceApplicationContext) (interface{}, error) {

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, err
	}

	status := requestData["x_response_code"]

	session, err := api.GetSessionByID(utils.InterfaceToString(requestData["x_session"]))
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "48f70911-836f-41ba-9ed9-b2afcb7ca462", "Wrong session ID")
	}
	context.SetSession(session)

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, err
	}

	checkoutOrder := currentCheckout.GetOrder()

	switch status {
	case ConstTransactionApproved:
		{
			currentCart := currentCheckout.GetCart()
			if currentCart == nil {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6244e778-a837-4425-849b-fbce26d5b095", "Cart is not specified")
			}
			if checkoutOrder != nil {

				result, err := currentCheckout.SubmitFinish(requestData)

				env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "TRANSACTION APPROVED: "+
					"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderID - "+checkoutOrder.GetID()+", "+
					"Card  - "+utils.InterfaceToString(requestData["x_card_type"])+" "+utils.InterfaceToString(requestData["x_account_number"])+", "+
					"Total - "+utils.InterfaceToString(requestData["x_amount"])+", "+
					"Transaction ID - "+utils.InterfaceToString(requestData["x_trans_id"]))

				return api.StructRestRedirect{Result: result, Location: app.GetStorefrontURL("account/order/" + checkoutOrder.GetID()), DoRedirect: true}, err
			}
		}
	case ConstTransactionDeclined:
	case ConstTransactionWaitingReview:
	default:
		{
			if checkoutOrder != nil {
				env.Log(ConstLogStorage, env.ConstLogPrefixError, "TRANSACTION NOT APPROVED: "+
					"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderID - "+checkoutOrder.GetID()+", "+
					"Card  - "+utils.InterfaceToString(requestData["x_card_type"])+" "+utils.InterfaceToString(requestData["x_account_number"])+", "+
					"Total - "+utils.InterfaceToString(requestData["x_amount"])+", "+
					"Transaction ID - "+utils.InterfaceToString(requestData["x_trans_id"]))
			}

			return []byte(`<html>
					 <head>
						 <noscript>
						 	<meta http-equiv='refresh' content='1;url=` + app.GetStorefrontURL("checkout") + `'>
						 </noscript>
					 </head>
					 <body>
					 	<h1>Something went wrong</h1>
					 	<p>` + utils.InterfaceToString(requestData["x_response_reason_text"]) + `</p>

						<p><a href="` + app.GetStorefrontURL("checkout") + `">Back to store</a></p>

					 </body>
				</html>`), nil
		}
	}
	if checkoutOrder != nil {
		env.Log(ConstLogStorage, env.ConstLogPrefixError, "TRANSACTION NOT APPROVED: (can't process authorize.net response) "+
			"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
			"OrderID - "+checkoutOrder.GetID()+", "+
			"Card  - "+utils.InterfaceToString(requestData["x_card_type"])+" "+utils.InterfaceToString(requestData["x_account_number"])+", "+
			"Total - "+utils.InterfaceToString(requestData["x_amount"])+", "+
			"Transaction ID - "+utils.InterfaceToString(requestData["x_trans_id"]))
	}
	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "770e9dec-8f59-4e98-857f-e8124bf6771e", "can't process authorize.net response")
}

// APIRelay processes Authorize.net relay response
//   - "x_session" should be specified in request contents with id of existing session
//   - refer to http://www.authorize.net/support/DirectPost_guide.pdf for other fields relay response should contain
func APIRelay(context api.InterfaceApplicationContext) (interface{}, error) {

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, err
	}

	status := requestData["x_response_code"]

	sessionInstance, err := api.GetSessionByID(utils.InterfaceToString(requestData["x_session"]))
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "48f70911-836f-41ba-9ed9-b2afcb7ca462", "Wrong session ID")
	}
	context.SetSession(sessionInstance)

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, err
	}

	checkoutOrder := currentCheckout.GetOrder()

	switch status {
	case ConstTransactionApproved:
		{
			currentCart := currentCheckout.GetCart()
			if currentCart == nil {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6244e778-a837-4425-849b-fbce26d5b095", "Cart is not specified")
			}
			if checkoutOrder != nil {

				result, err := currentCheckout.SubmitFinish(requestData)

				context.SetResponseContentType("text/plain")

				env.Log(ConstLogStorage, env.ConstLogPrefixInfo, "TRANSACTION APPROVED: "+
					"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderID - "+checkoutOrder.GetID()+", "+
					"Card  - "+utils.InterfaceToString(requestData["x_card_type"])+" "+utils.InterfaceToString(requestData["x_account_number"])+", "+
					"Total - "+utils.InterfaceToString(requestData["x_amount"])+", "+
					"Transaction ID - "+utils.InterfaceToString(requestData["x_trans_id"]))

				return api.StructRestRedirect{Result: result, Location: app.GetStorefrontURL("checkout/success/" + checkoutOrder.GetID()), DoRedirect: true}, err
			}
		}
	case ConstTransactionDeclined:
	case ConstTransactionWaitingReview:
	default:
		{
			if checkoutOrder != nil {
				env.Log(ConstLogStorage, env.ConstLogPrefixError, "TRANSACTION NOT APPROVED: "+
					"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderID - "+checkoutOrder.GetID()+", "+
					"Card  - "+utils.InterfaceToString(requestData["x_card_type"])+" "+utils.InterfaceToString(requestData["x_account_number"])+", "+
					"Total - "+utils.InterfaceToString(requestData["x_amount"])+", "+
					"Transaction ID - "+utils.InterfaceToString(requestData["x_trans_id"]))
			}
			return []byte(`<html>
					 <head>
						 <noscript>
						 	<meta http-equiv='refresh' content='1;url=` + app.GetStorefrontURL("checkout") + `'>
						 </noscript>
					 </head>
					 <body>
					 	<h1>Something went wrong</h1>
					 	<p>` + utils.InterfaceToString(requestData["x_response_reason_text"]) + `</p>

						<p><a href="` + app.GetStorefrontURL("checkout") + `">Back to store</a></p>

					 </body>
				</html>`), nil
		}
	}
	if checkoutOrder != nil {
		env.Log(ConstLogStorage, env.ConstLogPrefixError, "TRANSACTION NOT APPROVED: (can't process authorize.net response) "+
			"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
			"OrderID - "+checkoutOrder.GetID()+", "+
			"Card  - "+utils.InterfaceToString(requestData["x_card_type"])+" "+utils.InterfaceToString(requestData["x_account_number"])+", "+
			"Total - "+utils.InterfaceToString(requestData["x_amount"])+", "+
			"Transaction ID - "+utils.InterfaceToString(requestData["x_trans_id"]))
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "770e9dec-8f59-4e98-857f-e8124bf6771e", "can't process authorize.net response")
}
